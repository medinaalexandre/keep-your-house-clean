package task

import (
	"keep-your-house-clean/internal/domain"
	"time"
)

type CreateTaskRequest struct {
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	Points         int                    `json:"points"`
	Status         string                 `json:"status"`
	ScheduledTo    *time.Time             `json:"scheduled_to"`
	ScheduledById  *int64                 `json:"scheduled_by_id"`
	FrequencyValue int                    `json:"frequency_value"`
	FrequencyUnit  domain.FrequencyUnit   `json:"frequency_unit"`
}

type UpdateTaskRequest struct {
	Title          *string                `json:"title"`
	Description    *string                `json:"description"`
	Points         *int                   `json:"points"`
	Status         *string                `json:"status"`
	ScheduledTo    *time.Time             `json:"scheduled_to"`
	ScheduledById  *int64                 `json:"scheduled_by_id"`
	FrequencyValue *int                   `json:"frequency_value"`
	FrequencyUnit  *domain.FrequencyUnit  `json:"frequency_unit"`
	Completed      *bool                  `json:"completed"`
	CompletedAt    *time.Time             `json:"completed_at"`
}

type CompleteTaskRequest struct {
	CompletedById *int64 `json:"completed_by_id"`
}
