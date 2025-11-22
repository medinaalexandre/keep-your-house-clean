package compliment

import (
	"context"
	"keep-your-house-clean/internal/domain"
	"keep-your-house-clean/internal/events"
	"keep-your-house-clean/internal/platform/middleware"
	"time"
)

type Service struct {
	repo          domain.ComplimentRepository
	userRepo      domain.UserRepository
	dispatcher    events.EventDispatcher
}

func NewService(repo domain.ComplimentRepository, userRepo domain.UserRepository, dispatcher events.EventDispatcher) *Service {
	return &Service{
		repo:       repo,
		userRepo:   userRepo,
		dispatcher: dispatcher,
	}
}

func (s *Service) CreateCompliment(ctx context.Context, req CreateComplimentRequest) (*domain.Compliment, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if userID == 0 || tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	if req.Points < 0 || req.Points > 5 {
		return nil, ErrInvalidPoints
	}

	if req.ToUserID == userID {
		return nil, ErrInvalidUser
	}

	toUser, err := s.userRepo.GetByID(ctx, req.ToUserID)
	if err != nil {
		return nil, err
	}
	if toUser == nil || toUser.TenantID != tenantID {
		return nil, ErrUserNotFound
	}

	now := time.Now()
	compliment := &domain.Compliment{
		Title:       req.Title,
		Description: req.Description,
		Points:      req.Points,
		FromUserID:  userID,
		ToUserID:    req.ToUserID,
		TenantID:    tenantID,
		CreatedAt:   now,
		CreatedById: userID,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, compliment); err != nil {
		return nil, err
	}

	if req.Points > 0 {
		event := events.Event{
			Type: events.EventTypeComplimentReceived,
			Payload: events.ComplimentReceivedPayload{
				ToUser: req.ToUserID,
				Points: req.Points,
			},
			Timestamp: now,
		}
		if err := s.dispatcher.Dispatch(event); err != nil {
			return nil, err
		}
	}

	return compliment, nil
}

func (s *Service) GetComplimentByID(ctx context.Context, id int64) (*domain.Compliment, error) {
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	compliment, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}

	if compliment == nil {
		return nil, ErrComplimentNotFound
	}

	return compliment, nil
}

func (s *Service) ListCompliments(ctx context.Context) ([]domain.Compliment, error) {
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	compliments, err := s.repo.FetchAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	return compliments, nil
}

func (s *Service) GetLastReceivedCompliment(ctx context.Context) (*domain.ComplimentWithUser, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if userID == 0 || tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	compliment, err := s.repo.GetLastReceivedByUser(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}

	return compliment, nil
}

func (s *Service) GetUserComplimentsHistory(ctx context.Context) ([]domain.ComplimentWithUser, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if userID == 0 || tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	compliments, err := s.repo.GetUserComplimentsHistory(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}

	return compliments, nil
}

func (s *Service) GetUnviewedReceivedCompliments(ctx context.Context) ([]domain.ComplimentWithUser, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if userID == 0 || tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	compliments, err := s.repo.GetUnviewedReceivedCompliments(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}

	return compliments, nil
}

func (s *Service) MarkComplimentsAsViewed(ctx context.Context, ids []int64) error {
	userID := middleware.GetUserIDFromContext(ctx)
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if userID == 0 || tenantID == 0 {
		return ErrUserNotAuthenticated
	}

	return s.repo.MarkAsViewed(ctx, ids, userID, tenantID)
}

func (s *Service) DeleteCompliment(ctx context.Context, id int64) error {
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if tenantID == 0 {
		return ErrUserNotAuthenticated
	}

	return s.repo.Delete(ctx, id, tenantID)
}

