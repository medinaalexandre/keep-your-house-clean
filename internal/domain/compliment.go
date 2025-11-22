package domain

import (
	"context"
	"time"
)

type Compliment struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Points      int        `json:"points"`
	FromUserID  int64      `json:"from_user_id"`
	ToUserID    int64      `json:"to_user_id"`
	TenantID    int64      `json:"tenant_id"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedById int64      `json:"created_by_id"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UpdatedById *int64     `json:"updated_by_id"`
	DeletedAt   *time.Time `json:"deleted_at"`
	ViewedAt    *time.Time `json:"viewed_at"`
}

type ComplimentWithUser struct {
	Compliment
	FromUserName *string `json:"from_user_name"`
}

type ComplimentRepository interface {
	Create(ctx context.Context, compliment *Compliment) error
	GetByID(ctx context.Context, id int64, tenantID int64) (*Compliment, error)
	FetchAll(ctx context.Context, tenantID int64) ([]Compliment, error)
	GetLastReceivedByUser(ctx context.Context, userID int64, tenantID int64) (*ComplimentWithUser, error)
	GetUserComplimentsHistory(ctx context.Context, userID int64, tenantID int64) ([]ComplimentWithUser, error)
	GetUnviewedReceivedCompliments(ctx context.Context, userID int64, tenantID int64) ([]ComplimentWithUser, error)
	MarkAsViewed(ctx context.Context, ids []int64, userID int64, tenantID int64) error
	Delete(ctx context.Context, id int64, tenantID int64) error
}

