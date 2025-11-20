package task

import (
	"context"
	"errors"
	"keep-your-house-clean/internal/domain"
	"keep-your-house-clean/internal/platform/middleware"
	"keep-your-house-clean/internal/task/mocks"
	"testing"
	"time"
)

func createContextWithUserID(userID int64) context.Context {
	return middleware.SetUserIDInContext(context.Background(), userID)
}

func TestNewService(t *testing.T) {
	repo := &mocks.MockTaskRepository{}
	dispatcher := &mocks.MockDispatcher{}
	service := NewService(repo, dispatcher)

	if service == nil {
		t.Fatal("NewService retornou nil")
	}

	if service.repo != repo {
		t.Error("Repository não foi atribuído corretamente")
	}

	if service.dispatcher != dispatcher {
		t.Error("Dispatcher não foi atribuído corretamente")
	}
}

func TestService_CreateTask(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		req           CreateTaskRequest
		mockSetup     func(*mocks.MockTaskRepository)
		expectedError error
		validateTask  func(*testing.T, *domain.Task)
	}{
		{
			name: "sucesso ao criar tarefa",
			ctx:  createContextWithUserID(1),
			req: CreateTaskRequest{
				Title:          "Limpar casa",
				Description:    "Limpar todos os cômodos",
				Points:         10,
				Status:         "pending",
				FrequencyValue: 7,
				FrequencyUnit: domain.UnitDays,
			},
			mockSetup: func(m *mocks.MockTaskRepository) {
				m.CreateFunc = func(ctx context.Context, task *domain.Task) error {
					task.ID = 1
					return nil
				}
			},
			validateTask: func(t *testing.T, task *domain.Task) {
				if task.ID != 1 {
					t.Errorf("ID esperado 1, obtido %d", task.ID)
				}
				if task.Title != "Limpar casa" {
					t.Errorf("Title esperado 'Limpar casa', obtido '%s'", task.Title)
				}
				if task.CreatedById != 1 {
					t.Errorf("CreatedById esperado 1, obtido %d", task.CreatedById)
				}
				if task.Completed != false {
					t.Error("Completed deve ser false por padrão")
				}
			},
		},
		{
			name: "erro quando usuário não autenticado",
			ctx:  context.Background(),
			req: CreateTaskRequest{
				Title:          "Limpar casa",
				Description:    "Limpar todos os cômodos",
				Points:         10,
				Status:         "pending",
				FrequencyValue: 7,
				FrequencyUnit: domain.UnitDays,
			},
			expectedError: ErrUserNotAuthenticated,
		},
		{
			name: "erro ao salvar no repositório",
			ctx:  createContextWithUserID(1),
			req: CreateTaskRequest{
				Title:          "Limpar casa",
				Description:    "Limpar todos os cômodos",
				Points:         10,
				Status:         "pending",
				FrequencyValue: 7,
				FrequencyUnit: domain.UnitDays,
			},
			mockSetup: func(m *mocks.MockTaskRepository) {
				m.CreateFunc = func(ctx context.Context, task *domain.Task) error {
					return errors.New("database error")
				}
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockTaskRepository{}
			mockDispatcher := &mocks.MockDispatcher{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			service := NewService(mockRepo, mockDispatcher)
			ctx := tt.ctx
			if userID := middleware.GetUserIDFromContext(ctx); userID > 0 {
				ctx = middleware.SetTenantIDInContext(ctx, 1)
			}
			task, err := service.CreateTask(ctx, tt.req)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("erro esperado '%v', mas nenhum erro foi retornado", tt.expectedError)
					return
				}
				if !errors.Is(err, tt.expectedError) && err.Error() != tt.expectedError.Error() {
					t.Errorf("erro esperado '%v', obtido '%v'", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("erro inesperado: %v", err)
				return
			}

			if task == nil {
				t.Fatal("task não deveria ser nil")
			}

			if tt.validateTask != nil {
				tt.validateTask(t, task)
			}
		})
	}
}

func TestService_GetTaskByID(t *testing.T) {
	tests := []struct {
		name          string
		id            int64
		mockSetup     func(*mocks.MockTaskRepository)
		expectedError error
		validateTask  func(*testing.T, *domain.Task)
	}{
		{
			name: "sucesso ao buscar tarefa",
			id:   1,
			mockSetup: func(m *mocks.MockTaskRepository) {
				m.GetByIDFunc = func(ctx context.Context, id int64, tenantID int64) (*domain.Task, error) {
					return &domain.Task{
						ID:          1,
						Title:       "Limpar casa",
						Description: "Limpar todos os cômodos",
						Points:      10,
					}, nil
				}
			},
			validateTask: func(t *testing.T, task *domain.Task) {
				if task.ID != 1 {
					t.Errorf("ID esperado 1, obtido %d", task.ID)
				}
				if task.Title != "Limpar casa" {
					t.Errorf("Title esperado 'Limpar casa', obtido '%s'", task.Title)
				}
			},
		},
		{
			name: "erro quando tarefa não encontrada",
			id:   999,
			mockSetup: func(m *mocks.MockTaskRepository) {
				m.GetByIDFunc = func(ctx context.Context, id int64, tenantID int64) (*domain.Task, error) {
					return nil, nil
				}
			},
			expectedError: ErrTaskNotFound,
		},
		{
			name: "erro do repositório",
			id:   1,
			mockSetup: func(m *mocks.MockTaskRepository) {
				m.GetByIDFunc = func(ctx context.Context, id int64, tenantID int64) (*domain.Task, error) {
					return nil, errors.New("database error")
				}
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockTaskRepository{}
			mockDispatcher := &mocks.MockDispatcher{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			service := NewService(mockRepo, mockDispatcher)
			ctx := createContextWithUserID(1)
			ctx = middleware.SetTenantIDInContext(ctx, 1)
			task, err := service.GetTaskByID(ctx, tt.id)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("erro esperado '%v', mas nenhum erro foi retornado", tt.expectedError)
					return
				}
				if !errors.Is(err, tt.expectedError) && err.Error() != tt.expectedError.Error() {
					t.Errorf("erro esperado '%v', obtido '%v'", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("erro inesperado: %v", err)
				return
			}

			if task == nil {
				t.Fatal("task não deveria ser nil")
			}

			if tt.validateTask != nil {
				tt.validateTask(t, task)
			}
		})
	}
}

func TestService_ListTasks(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockTaskRepository)
		expectedError  error
		expectedLength int
	}{
		{
			name: "sucesso ao listar tarefas",
			mockSetup: func(m *mocks.MockTaskRepository) {
				m.FetchAllFunc = func(ctx context.Context, tenantID int64) ([]domain.Task, error) {
					return []domain.Task{
						{ID: 1, Title: "Tarefa 1"},
						{ID: 2, Title: "Tarefa 2"},
					}, nil
				}
			},
			expectedLength: 2,
		},
		{
			name: "retorna lista vazia",
			mockSetup: func(m *mocks.MockTaskRepository) {
				m.FetchAllFunc = func(ctx context.Context, tenantID int64) ([]domain.Task, error) {
					return []domain.Task{}, nil
				}
			},
			expectedLength: 0,
		},
		{
			name: "erro do repositório",
			mockSetup: func(m *mocks.MockTaskRepository) {
				m.FetchAllFunc = func(ctx context.Context, tenantID int64) ([]domain.Task, error) {
					return nil, errors.New("database error")
				}
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockTaskRepository{}
			mockDispatcher := &mocks.MockDispatcher{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			service := NewService(mockRepo, mockDispatcher)
			ctx := createContextWithUserID(1)
			ctx = middleware.SetTenantIDInContext(ctx, 1)
			tasks, err := service.ListTasks(ctx)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("erro esperado '%v', mas nenhum erro foi retornado", tt.expectedError)
					return
				}
				if !errors.Is(err, tt.expectedError) && err.Error() != tt.expectedError.Error() {
					t.Errorf("erro esperado '%v', obtido '%v'", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("erro inesperado: %v", err)
				return
			}

			if len(tasks) != tt.expectedLength {
				t.Errorf("número de tarefas esperado %d, obtido %d", tt.expectedLength, len(tasks))
			}
		})
	}
}

func TestService_UpdateTask(t *testing.T) {
	now := time.Now()
	completedTime := now.Add(12 * time.Hour)

	tests := []struct {
		name          string
		ctx           context.Context
		id            int64
		req           UpdateTaskRequest
		mockSetup     func(*mocks.MockTaskRepository)
		expectedError error
		validateTask  func(*testing.T, *domain.Task)
	}{
		{
			name: "sucesso ao atualizar título",
			ctx:  createContextWithUserID(1),
			id:   1,
			req: UpdateTaskRequest{
				Title: stringPtr("Novo título"),
			},
			mockSetup: func(m *mocks.MockTaskRepository) {
				m.GetByIDFunc = func(ctx context.Context, id int64, tenantID int64) (*domain.Task, error) {
					return &domain.Task{
						ID:          1,
						Title:       "Título antigo",
						Description: "Descrição",
						Points:      10,
						Status:      "pending",
					}, nil
				}
				m.UpdateFunc = func(ctx context.Context, task *domain.Task) error {
					return nil
				}
			},
			validateTask: func(t *testing.T, task *domain.Task) {
				if task.Title != "Novo título" {
					t.Errorf("Title esperado 'Novo título', obtido '%s'", task.Title)
				}
			},
		},
		{
			name: "sucesso ao marcar como completa e calcular próxima data",
			ctx:  createContextWithUserID(1),
			id:   1,
			req: UpdateTaskRequest{
				Completed:   boolPtr(true),
				CompletedAt: &completedTime,
			},
			mockSetup: func(m *mocks.MockTaskRepository) {
				m.GetByIDFunc = func(ctx context.Context, id int64, tenantID int64) (*domain.Task, error) {
					return &domain.Task{
						ID:             1,
						Title:          "Tarefa",
						FrequencyValue: 7,
						FrequencyUnit:  domain.UnitDays,
						Completed:      false,
					}, nil
				}
				m.UpdateFunc = func(ctx context.Context, task *domain.Task) error {
					return nil
				}
			},
			validateTask: func(t *testing.T, task *domain.Task) {
				if !task.Completed {
					t.Error("Completed deve ser true")
				}
				if task.CompletedById == nil || *task.CompletedById != 1 {
					t.Error("CompletedById deve ser 1")
				}
				if task.ScheduledTo == nil {
					t.Error("ScheduledTo não deve ser nil após completar")
				}
				expectedNextDate := completedTime.AddDate(0, 0, 7)
				if task.ScheduledTo != nil && !task.ScheduledTo.Equal(expectedNextDate) {
					t.Errorf("ScheduledTo esperado %v, obtido %v", expectedNextDate, task.ScheduledTo)
				}
			},
		},
		{
			name: "sucesso ao desmarcar como completa",
			ctx:  createContextWithUserID(1),
			id:   1,
			req: UpdateTaskRequest{
				Completed: boolPtr(false),
			},
			mockSetup: func(m *mocks.MockTaskRepository) {
				userID := int64(1)
				m.GetByIDFunc = func(ctx context.Context, id int64, tenantID int64) (*domain.Task, error) {
					return &domain.Task{
						ID:           1,
						Title:        "Tarefa",
						Completed:    true,
						CompletedById: &userID,
					}, nil
				}
				m.UpdateFunc = func(ctx context.Context, task *domain.Task) error {
					return nil
				}
			},
			validateTask: func(t *testing.T, task *domain.Task) {
				if task.Completed {
					t.Error("Completed deve ser false")
				}
				if task.CompletedById != nil {
					t.Error("CompletedById deve ser nil")
				}
			},
		},
		{
			name: "erro quando usuário não autenticado",
			ctx:  context.Background(),
			id:   1,
			req: UpdateTaskRequest{
				Title: stringPtr("Novo título"),
			},
			expectedError: ErrUserNotAuthenticated,
		},
		{
			name: "erro quando tarefa não encontrada",
			ctx:  createContextWithUserID(1),
			id:   999,
			req: UpdateTaskRequest{
				Title: stringPtr("Novo título"),
			},
			mockSetup: func(m *mocks.MockTaskRepository) {
				m.GetByIDFunc = func(ctx context.Context, id int64, tenantID int64) (*domain.Task, error) {
					return nil, nil
				}
			},
			expectedError: ErrTaskNotFound,
		},
		{
			name: "erro ao atualizar no repositório",
			ctx:  createContextWithUserID(1),
			id:   1,
			req: UpdateTaskRequest{
				Title: stringPtr("Novo título"),
			},
			mockSetup: func(m *mocks.MockTaskRepository) {
				m.GetByIDFunc = func(ctx context.Context, id int64, tenantID int64) (*domain.Task, error) {
					return &domain.Task{ID: 1, Title: "Tarefa"}, nil
				}
				m.UpdateFunc = func(ctx context.Context, task *domain.Task) error {
					return errors.New("database error")
				}
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockTaskRepository{}
			mockDispatcher := &mocks.MockDispatcher{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			service := NewService(mockRepo, mockDispatcher)
			ctx := tt.ctx
			if userID := middleware.GetUserIDFromContext(ctx); userID > 0 {
				ctx = middleware.SetTenantIDInContext(ctx, 1)
			}
			task, err := service.UpdateTask(ctx, tt.id, tt.req)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("erro esperado '%v', mas nenhum erro foi retornado", tt.expectedError)
					return
				}
				if !errors.Is(err, tt.expectedError) && err.Error() != tt.expectedError.Error() {
					t.Errorf("erro esperado '%v', obtido '%v'", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("erro inesperado: %v", err)
				return
			}

			if task == nil {
				t.Fatal("task não deveria ser nil")
			}

			if tt.validateTask != nil {
				tt.validateTask(t, task)
			}
		})
	}
}

func TestService_DeleteTask(t *testing.T) {
	tests := []struct {
		name          string
		id            int64
		mockSetup     func(*mocks.MockTaskRepository)
		expectedError error
	}{
		{
			name: "sucesso ao deletar tarefa",
			id:   1,
			mockSetup: func(m *mocks.MockTaskRepository) {
				m.DeleteFunc = func(ctx context.Context, id int64, tenantID int64) error {
					return nil
				}
			},
		},
		{
			name: "erro do repositório",
			id:   1,
			mockSetup: func(m *mocks.MockTaskRepository) {
				m.DeleteFunc = func(ctx context.Context, id int64, tenantID int64) error {
					return errors.New("database error")
				}
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockTaskRepository{}
			mockDispatcher := &mocks.MockDispatcher{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			service := NewService(mockRepo, mockDispatcher)
			ctx := createContextWithUserID(1)
			ctx = middleware.SetTenantIDInContext(ctx, 1)
			err := service.DeleteTask(ctx, tt.id)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("erro esperado '%v', mas nenhum erro foi retornado", tt.expectedError)
					return
				}
				if !errors.Is(err, tt.expectedError) && err.Error() != tt.expectedError.Error() {
					t.Errorf("erro esperado '%v', obtido '%v'", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("erro inesperado: %v", err)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
