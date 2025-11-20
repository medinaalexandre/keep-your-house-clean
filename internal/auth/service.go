package auth

import (
	"context"
	"keep-your-house-clean/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo   domain.UserRepository
	tenantRepo domain.TenantRepository
	jwtSecret  string
}

func NewService(userRepo domain.UserRepository, tenantRepo domain.TenantRepository, jwtSecret string) *Service {
	return &Service{
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
		jwtSecret:  jwtSecret,
	}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if user.Status != "active" {
		return nil, ErrUserInactive
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	token, err := s.generateToken(user.ID, user.TenantID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:    token,
		UserID:   user.ID,
		TenantID: user.TenantID,
		Email:    user.Email,
		Name:     user.Name,
	}, nil
}

func (s *Service) generateToken(userID, tenantID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"tenant_id": tenantID,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*LoginResponse, error) {
	existingTenant, err := s.tenantRepo.GetByDomain(ctx, req.TenantDomain)
	if err != nil {
		return nil, err
	}
	if existingTenant != nil {
		return nil, ErrDomainExists
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailExists
	}

	now := time.Now()
	tenant := &domain.Tenant{
		Name:      req.TenantName,
		Domain:    req.TenantDomain,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, err
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, ErrPasswordHashFailed
	}

	user := &domain.User{
		Name:       req.UserName,
		Email:      req.Email,
		Password:   hashedPassword,
		TenantID:   tenant.ID,
		Points:     0,
		Role:       "admin",
		Status:     "active",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := s.generateToken(user.ID, tenant.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:    token,
		UserID:   user.ID,
		TenantID: tenant.ID,
		Email:    user.Email,
		Name:     user.Name,
	}, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
