package user

import "time"

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	TenantID int64  `json:"tenant_id"`
	Points   int    `json:"points"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type UpdateUserRequest struct {
	Name        *string    `json:"name"`
	Email       *string    `json:"email"`
	Password    *string    `json:"password"`
	Points      *int       `json:"points"`
	Role        *string    `json:"role"`
	Status      *string    `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at"`
}
