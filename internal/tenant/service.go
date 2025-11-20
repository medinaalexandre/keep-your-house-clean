package tenant

import (
	"context"
	"errors"
	"keep-your-house-clean/internal/domain"
	"time"
)

type Service struct {
	repo domain.TenantRepository
}

func NewService(repo domain.TenantRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTenant(ctx context.Context, req CreateTenantRequest) (*domain.Tenant, error) {
	existingTenant, err := s.repo.GetByDomain(ctx, req.Domain)
	if err != nil {
		return nil, err
	}
	if existingTenant != nil {
		return nil, errors.New("domain already exists")
	}

	now := time.Now()
	tenant := &domain.Tenant{
		Name:      req.Name,
		Domain:    req.Domain,
		Status:    getStatusOrDefault(req.Status),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

func (s *Service) GetTenantByID(ctx context.Context, id int64) (*domain.Tenant, error) {
	tenant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if tenant == nil {
		return nil, errors.New("tenant not found")
	}

	return tenant, nil
}

func (s *Service) GetTenantByDomain(ctx context.Context, domain string) (*domain.Tenant, error) {
	tenant, err := s.repo.GetByDomain(ctx, domain)
	if err != nil {
		return nil, err
	}

	if tenant == nil {
		return nil, errors.New("tenant not found")
	}

	return tenant, nil
}

func (s *Service) ListTenants(ctx context.Context) ([]domain.Tenant, error) {
	tenants, err := s.repo.FetchAll(ctx)
	if err != nil {
		return nil, err
	}

	return tenants, nil
}

func (s *Service) UpdateTenant(ctx context.Context, id int64, req UpdateTenantRequest) (*domain.Tenant, error) {
	tenant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if tenant == nil {
		return nil, errors.New("tenant not found")
	}

	if req.Name != nil {
		tenant.Name = *req.Name
	}
	if req.Domain != nil {
		existingTenant, err := s.repo.GetByDomain(ctx, *req.Domain)
		if err != nil {
			return nil, err
		}
		if existingTenant != nil && existingTenant.ID != id {
			return nil, errors.New("domain already exists")
		}
		tenant.Domain = *req.Domain
	}
	if req.Status != nil {
		tenant.Status = *req.Status
	}

	tenant.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

func (s *Service) DeleteTenant(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func getStatusOrDefault(status string) string {
	if status == "" {
		return "active"
	}
	return status
}
