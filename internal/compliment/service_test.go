package compliment

import (
	"context"
	"errors"
	"keep-your-house-clean/internal/domain"
	"keep-your-house-clean/internal/events"
	"keep-your-house-clean/internal/platform/middleware"
	"keep-your-house-clean/internal/compliment/mocks"
	"testing"
)

func createContextWithUserID(userID int64) context.Context {
	return middleware.SetUserIDInContext(context.Background(), userID)
}

func TestNewService(t *testing.T) {
	repo := &mocks.MockComplimentRepository{}
	userRepo := &mocks.MockUserRepository{}
	dispatcher := &mocks.MockDispatcher{}
	service := NewService(repo, userRepo, dispatcher)

	if service == nil {
		t.Fatal("NewService retornou nil")
	}

	if service.repo != repo {
		t.Error("Repository não foi atribuído corretamente")
	}

	if service.userRepo != userRepo {
		t.Error("UserRepository não foi atribuído corretamente")
	}

	if service.dispatcher != dispatcher {
		t.Error("Dispatcher não foi atribuído corretamente")
	}
}

func TestService_CreateCompliment(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		req           CreateComplimentRequest
		mockSetup     func(*mocks.MockComplimentRepository, *mocks.MockUserRepository, *mocks.MockDispatcher)
		expectedError error
		validateCompliment func(*testing.T, *domain.Compliment)
	}{
		{
			name: "sucesso ao criar elogio",
			ctx:  createContextWithUserID(1),
			req: CreateComplimentRequest{
				Title:       "Ótimo trabalho",
				Description: "Você fez um excelente trabalho",
				Points:      5,
				ToUserID:    2,
			},
			mockSetup: func(cr *mocks.MockComplimentRepository, ur *mocks.MockUserRepository, d *mocks.MockDispatcher) {
				ur.GetByIDFunc = func(ctx context.Context, id int64) (*domain.User, error) {
					return &domain.User{
						ID:       2,
						TenantID: 1,
					}, nil
				}
				cr.CreateFunc = func(ctx context.Context, compliment *domain.Compliment) error {
					compliment.ID = 1
					return nil
				}
				d.DispatchFunc = func(event events.Event) error {
					return nil
				}
			},
			validateCompliment: func(t *testing.T, compliment *domain.Compliment) {
				if compliment.ID != 1 {
					t.Errorf("ID esperado 1, obtido %d", compliment.ID)
				}
				if compliment.Title != "Ótimo trabalho" {
					t.Errorf("Title esperado 'Ótimo trabalho', obtido '%s'", compliment.Title)
				}
				if compliment.Points != 5 {
					t.Errorf("Points esperado 5, obtido %d", compliment.Points)
				}
				if compliment.FromUserID != 1 {
					t.Errorf("FromUserID esperado 1, obtido %d", compliment.FromUserID)
				}
				if compliment.ToUserID != 2 {
					t.Errorf("ToUserID esperado 2, obtido %d", compliment.ToUserID)
				}
			},
		},
		{
			name: "erro quando usuário não autenticado",
			ctx:  context.Background(),
			req: CreateComplimentRequest{
				Title:       "Ótimo trabalho",
				Description: "Você fez um excelente trabalho",
				Points:      5,
				ToUserID:    2,
			},
			expectedError: ErrUserNotAuthenticated,
		},
		{
			name: "erro quando pontos inválidos (maior que 5)",
			ctx:  createContextWithUserID(1),
			req: CreateComplimentRequest{
				Title:       "Ótimo trabalho",
				Description: "Você fez um excelente trabalho",
				Points:      6,
				ToUserID:    2,
			},
			expectedError: ErrInvalidPoints,
		},
		{
			name: "erro quando pontos inválidos (menor que 0)",
			ctx:  createContextWithUserID(1),
			req: CreateComplimentRequest{
				Title:       "Ótimo trabalho",
				Description: "Você fez um excelente trabalho",
				Points:      -1,
				ToUserID:    2,
			},
			expectedError: ErrInvalidPoints,
		},
		{
			name: "erro quando usuário tenta elogiar a si mesmo",
			ctx:  createContextWithUserID(1),
			req: CreateComplimentRequest{
				Title:       "Ótimo trabalho",
				Description: "Você fez um excelente trabalho",
				Points:      5,
				ToUserID:    1,
			},
			expectedError: ErrInvalidUser,
		},
		{
			name: "erro quando usuário não encontrado",
			ctx:  createContextWithUserID(1),
			req: CreateComplimentRequest{
				Title:       "Ótimo trabalho",
				Description: "Você fez um excelente trabalho",
				Points:      5,
				ToUserID:    999,
			},
			mockSetup: func(cr *mocks.MockComplimentRepository, ur *mocks.MockUserRepository, d *mocks.MockDispatcher) {
				ur.GetByIDFunc = func(ctx context.Context, id int64) (*domain.User, error) {
					return nil, nil
				}
			},
			expectedError: ErrUserNotFound,
		},
		{
			name: "erro quando usuário de outro tenant",
			ctx:  createContextWithUserID(1),
			req: CreateComplimentRequest{
				Title:       "Ótimo trabalho",
				Description: "Você fez um excelente trabalho",
				Points:      5,
				ToUserID:    2,
			},
			mockSetup: func(cr *mocks.MockComplimentRepository, ur *mocks.MockUserRepository, d *mocks.MockDispatcher) {
				ur.GetByIDFunc = func(ctx context.Context, id int64) (*domain.User, error) {
					return &domain.User{
						ID:       2,
						TenantID: 999,
					}, nil
				}
			},
			expectedError: ErrUserNotFound,
		},
		{
			name: "erro ao salvar no repositório",
			ctx:  createContextWithUserID(1),
			req: CreateComplimentRequest{
				Title:       "Ótimo trabalho",
				Description: "Você fez um excelente trabalho",
				Points:      5,
				ToUserID:    2,
			},
			mockSetup: func(cr *mocks.MockComplimentRepository, ur *mocks.MockUserRepository, d *mocks.MockDispatcher) {
				ur.GetByIDFunc = func(ctx context.Context, id int64) (*domain.User, error) {
					return &domain.User{
						ID:       2,
						TenantID: 1,
					}, nil
				}
				cr.CreateFunc = func(ctx context.Context, compliment *domain.Compliment) error {
					return errors.New("database error")
				}
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockComplimentRepository{}
			mockUserRepo := &mocks.MockUserRepository{}
			mockDispatcher := &mocks.MockDispatcher{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo, mockUserRepo, mockDispatcher)
			}

			service := NewService(mockRepo, mockUserRepo, mockDispatcher)
			ctx := tt.ctx
			if userID := middleware.GetUserIDFromContext(ctx); userID > 0 {
				ctx = middleware.SetTenantIDInContext(ctx, 1)
			}
			compliment, err := service.CreateCompliment(ctx, tt.req)

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

			if compliment == nil {
				t.Fatal("compliment não deveria ser nil")
			}

			if tt.validateCompliment != nil {
				tt.validateCompliment(t, compliment)
			}
		})
	}
}

func TestService_GetComplimentByID(t *testing.T) {
	tests := []struct {
		name          string
		id            int64
		mockSetup     func(*mocks.MockComplimentRepository)
		expectedError error
		validateCompliment func(*testing.T, *domain.Compliment)
	}{
		{
			name: "sucesso ao buscar elogio",
			id:   1,
			mockSetup: func(m *mocks.MockComplimentRepository) {
				m.GetByIDFunc = func(ctx context.Context, id int64, tenantID int64) (*domain.Compliment, error) {
					return &domain.Compliment{
						ID:          1,
						Title:       "Ótimo trabalho",
						Description: "Excelente",
						Points:      5,
					}, nil
				}
			},
			validateCompliment: func(t *testing.T, compliment *domain.Compliment) {
				if compliment.ID != 1 {
					t.Errorf("ID esperado 1, obtido %d", compliment.ID)
				}
				if compliment.Title != "Ótimo trabalho" {
					t.Errorf("Title esperado 'Ótimo trabalho', obtido '%s'", compliment.Title)
				}
			},
		},
		{
			name: "erro quando elogio não encontrado",
			id:   999,
			mockSetup: func(m *mocks.MockComplimentRepository) {
				m.GetByIDFunc = func(ctx context.Context, id int64, tenantID int64) (*domain.Compliment, error) {
					return nil, nil
				}
			},
			expectedError: ErrComplimentNotFound,
		},
		{
			name: "erro do repositório",
			id:   1,
			mockSetup: func(m *mocks.MockComplimentRepository) {
				m.GetByIDFunc = func(ctx context.Context, id int64, tenantID int64) (*domain.Compliment, error) {
					return nil, errors.New("database error")
				}
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockComplimentRepository{}
			mockUserRepo := &mocks.MockUserRepository{}
			mockDispatcher := &mocks.MockDispatcher{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			service := NewService(mockRepo, mockUserRepo, mockDispatcher)
			ctx := createContextWithUserID(1)
			ctx = middleware.SetTenantIDInContext(ctx, 1)
			compliment, err := service.GetComplimentByID(ctx, tt.id)

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

			if compliment == nil {
				t.Fatal("compliment não deveria ser nil")
			}

			if tt.validateCompliment != nil {
				tt.validateCompliment(t, compliment)
			}
		})
	}
}

func TestService_ListCompliments(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockComplimentRepository)
		expectedError  error
		expectedLength int
	}{
		{
			name: "sucesso ao listar elogios",
			mockSetup: func(m *mocks.MockComplimentRepository) {
				m.FetchAllFunc = func(ctx context.Context, tenantID int64) ([]domain.Compliment, error) {
					return []domain.Compliment{
						{ID: 1, Title: "Elogio 1"},
						{ID: 2, Title: "Elogio 2"},
					}, nil
				}
			},
			expectedLength: 2,
		},
		{
			name: "retorna lista vazia",
			mockSetup: func(m *mocks.MockComplimentRepository) {
				m.FetchAllFunc = func(ctx context.Context, tenantID int64) ([]domain.Compliment, error) {
					return []domain.Compliment{}, nil
				}
			},
			expectedLength: 0,
		},
		{
			name: "erro do repositório",
			mockSetup: func(m *mocks.MockComplimentRepository) {
				m.FetchAllFunc = func(ctx context.Context, tenantID int64) ([]domain.Compliment, error) {
					return nil, errors.New("database error")
				}
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockComplimentRepository{}
			mockUserRepo := &mocks.MockUserRepository{}
			mockDispatcher := &mocks.MockDispatcher{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			service := NewService(mockRepo, mockUserRepo, mockDispatcher)
			ctx := createContextWithUserID(1)
			ctx = middleware.SetTenantIDInContext(ctx, 1)
			compliments, err := service.ListCompliments(ctx)

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

			if len(compliments) != tt.expectedLength {
				t.Errorf("número de elogios esperado %d, obtido %d", tt.expectedLength, len(compliments))
			}
		})
	}
}

func TestService_DeleteCompliment(t *testing.T) {
	tests := []struct {
		name          string
		id            int64
		mockSetup     func(*mocks.MockComplimentRepository)
		expectedError error
	}{
		{
			name: "sucesso ao deletar elogio",
			id:   1,
			mockSetup: func(m *mocks.MockComplimentRepository) {
				m.DeleteFunc = func(ctx context.Context, id int64, tenantID int64) error {
					return nil
				}
			},
		},
		{
			name: "erro do repositório",
			id:   1,
			mockSetup: func(m *mocks.MockComplimentRepository) {
				m.DeleteFunc = func(ctx context.Context, id int64, tenantID int64) error {
					return errors.New("database error")
				}
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockComplimentRepository{}
			mockUserRepo := &mocks.MockUserRepository{}
			mockDispatcher := &mocks.MockDispatcher{}
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			service := NewService(mockRepo, mockUserRepo, mockDispatcher)
			ctx := createContextWithUserID(1)
			ctx = middleware.SetTenantIDInContext(ctx, 1)
			err := service.DeleteCompliment(ctx, tt.id)

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

