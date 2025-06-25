package main

import (
	"fmt"
	"os"
)

// Efinir var window PowerShell
// $env:MY_VARIABLE = "meu_valor"

func main() {
	connectionString := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	queueName := os.Getenv("SERVICEBUS_QUEUE_NAME")
	workerCount := os.Getenv("WORKER_COUNT")

	fmt.Println("connectionString", connectionString)
	fmt.Println("queueName", queueName)
	fmt.Println("workerCount", workerCount)
}
