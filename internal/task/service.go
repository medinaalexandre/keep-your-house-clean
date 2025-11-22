package task

import (
	"context"
	"keep-your-house-clean/internal/domain"
	"keep-your-house-clean/internal/events"
	"keep-your-house-clean/internal/platform/middleware"
	"time"
)

type Service struct {
	repo      domain.TaskRepository
	dispatcher events.EventDispatcher
}

func NewService(repo domain.TaskRepository, dispatcher events.EventDispatcher) *Service {
	return &Service{
		repo:      repo,
		dispatcher: dispatcher,
	}
}

func (s *Service) CreateTask(ctx context.Context, req CreateTaskRequest) (*domain.Task, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if userID == 0 || tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	now := time.Now()
	task := &domain.Task{
		Title:          req.Title,
		Description:    req.Description,
		Points:         req.Points,
		Status:         req.Status,
		ScheduledTo:    req.ScheduledTo,
		ScheduledById:  req.ScheduledById,
		FrequencyValue: req.FrequencyValue,
		FrequencyUnit:  req.FrequencyUnit,
		Completed:      false,
		TenantID:       tenantID,
		CreatedAt:      now,
		CreatedById:    userID,
		UpdatedAt:      now,
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) GetTaskByID(ctx context.Context, id int64) (*domain.Task, error) {
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	task, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

func (s *Service) ListTasks(ctx context.Context) ([]domain.Task, error) {
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	tasks, err := s.repo.FetchAll(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Service) UpdateTask(ctx context.Context, id int64, req UpdateTaskRequest) (*domain.Task, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if userID == 0 || tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	task, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, ErrTaskNotFound
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Points != nil {
		task.Points = *req.Points
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.ScheduledTo != nil {
		task.ScheduledTo = req.ScheduledTo
	}
	if req.ScheduledById != nil {
		task.ScheduledById = req.ScheduledById
	}
	if req.FrequencyValue != nil {
		task.FrequencyValue = *req.FrequencyValue
	}
	if req.FrequencyUnit != nil {
		task.FrequencyUnit = *req.FrequencyUnit
	}
	if req.Completed != nil {
		task.Completed = *req.Completed
		if *req.Completed && req.CompletedAt != nil {
			task.CompletedById = &userID
			nextDueDate, err := task.CalculateNextDueDate(*req.CompletedAt)
			if err != nil {
				return nil, err
			}
			task.ScheduledTo = &nextDueDate
		} else if !*req.Completed {
			task.CompletedById = nil
		}
	}

	task.UpdatedAt = time.Now()
	task.UpdatedById = &userID

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) CompleteTask(ctx context.Context, id int64, req CompleteTaskRequest) (*domain.Task, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if userID == 0 || tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	task, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, ErrTaskNotFound
	}

	if task.Completed {
		return nil, ErrTaskAlreadyCompleted
	}

	completedByID := userID
	if req.CompletedById != nil {
		completedByID = *req.CompletedById
	}

	now := time.Now()
	task.Completed = true
	task.CompletedById = &completedByID
	task.UpdatedAt = now
	task.UpdatedById = &userID

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	if task.Points > 0 {
		event := events.Event{
			Type:      events.EventTypeTaskCompleted,
			Payload:   events.TaskCompletedPayload{
				CompletedBy: completedByID,
				Points:      task.Points,
			},
			Timestamp: now,
		}
		if err := s.dispatcher.Dispatch(event); err != nil {
			return nil, err
		}
	}

	if task.FrequencyValue > 0 && task.FrequencyUnit != "" {
		nextScheduledDate, err := task.CalculateNextDueDate(now)
		if err != nil {
			return nil, err
		}
		
		newTask := &domain.Task{
			Title:          task.Title,
			Description:    task.Description,
			Points:         task.Points,
			Status:         task.Status,
			ScheduledTo:    &nextScheduledDate,
			ScheduledById:  task.ScheduledById,
			FrequencyValue: task.FrequencyValue,
			FrequencyUnit:  task.FrequencyUnit,
			Completed:      false,
			CompletedById:  nil,
			TenantID:       task.TenantID,
			CreatedAt:      now,
			CreatedById:    userID,
			UpdatedAt:      now,
		}

		if err := s.repo.Create(ctx, newTask); err != nil {
			return nil, err
		}
	}

	return task, nil
}

func (s *Service) DeleteTask(ctx context.Context, id int64) error {
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if tenantID == 0 {
		return ErrUserNotAuthenticated
	}

	return s.repo.Delete(ctx, id, tenantID)
}

func (s *Service) GetUpcomingTasks(ctx context.Context, limit int, offset int) ([]domain.Task, error) {
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	return s.repo.GetUpcomingTasks(ctx, tenantID, limit, offset)
}

func (s *Service) GetCompletedTasksHistory(ctx context.Context, limit int) ([]domain.TaskWithUser, error) {
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	return s.repo.GetCompletedTasksHistory(ctx, tenantID, limit)
}

func (s *Service) GetCompletedTasksByUser(ctx context.Context, userID int64, limit int, offset int) ([]domain.TaskWithUser, error) {
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	return s.repo.GetCompletedTasksByUser(ctx, userID, tenantID, limit, offset)
}

func (s *Service) UndoCompleteTask(ctx context.Context, id int64) (*domain.Task, error) {
	userID := middleware.GetUserIDFromContext(ctx)
	tenantID := middleware.GetTenantIDFromContext(ctx)
	if userID == 0 || tenantID == 0 {
		return nil, ErrUserNotAuthenticated
	}

	task, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, ErrTaskNotFound
	}

	if !task.Completed {
		return nil, ErrTaskNotCompleted
	}

	completedByID := task.CompletedById
	if completedByID == nil {
		return nil, ErrTaskNotCompleted
	}

	now := time.Now()
	completionTime := task.UpdatedAt

	if task.FrequencyValue > 0 && task.FrequencyUnit != "" {
		createdTask, err := s.repo.FindTaskCreatedAfterCompletion(ctx, task, completionTime)
		if err != nil {
			return nil, err
		}

		if createdTask != nil {
			if createdTask.ScheduledTo != nil {
				originalScheduledDate, err := task.CalculatePreviousDueDate(*createdTask.ScheduledTo)
				if err == nil {
					task.ScheduledTo = &originalScheduledDate
				}
			}
			
			if err := s.repo.Delete(ctx, createdTask.ID, tenantID); err != nil {
				return nil, err
			}
		}
	}

	task.Completed = false
	task.CompletedById = nil
	task.UpdatedAt = now
	task.UpdatedById = &userID

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	if task.Points > 0 {
		event := events.Event{
			Type: events.EventTypeTaskUndone,
			Payload: events.TaskUndonePayload{
				CompletedBy: *completedByID,
				Points:      task.Points,
			},
			Timestamp: now,
		}
		if err := s.dispatcher.Dispatch(event); err != nil {
			return nil, err
		}
	}

	return task, nil
}
