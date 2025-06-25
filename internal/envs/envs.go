package envs

import (
	"log"
	"os"
	"strconv"
)

var (
	ConnectionString string
	QueueName        string
	WorkerCount      int
)

func init() {

	ConnectionString = os.Getenv("SERVICEBUS_CONNECTION_STRING")
	QueueName = os.Getenv("SERVICEBUS_QUEUE_NAME")
	WorkerCount = 100

	if ConnectionString == "" {
		log.Fatal("A variável de ambiente SERVICEBUS_CONNECTION_STRING não foi definida.")
	}
	if QueueName == "" {
		log.Fatal("A variável de ambiente SERVICEBUS_QUEUE_NAME não foi definida.")
	}
	if os.Getenv("WORKER_COUNT") != "" {
		WorkerCount = stringToInt(os.Getenv("WORKER_COUNT"))
	}
}

func stringToInt(workerCount string) int {
	i, err := strconv.Atoi(workerCount)
	if err != nil {
		log.Fatalf("Erro ao converter WORKER_COUNT para int: %v", err)
	}
	return i
}
