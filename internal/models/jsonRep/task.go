package jsonrep

import "github.com/google/uuid"

type TaskRequest struct {
	Translator string `json:"translator" binding:"required"`
	Code       string `json:"code" binding:"required"`
}

type TaskJSON struct {
	ID           uuid.UUID `json:"id"`
	Status       string    `json:"status"`
	Code         string    `json:"code"`
	CompilerName string    `json:"compiler_name"`
	Result       string    `json:"result"`
}
