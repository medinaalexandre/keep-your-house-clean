package auth

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	TenantName string `json:"tenant_name"`
	TenantDomain string `json:"tenant_domain"`
	UserName   string `json:"user_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserID   int64  `json:"user_id"`
	TenantID int64  `json:"tenant_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}
