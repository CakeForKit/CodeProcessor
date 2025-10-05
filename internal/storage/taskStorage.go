package storage

import (
	"codeProcessor/internal/models"
	"errors"

	"github.com/google/uuid"
)

type TaskStorage interface {
	GetTaskByID(id uuid.UUID) (models.Task, error)
	AddTask(task models.Task) error
	UpdateTask(task models.Task) error
}

func NewTaskStorage() (TaskStorage, error) {
	return &taskStorage{
		tasks: make(map[uuid.UUID]models.Task),
	}, nil
}

type taskStorage struct {
	tasks map[uuid.UUID]models.Task
}

var (
	ErrNoTask            = errors.New("no task in map")
	ErrTaskIDAlreadExist = errors.New("task id already exist")
)

func (ts *taskStorage) GetTaskByID(id uuid.UUID) (models.Task, error) {
	if t, ok := ts.tasks[id]; ok {
		return t, nil
	} else {
		return models.Task{}, ErrNoTask
	}
}

func (ts *taskStorage) AddTask(task models.Task) error {
	if _, ok := ts.tasks[task.ID()]; ok {
		return ErrTaskIDAlreadExist
	}
	// task.SetStatus(models.StatusInProcess)
	ts.tasks[task.ID()] = task
	return nil
}

func (ts *taskStorage) UpdateTask(task models.Task) error {
	ts.tasks[task.ID()] = task
	return nil
}
