package tenant

type CreateTenantRequest struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Status string `json:"status"`
}

type UpdateTenantRequest struct {
	Name   *string `json:"name"`
	Domain *string `json:"domain"`
	Status *string `json:"status"`
}
