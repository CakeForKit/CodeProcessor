package models

import (
	jsonrep "codeProcessor/internal/models/jsonRep"
	"errors"
	"strings"

	"github.com/google/uuid"
)

const (
	StatusInProcess = "in_progress"
	StatusReady     = "ready"
)

type TaskStatus string

type Task struct {
	id           uuid.UUID
	status       TaskStatus
	code         string
	compilerName string
	result       string
}

var (
	ErrEmptyCompilerName = errors.New("compiler name cannot be empty")
	ErrValidate          = errors.New("not valid datas")
)

func NewTask(
	id uuid.UUID,
	status TaskStatus,
	code string,
	compilerName string,
	result string,
) (Task, error) {
	t := Task{
		id:           id,
		status:       status,
		code:         code,
		compilerName: compilerName,
		result:       result,
	}
	if !t.Validate() {
		return Task{}, ErrValidate
	}
	return t, nil
}

func (t *Task) Validate() bool {
	if !(t.status == StatusInProcess || t.status == StatusReady) ||
		strings.TrimSpace(t.code) == "" ||
		strings.TrimSpace(t.compilerName) == "" {
		return false
	}
	return true
}

func (t *Task) ToTaskJson() jsonrep.TaskJSON {
	return jsonrep.TaskJSON{
		ID:           t.id,
		Status:       string(t.status),
		Code:         t.code,
		CompilerName: t.compilerName,
		Result:       t.result,
	}
}

func (t *Task) ID() uuid.UUID {
	return t.id
}

func (t *Task) Status() TaskStatus {
	return t.status
}

func (t *Task) Code() string {
	return t.code
}

func (t *Task) CompilerName() string {
	return t.compilerName
}

func (t *Task) Result() string {
	return t.result
}

func (t *Task) SetID(id uuid.UUID) {
	t.id = id
}

func (t *Task) SetStatus(s TaskStatus) {
	t.status = s
}

func (t *Task) SetCode(code string) {
	t.code = code
}

func (t *Task) SetCompilerNameWithValidation(compilerName string) error {
	if strings.TrimSpace(compilerName) == "" {
		return ErrEmptyCompilerName
	}
	t.compilerName = compilerName
	return nil
}

func (t *Task) SetResult(r string) {
	t.result = r
}
