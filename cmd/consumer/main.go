package main

import (
	"context"
	"log"
	"os"
	"strconv"

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
}

func stringToInt(workerCount string) int {
	i, err := strconv.Atoi(workerCount)
	if err != nil {
		log.Fatalf("Erro ao converter WORKER_COUNT para int: %v", err)
	}
	return i
}
