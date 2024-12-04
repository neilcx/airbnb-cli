package main

import (
	"airbnb-cli/kafka"
	"airbnb-cli/tasks"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	// airbnb-cli inputs with default values
	command := flag.String("command", "", "start producer or consumer")
	data := flag.String("data", "tasks.json", "path to tasks JSON (for producer)")
	kafkaBroker := flag.String("kafka", "localhost:9092", "Kafka broker address")
	topic := flag.String("topic", "airbnb-topic", "Kafka topic")
	workers := flag.Int("workers", 10, "number of workers for consumer")
	flag.Parse()

	// read tasks data
	// since we are running in a single local machine, will hanlde data read first
	// will need to modify logic if

	switch *command {
	case "producer":
		var taskData tasks.TaskData
		fileData, err := os.ReadFile(*data)
		if err != nil {
			log.Fatalf("Error reading tasks file: %v\n", err)
		}

		if err := json.Unmarshal(fileData, &taskData); err != nil {
			log.Fatalf("Error unmarshalling tasks: %v\n", err)
		}

		producer, err := kafka.NewKafkaProducer(*kafkaBroker, *topic)
		if err != nil {
			log.Fatalf("Error creating Kafka producer: %v\n", err)
		}
		defer producer.Close()

		for _, task := range taskData.Tasks {
			err := producer.ProcessTask(task)
			if err != nil {
				log.Printf("Error processing task '%s': %v\n", task.Name, err)
			}
		}
		fmt.Println("All tasks processed.")

	case "consumer":

		consumer, err := kafka.NewKafkaConsumer(*kafkaBroker, "airbnb-consumer-group", *topic)
		if err != nil {
			log.Fatalf("Error creating Kafka consumer: %v\n", err)
		}
		defer consumer.consumer.Close()

		var wg sync.WaitGroup
		consumer.StartConsuming(*workers, &wg)
		wg.Wait()

	default:
		fmt.Println("Unknown command. Use 'producer' or 'consumer'.")
	}
}
