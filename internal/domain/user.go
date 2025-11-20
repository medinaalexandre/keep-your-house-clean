package domain

import (
	"context"
	"time"
)

type User struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Password   string     `json:"-" db:"password"`
	TenantID   int64      `json:"tenant_id"`
	Points     int        `json:"points"`
	Role       string     `json:"role"`
	Status     string     `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByEmailAndTenant(ctx context.Context, email string, tenantID int64) (*User, error)
	FetchAll(ctx context.Context, tenantID int64) ([]User, error)
	GetTopUsersByPoints(ctx context.Context, tenantID int64, limit int) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}
