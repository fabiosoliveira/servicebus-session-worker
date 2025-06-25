package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/fabiosoliveira/servicebus-session-worker/internal/envs"
)

func main() {

	client, err := azservicebus.NewClientFromConnectionString(envs.ConnectionString, nil)
	if err != nil {
		log.Fatalf("Erro ao criar o client: %v", err)
	}
	defer client.Close(context.Background())

	sem := make(chan struct{}, envs.WorkerCount)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown setup
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		log.Println("Sinal de interrupção recebido. Encerrando com graceful shutdown...")
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Contexto cancelado. Encerrando loop principal.")
			return
		default:
			sem <- struct{}{}

			sessionReceiver, err := client.AcceptNextSessionForQueue(ctx, envs.QueueName, nil)
			if err != nil {
				var sbErr *azservicebus.Error
				if errors.As(err, &sbErr) && sbErr.Code == azservicebus.CodeTimeout {
					log.Println("Nenhuma sessão disponível, tentando novamente em 5s...")
					time.Sleep(5 * time.Second)
					<-sem
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
					<-sem
				}()

				const maxIdleTries = 3
				idleCount := 0

				for {
					select {
					case <-ctx.Done():
						log.Printf("Contexto cancelado. Encerrando receiver da sessão '%s'\n", receiver.SessionID())
						return
					default:
						innerCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
						messages, err := receiver.ReceiveMessages(innerCtx, 1, nil)
						cancel()

						if err != nil {
							log.Printf("Erro ao receber mensagens da sessão '%s': %v\n", receiver.SessionID(), err)
							break
						}

						if len(messages) == 0 {
							idleCount++
							log.Printf("Sessão '%s' ociosa (%d/3)\n", receiver.SessionID(), idleCount)
							if idleCount >= maxIdleTries {
								log.Printf("Encerrando sessão '%s' por inatividade\n", receiver.SessionID())
								break
							}
							continue
						}

						idleCount = 0
						for _, msg := range messages {
							if err := processMessageFromSession(ctx, receiver, msg); err != nil {
								log.Printf("Erro ao processar mensagem: %v", err)
							}
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
