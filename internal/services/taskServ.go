package services

import (
	"codeProcessor/internal/cnfg"
	"codeProcessor/internal/models"
	"codeProcessor/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type TaskServ interface {
	GetTaskStatus(taskID uuid.UUID) (models.TaskStatus, error)
	GetTaskResult(taskID uuid.UUID) (string, error)
	AddTask(ctx context.Context, task models.Task) error
	UpdateTask(taskID uuid.UUID, status string, result string) error
	Close()
}

var (
	ErrTaskServ    = errors.New("taskServ: ")
	ErrTaskNotRead = errors.New("task not ready")
)

type taskServ struct {
	taskStorage storage.TaskStorage
	conn        *amqp.Connection
	ch          *amqp.Channel
	q           amqp.Queue
}

func NewTaskServ(taskStorage storage.TaskStorage, rabbitMqCnfg *cnfg.RabbitMQConfig) (TaskServ, error) {
	conn, err := amqp.Dial(rabbitMqCnfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrTaskServ, err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrTaskServ, err)
	}
	q, err := ch.QueueDeclare(
		"queue of tasks", // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrTaskServ, err)
	}

	return &taskServ{
		taskStorage: taskStorage,
		conn:        conn,
		ch:          ch,
		q:           q,
	}, nil
}

func (ts *taskServ) AddTask(ctx context.Context, task models.Task) error {
	if err := ts.taskStorage.AddTask(task); err != nil {
		return fmt.Errorf("%w: %w", ErrTaskServ, err)
	}

	// Сериализация в JSON
	body, err := json.Marshal(task.ToTaskJson())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTaskServ, err)
	}
	err = ts.ch.PublishWithContext(ctx,
		"",        // exchange
		ts.q.Name, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json", // меняем на JSON
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTaskServ, err)
	}
	return nil
}

func (ts *taskServ) GetTaskStatus(taskID uuid.UUID) (models.TaskStatus, error) {
	t, err := ts.taskStorage.GetTaskByID(taskID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTaskServ, err)
	}
	return t.Status(), nil
}

func (ts *taskServ) GetTaskResult(taskID uuid.UUID) (string, error) {
	t, err := ts.taskStorage.GetTaskByID(taskID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTaskServ, err)
	}
	if t.Status() != models.StatusReady {
		return "", fmt.Errorf("%w: %w", ErrTaskServ, ErrTaskNotRead)
	}
	return t.Result(), nil
}

func (ts *taskServ) UpdateTask(taskID uuid.UUID, status string, result string) error {
	task, err := ts.taskStorage.GetTaskByID(taskID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTaskServ, err)
	}
	task.SetStatus(models.TaskStatus(status))
	task.SetResult(result)
	return ts.taskStorage.UpdateTask(task)
}

func (ts *taskServ) Close() {
	ts.conn.Close()
	ts.ch.Close()
}
