package domain

import (
	"context"
	"errors"
	"time"
)

type FrequencyUnit string

const (
	UnitDays   FrequencyUnit = "days"
	UnitWeeks  FrequencyUnit = "weeks"
	UnitMonths FrequencyUnit = "months"
)

type Task struct {
	ID             int64          `json:"id"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Points         int            `json:"points"`
	Status         string         `json:"status"`
	ScheduledTo    *time.Time     `json:"scheduled_to"`
	ScheduledById  *int64         `json:"scheduled_by_id"`
	FrequencyValue int            `json:"frequency_value"`
	FrequencyUnit  FrequencyUnit  `json:"frequency_unit"`
	Completed      bool           `json:"completed"`
	CompletedById  *int64         `json:"completed_by_id"`
	TenantID       int64          `json:"tenant_id"`
	CreatedAt      time.Time      `json:"created_at"`
	CreatedById    int64          `json:"created_by_id"`
	UpdatedAt      time.Time      `json:"updated_at"`
	UpdatedById    *int64         `json:"updated_by_id"`
	DeletedAt      *time.Time     `json:"deleted_at"`
}

func (t *Task) CalculateNextDueDate(completionDate time.Time) (time.Time, error) {
	if t.FrequencyUnit == "" || t.FrequencyValue <= 0 {
		return time.Time{}, errors.New("frequency unit or frequency value not defined")
	}

	switch t.FrequencyUnit {
	case UnitDays:
		return completionDate.AddDate(0, 0, t.FrequencyValue), nil
	case UnitWeeks:
		return completionDate.AddDate(0, 0, t.FrequencyValue*7), nil
	case UnitMonths:
		return completionDate.AddDate(0, t.FrequencyValue, 0), nil
	default:
		return time.Time{}, errors.New("invalid frequency unit")
	}
}

func (t *Task) CalculatePreviousDueDate(nextDueDate time.Time) (time.Time, error) {
	if t.FrequencyUnit == "" || t.FrequencyValue <= 0 {
		return time.Time{}, errors.New("frequency unit or frequency value not defined")
	}

	switch t.FrequencyUnit {
	case UnitDays:
		return nextDueDate.AddDate(0, 0, -t.FrequencyValue), nil
	case UnitWeeks:
		return nextDueDate.AddDate(0, 0, -t.FrequencyValue*7), nil
	case UnitMonths:
		return nextDueDate.AddDate(0, -t.FrequencyValue, 0), nil
	default:
		return time.Time{}, errors.New("invalid frequency unit")
	}
}

type TaskWithUser struct {
	Task
	CompletedByName *string `json:"completed_by_name"`
}

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	FetchAll(ctx context.Context, tenantID int64) ([]Task, error)
	GetByID(ctx context.Context, id int64, tenantID int64) (*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id int64, tenantID int64) error
	GetUpcomingTasks(ctx context.Context, tenantID int64, limit int, offset int) ([]Task, error)
	GetCompletedTasksHistory(ctx context.Context, tenantID int64, limit int) ([]TaskWithUser, error)
	GetCompletedTasksByUser(ctx context.Context, userID int64, tenantID int64, limit int, offset int) ([]TaskWithUser, error)
	FindTaskCreatedAfterCompletion(ctx context.Context, originalTask *Task, completionTime time.Time) (*Task, error)
}