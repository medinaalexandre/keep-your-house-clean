package mocks

import (
	"context"
	"keep-your-house-clean/internal/domain"
	"keep-your-house-clean/internal/events"
)

type MockComplimentRepository struct {
	CreateFunc                    func(ctx context.Context, compliment *domain.Compliment) error
	GetByIDFunc                   func(ctx context.Context, id int64, tenantID int64) (*domain.Compliment, error)
	FetchAllFunc                  func(ctx context.Context, tenantID int64) ([]domain.Compliment, error)
	GetLastReceivedByUserFunc     func(ctx context.Context, userID int64, tenantID int64) (*domain.ComplimentWithUser, error)
	GetUserComplimentsHistoryFunc func(ctx context.Context, userID int64, tenantID int64) ([]domain.ComplimentWithUser, error)
	DeleteFunc                    func(ctx context.Context, id int64, tenantID int64) error
}

func (m *MockComplimentRepository) Create(ctx context.Context, compliment *domain.Compliment) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, compliment)
	}
	return nil
}

func (m *MockComplimentRepository) GetByID(ctx context.Context, id int64, tenantID int64) (*domain.Compliment, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id, tenantID)
	}
	return nil, nil
}

func (m *MockComplimentRepository) FetchAll(ctx context.Context, tenantID int64) ([]domain.Compliment, error) {
	if m.FetchAllFunc != nil {
		return m.FetchAllFunc(ctx, tenantID)
	}
	return []domain.Compliment{}, nil
}

func (m *MockComplimentRepository) GetLastReceivedByUser(ctx context.Context, userID int64, tenantID int64) (*domain.ComplimentWithUser, error) {
	if m.GetLastReceivedByUserFunc != nil {
		return m.GetLastReceivedByUserFunc(ctx, userID, tenantID)
	}
	return nil, nil
}

func (m *MockComplimentRepository) GetUserComplimentsHistory(ctx context.Context, userID int64, tenantID int64) ([]domain.ComplimentWithUser, error) {
	if m.GetUserComplimentsHistoryFunc != nil {
		return m.GetUserComplimentsHistoryFunc(ctx, userID, tenantID)
	}
	return []domain.ComplimentWithUser{}, nil
}

func (m *MockComplimentRepository) Delete(ctx context.Context, id int64, tenantID int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id, tenantID)
	}
	return nil
}

type MockUserRepository struct {
	GetByIDFunc func(ctx context.Context, id int64) (*domain.User, error)
	UpdateFunc  func(ctx context.Context, user *domain.User) error
	CreateFunc  func(ctx context.Context, user *domain.User) error
	GetByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
	GetByEmailAndTenantFunc func(ctx context.Context, email string, tenantID int64) (*domain.User, error)
	FetchAllFunc func(ctx context.Context, tenantID int64) ([]domain.User, error)
	GetTopUsersByPointsFunc func(ctx context.Context, tenantID int64, limit int) ([]domain.User, error)
	DeleteFunc func(ctx context.Context, id int64) error
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return &domain.User{ID: id, Points: 0}, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) GetByEmailAndTenant(ctx context.Context, email string, tenantID int64) (*domain.User, error) {
	if m.GetByEmailAndTenantFunc != nil {
		return m.GetByEmailAndTenantFunc(ctx, email, tenantID)
	}
	return nil, nil
}

func (m *MockUserRepository) FetchAll(ctx context.Context, tenantID int64) ([]domain.User, error) {
	if m.FetchAllFunc != nil {
		return m.FetchAllFunc(ctx, tenantID)
	}
	return nil, nil
}

func (m *MockUserRepository) GetTopUsersByPoints(ctx context.Context, tenantID int64, limit int) ([]domain.User, error) {
	if m.GetTopUsersByPointsFunc != nil {
		return m.GetTopUsersByPointsFunc(ctx, tenantID, limit)
	}
	return nil, nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

type MockDispatcher struct {
	DispatchFunc       func(event events.Event) error
	RegisterHandlerFunc func(eventType events.EventType, handler events.EventHandler)
	StartFunc          func()
	StopFunc           func()
}

func (m *MockDispatcher) Dispatch(event events.Event) error {
	if m.DispatchFunc != nil {
		return m.DispatchFunc(event)
	}
	return nil
}

func (m *MockDispatcher) RegisterHandler(eventType events.EventType, handler events.EventHandler) {
	if m.RegisterHandlerFunc != nil {
		m.RegisterHandlerFunc(eventType, handler)
	}
}

func (m *MockDispatcher) Start() {
	if m.StartFunc != nil {
		m.StartFunc()
	}
}

func (m *MockDispatcher) Stop() {
	if m.StopFunc != nil {
		m.StopFunc()
	}
}

