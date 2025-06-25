package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func main() {
	connectionString := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	queueName := os.Getenv("SERVICEBUS_QUEUE_NAME")
	workerCount := os.Getenv("WORKER_COUNT")

	if connectionString == "" {
		log.Fatal("A variável de ambiente SERVICEBUS_CONNECTION_STRING não foi definida.")
	}
	if queueName == "" {
		log.Fatal("A variável de ambiente SERVICEBUS_QUEUE_NAME não foi definida.")
	}
	if workerCount == "" {
		workerCount = "100"
	}

	client, err := azservicebus.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		log.Fatalf("Erro ao criar o client: %v", err)
	}
	defer client.Close(context.Background())

	sender, err := client.NewSender(queueName, nil)
	if err != nil {
		log.Fatalf("Erro ao criar sender: %v", err)
	}
	defer sender.Close(context.Background())

	// Sessões de exemplo
	sessionIDs := make([]string, 0, 100)
	for i := 1; i <= 100; i++ {
		sessionIDs = append(sessionIDs, fmt.Sprintf("cliente-%d", i))
	}

	for i := 1; i <= 10; i++ {
		for _, sessionID := range sessionIDs {
			msg := &azservicebus.Message{
				Body:      []byte(fmt.Sprintf("Mensagem %d para sessão %s", i, sessionID)),
				SessionID: &sessionID,
			}

			err := sender.SendMessage(context.Background(), msg, nil)
			if err != nil {
				log.Printf("Erro ao enviar mensagem para sessão %s: %v", sessionID, err)
			} else {
				log.Printf("Mensagem enviada para sessão %s", sessionID)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
