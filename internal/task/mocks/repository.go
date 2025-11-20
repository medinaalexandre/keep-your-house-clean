package mocks

import (
	"context"
	"keep-your-house-clean/internal/domain"
	"keep-your-house-clean/internal/events"
)

type MockTaskRepository struct {
	CreateFunc                func(ctx context.Context, task *domain.Task) error
	FetchAllFunc              func(ctx context.Context, tenantID int64) ([]domain.Task, error)
	GetByIDFunc               func(ctx context.Context, id int64, tenantID int64) (*domain.Task, error)
	UpdateFunc                func(ctx context.Context, task *domain.Task) error
	DeleteFunc                func(ctx context.Context, id int64, tenantID int64) error
	GetUpcomingTasksFunc      func(ctx context.Context, tenantID int64, limit int, offset int) ([]domain.Task, error)
	GetCompletedTasksHistoryFunc func(ctx context.Context, tenantID int64, limit int) ([]domain.TaskWithUser, error)
}

func (m *MockTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, task)
	}
	return nil
}

func (m *MockTaskRepository) FetchAll(ctx context.Context, tenantID int64) ([]domain.Task, error) {
	if m.FetchAllFunc != nil {
		return m.FetchAllFunc(ctx, tenantID)
	}
	return []domain.Task{}, nil
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id int64, tenantID int64) (*domain.Task, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id, tenantID)
	}
	return nil, nil
}

func (m *MockTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, task)
	}
	return nil
}

func (m *MockTaskRepository) Delete(ctx context.Context, id int64, tenantID int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id, tenantID)
	}
	return nil
}

func (m *MockTaskRepository) GetUpcomingTasks(ctx context.Context, tenantID int64, limit int, offset int) ([]domain.Task, error) {
	if m.GetUpcomingTasksFunc != nil {
		return m.GetUpcomingTasksFunc(ctx, tenantID, limit, offset)
	}
	return []domain.Task{}, nil
}

func (m *MockTaskRepository) GetCompletedTasksHistory(ctx context.Context, tenantID int64, limit int) ([]domain.TaskWithUser, error) {
	if m.GetCompletedTasksHistoryFunc != nil {
		return m.GetCompletedTasksHistoryFunc(ctx, tenantID, limit)
	}
	return []domain.TaskWithUser{}, nil
}

type MockUserRepository struct {
	GetByIDFunc func(ctx context.Context, id int64) (*domain.User, error)
	UpdateFunc  func(ctx context.Context, user *domain.User) error
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
	return nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (m *MockUserRepository) GetByEmailAndTenant(ctx context.Context, email string, tenantID int64) (*domain.User, error) {
	return nil, nil
}

func (m *MockUserRepository) FetchAll(ctx context.Context, tenantID int64) ([]domain.User, error) {
	return nil, nil
}

func (m *MockUserRepository) GetTopUsersByPoints(ctx context.Context, tenantID int64, limit int) ([]domain.User, error) {
	return nil, nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
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
