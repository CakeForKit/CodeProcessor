package main

import (
	"bytes"
	"codeProcessor/internal/cnfg"
	"codeProcessor/internal/models"
	jsonrep "codeProcessor/internal/models/jsonRep"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// Config
	appCnfg, err := cnfg.LoadAppConfig("./configs/", "app", "yaml")
	if err != nil {
		panic(fmt.Errorf("LoadAppConfig: %v", err))
	}
	rabbitMqCnfg, err := cnfg.LoadRabbitMQConfig("./configs/", "rabbitmq", "env")
	if err != nil {
		panic(fmt.Errorf("LoadRabbitMQConfig: %v", err))
	}

	// RabbitMQ
	fmt.Printf("%s\n\n", rabbitMqCnfg.RabbitMQURL)

	conn, err := amqp.Dial(rabbitMqCnfg.RabbitMQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"queue of tasks", // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message with content type: %s", d.ContentType)

			// Проверяем content type
			if d.ContentType != "application/json" {
				log.Printf("Unexpected content type: %s", d.ContentType)
				d.Nack(false, false) // отбрасываем сообщение
				continue
			}

			// Десериализация JSON
			var taskJSON jsonrep.TaskJSON
			if err := json.Unmarshal(d.Body, &taskJSON); err != nil {
				log.Printf("Error decoding JSON: %v", err)
				log.Printf("Raw message: %s", string(d.Body))
				d.Nack(false, false) // отбрасываем некорректное сообщение
				continue
			}

			log.Printf("Successfully received task: ID=%s, Compiler=%s, CodeLength=%d",
				taskJSON.ID, taskJSON.CompilerName, len(taskJSON.Code))

			// Обработка задачи
			if err := processTask(&taskJSON, *appCnfg); err != nil {
				log.Printf("Error processing task %s: %v", taskJSON.ID, err)
				d.Nack(false, true) // повторная попытка
			} else {
				d.Ack(false) // Подтверждаем успешную обработку
				log.Printf("Task %s processed successfully", taskJSON.ID)
				log.Printf("%v", taskJSON)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func processTask(task *jsonrep.TaskJSON, appCnfg cnfg.AppConfig) error {
	log.Printf("Processing task %s with compiler %s", task.ID, task.CompilerName)
	task.Result = fmt.Sprintf("(%s) here is result! ", task.CompilerName)
	task.Status = models.StatusReady

	jsonData, _ := json.Marshal(*task)
	resp, err := http.Post(
		fmt.Sprintf("http://code_processor_app:%d/commit", appCnfg.Port),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		log.Printf("resp: %d, %v", resp.StatusCode, resp.Body)
		return errors.New("resp.StatusCode != 200 ")
	}
	return nil
}
