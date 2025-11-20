package domain

import (
	"context"
	"time"
)

type Tenant struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Domain    string     `json:"domain"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type TenantRepository interface {
	Create(ctx context.Context, tenant *Tenant) error
	GetByID(ctx context.Context, id int64) (*Tenant, error)
	GetByDomain(ctx context.Context, domain string) (*Tenant, error)
	FetchAll(ctx context.Context) ([]Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id int64) error
}
