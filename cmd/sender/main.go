package main

import (
	"context"
	"fmt"
	"log"
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

	sender, err := client.NewSender(envs.QueueName, nil)
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
