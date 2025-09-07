package services

import (
	"codeProcessor/internal/models"
	"codeProcessor/internal/storage"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type TaskServ interface {
	GetTaskStatus(taskID uuid.UUID) (models.TaskStatus, error)
	GetTaskResult(taskID uuid.UUID) (string, error)
	AddTask(task models.Task) error
}

var (
	ErrTaskServ    = errors.New("taskServ: ")
	ErrTaskNotRead = errors.New("task not ready")
)

type taskServ struct {
	taskStorage storage.TaskStorage
}

func NewTaskServ(taskStorage storage.TaskStorage) (TaskServ, error) {
	return &taskServ{
		taskStorage: taskStorage,
	}, nil
}

func (ts *taskServ) AddTask(task models.Task) error {
	if err := ts.taskStorage.AddTask(task); err != nil {
		return fmt.Errorf("%w: %v", ErrTaskServ, err)
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
