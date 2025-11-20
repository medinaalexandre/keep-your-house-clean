package user

import (
	"context"
	"errors"
	"keep-your-house-clean/internal/auth"
	"keep-your-house-clean/internal/domain"
	"time"
)

type Service struct {
	repo domain.UserRepository
}

func NewService(repo domain.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*domain.User, error) {
	existingUser, err := s.repo.GetByEmailAndTenant(ctx, req.Email, req.TenantID)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists for this tenant")
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	now := time.Now()
	user := &domain.User{
		Name:       req.Name,
		Email:      req.Email,
		Password:   hashedPassword,
		TenantID:   req.TenantID,
		Points:     req.Points,
		Role:       getRoleOrDefault(req.Role),
		Status:     getStatusOrDefault(req.Status),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Password = ""
	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *Service) GetUserByEmailAndTenant(ctx context.Context, email string, tenantID int64) (*domain.User, error) {
	user, err := s.repo.GetByEmailAndTenant(ctx, email, tenantID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *Service) ListUsers(ctx context.Context, tenantID int64) ([]domain.User, error) {
	users, err := s.repo.FetchAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (s *Service) GetTopUsersByPoints(ctx context.Context, tenantID int64, limit int) ([]domain.User, error) {
	users, err := s.repo.GetTopUsersByPoints(ctx, tenantID, limit)
	if err != nil {
		return nil, err
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (s *Service) UpdateUser(ctx context.Context, id int64, req UpdateUserRequest) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		existingUser, err := s.repo.GetByEmailAndTenant(ctx, *req.Email, user.TenantID)
		if err != nil {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already exists for this tenant")
		}
		user.Email = *req.Email
	}
	if req.Password != nil {
		hashedPassword, err := auth.HashPassword(*req.Password)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		user.Password = hashedPassword
	}
	if req.Points != nil {
		user.Points = *req.Points
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.LastLoginAt != nil {
		user.LastLoginAt = req.LastLoginAt
	}

	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func getRoleOrDefault(role string) string {
	if role == "" {
		return "user"
	}
	return role
}

func getStatusOrDefault(status string) string {
	if status == "" {
		return "active"
	}
	return status
}
