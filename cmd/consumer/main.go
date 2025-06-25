package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

// Efinir var window PowerShell
// $env:MY_VARIABLE = "meu_valor"

func main() {
	connectionString := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	queueName := os.Getenv("SERVICEBUS_QUEUE_NAME")
	workerCount := os.Getenv("WORKER_COUNT")

	if connectionString == "" {
		log.Fatal("A variável de ambiente SERVICEBUS_CONNECTION_STRING não foi definida.")
	}
	if queueName == "" {
		log.Fatal("A variável de ambiente SERVICEBUS_QUEUE_NAME não foi definida.")
	}
	if workerCount == "" {
		workerCount = "100"
	}

	client, err := azservicebus.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		log.Fatalf("Erro ao criar o client: %v", err)
	}
	defer client.Close(context.Background())

	// Controla o número máximo de sessões/goroutines em paralelo
	maxConcurrentSessions := stringToInt(workerCount)
	sem := make(chan struct{}, maxConcurrentSessions)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// aguarda espaço no semáforo
			sem <- struct{}{}

			// Tenta aceitar nova sessão
			sessionReceiver, err := client.AcceptNextSessionForQueue(ctx, queueName, nil)
			if err != nil {
				var sbErr *azservicebus.Error
				if errors.As(err, &sbErr) && sbErr.Code == azservicebus.CodeTimeout {
					log.Println("Nenhuma sessão disponível, tentando novamente em 5s...")
					time.Sleep(5 * time.Second)
					<-sem // libera o slot reservado no semáforo
					continue
				}
				<-sem
				log.Printf("Erro inesperado ao aceitar sessão: %v", err)
				continue
			}

			log.Printf("Recebendo sessão '%s'\n", sessionReceiver.SessionID())

			go func(receiver *azservicebus.SessionReceiver) {
				defer func() {
					log.Printf("Fechando receiver da sessão '%s'\n", receiver.SessionID())
					err := receiver.Close(context.Background())
					if err != nil {
						log.Printf("Erro ao fechar receiver: %v", err)
					}
					<-sem // libera o slot no semáforo ao final da goroutine
				}()

				const maxIdleTries = 3
				idleCount := 0

				for {
					innerCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
					messages, err := receiver.ReceiveMessages(innerCtx, 1, nil)
					cancel()

					if err != nil {
						log.Printf("Erro ao receber mensagens da sessão '%s': %v\n", receiver.SessionID(), err)
						break
					}

					if len(messages) == 0 {
						idleCount++
						log.Printf("Sessão '%s' ociosa (%d/%d)\n", receiver.SessionID(), idleCount, maxIdleTries)
						if idleCount >= maxIdleTries {
							log.Printf("Encerrando sessão '%s' por inatividade\n", receiver.SessionID())
							break
						}
						continue
					}

					idleCount = 0 // reseta contador de ociosidade
					for _, msg := range messages {
						if err := processMessageFromSession(ctx, receiver, msg); err != nil {
							log.Printf("Erro ao processar mensagem: %v", err)
						}
					}
				}
			}(sessionReceiver)
		}
	}
}

func processMessageFromSession(ctx context.Context, receiver *azservicebus.SessionReceiver, msg *azservicebus.ReceivedMessage) error {
	log.Printf("Mensagem recebida da sessão %s, ID %s\n", *msg.SessionID, msg.MessageID)
	return receiver.CompleteMessage(ctx, msg, nil)
}

func stringToInt(workerCount string) int {
	i, err := strconv.Atoi(workerCount)
	if err != nil {
		log.Fatalf("Erro ao converter WORKER_COUNT para int: %v", err)
	}
	return i
}
