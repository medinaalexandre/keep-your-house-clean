package events

import (
	"context"
	"time"
)

type EventType string

const (
	EventTypeTaskCompleted EventType = "task.completed"
	EventTypeTaskUndone   EventType = "task.undone"
)

type Event struct {
	Type      EventType
	Payload   interface{}
	Timestamp time.Time
}

type TaskCompletedPayload struct {
	CompletedBy int64
	Points      int
}

type TaskUndonePayload struct {
	CompletedBy int64
	Points      int
}

type EventHandler func(ctx context.Context, event Event) error

