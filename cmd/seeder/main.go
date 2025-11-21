package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"keep-your-house-clean/internal/auth"
	"keep-your-house-clean/internal/domain"
	"keep-your-house-clean/internal/platform/database"
)

func main() {
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "keep_your_house_clean"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := waitForTables(db); err != nil {
		log.Fatalf("Failed to verify database tables: %v", err)
	}

	ctx := context.Background()

	tenantRepo := database.NewTenantRepository(db)
	userRepo := database.NewUserRepository(db)
	taskRepo := database.NewTaskRepository(db)

	tenant, err := seedTenant(ctx, tenantRepo)
	if err != nil {
		log.Fatalf("Failed to seed tenant: %v", err)
	}
	fmt.Printf("✓ Tenant criado: %s (ID: %d)\n", tenant.Name, tenant.ID)

	user, err := seedUser(ctx, userRepo, tenant.ID)
	if err != nil {
		log.Fatalf("Failed to seed user: %v", err)
	}
	fmt.Printf("✓ Usuário criado: %s (ID: %d, Email: %s)\n", user.Name, user.ID, user.Email)

	tasks, err := seedTasks(ctx, taskRepo, tenant.ID, user.ID)
	if err != nil {
		log.Fatalf("Failed to seed tasks: %v", err)
	}
	fmt.Printf("✓ %d tarefas criadas\n", len(tasks))

	fmt.Println("\nSeeder executado com sucesso!")
	fmt.Printf("\nCredenciais de acesso:\n")
	fmt.Printf("  Email: %s\n", user.Email)
	fmt.Printf("  Senha: password123\n")
}

func seedTenant(ctx context.Context, repo domain.TenantRepository) (*domain.Tenant, error) {
	existingTenant, err := repo.GetByDomain(ctx, "example.com")
	if err != nil {
		return nil, err
	}

	if existingTenant != nil {
		fmt.Println("⚠ Tenant já existe, usando existente")
		return existingTenant, nil
	}

	now := time.Now()
	tenant := &domain.Tenant{
		Name:      "Tenant de Exemplo",
		Domain:    "example.com",
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := repo.Create(ctx, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

func seedUser(ctx context.Context, repo domain.UserRepository, tenantID int64) (*domain.User, error) {
	existingUser, err := repo.GetByEmailAndTenant(ctx, "admin@example.com", tenantID)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		fmt.Println("⚠ Usuário já existe, usando existente")
		return existingUser, nil
	}

	hashedPassword, err := auth.HashPassword("password123")
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		Name:       "Admin Usuário",
		Email:      "admin@example.com",
		Password:   hashedPassword,
		TenantID:   tenantID,
		Points:     0,
		Role:       "admin",
		Status:     "active",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func seedTasks(ctx context.Context, repo domain.TaskRepository, tenantID, userID int64) ([]domain.Task, error) {
	now := time.Now()
	tasks := []struct {
		Title          string
		Description    string
		Points         int
		Status         string
		ScheduledTo    *time.Time
		FrequencyValue int
		FrequencyUnit  domain.FrequencyUnit
	}{
		{
			Title:          "Limpar a cozinha",
			Description:    "Lavar a louça, limpar a bancada e organizar os armários",
			Points:         10,
			Status:         "pending",
			ScheduledTo:    timePtr(now.Add(24 * time.Hour)),
			FrequencyValue: 1,
			FrequencyUnit:  domain.UnitDays,
		},
		{
			Title:          "Aspirar a sala",
			Description:    "Aspirar o tapete e os móveis",
			Points:         5,
			Status:         "pending",
			ScheduledTo:    timePtr(now.Add(48 * time.Hour)),
			FrequencyValue: 2,
			FrequencyUnit:  domain.UnitDays,
		},
		{
			Title:          "Lavar as roupas",
			Description:    "Separar, lavar e dobrar as roupas",
			Points:         15,
			Status:         "pending",
			ScheduledTo:    timePtr(now.Add(72 * time.Hour)),
			FrequencyValue: 3,
			FrequencyUnit:  domain.UnitDays,
		},
		{
			Title:          "Limpar o banheiro",
			Description:    "Limpar pia, vaso sanitário e box",
			Points:         12,
			Status:         "pending",
			ScheduledTo:    timePtr(now.Add(4 * 24 * time.Hour)),
			FrequencyValue: 1,
			FrequencyUnit:  domain.UnitWeeks,
		},
		{
			Title:          "Organizar o quarto",
			Description:    "Arrumar a cama, organizar o guarda-roupa e limpar o chão",
			Points:         8,
			Status:         "pending",
			ScheduledTo:    timePtr(now.Add(5 * 24 * time.Hour)),
			FrequencyValue: 1,
			FrequencyUnit:  domain.UnitWeeks,
		},
		{
			Title:          "Limpar as janelas",
			Description:    "Limpar todas as janelas da casa",
			Points:         20,
			Status:         "pending",
			ScheduledTo:    timePtr(now.Add(7 * 24 * time.Hour)),
			FrequencyValue: 1,
			FrequencyUnit:  domain.UnitMonths,
		},
		{
			Title:          "Varrer o quintal",
			Description:    "Varrer folhas e limpar a área externa",
			Points:         7,
			Status:         "pending",
			ScheduledTo:    timePtr(now.Add(8 * 24 * time.Hour)),
			FrequencyValue: 1,
			FrequencyUnit:  domain.UnitWeeks,
		},
		{
			Title:          "Limpar a geladeira",
			Description:    "Organizar e limpar o interior da geladeira",
			Points:         10,
			Status:         "pending",
			ScheduledTo:    timePtr(now.Add(10 * 24 * time.Hour)),
			FrequencyValue: 1,
			FrequencyUnit:  domain.UnitMonths,
		},
		{
			Title:          "Trocar lençóis",
			Description:    "Trocar lençóis de todas as camas",
			Points:         6,
			Status:         "pending",
			ScheduledTo:    timePtr(now.Add(12 * 24 * time.Hour)),
			FrequencyValue: 2,
			FrequencyUnit:  domain.UnitWeeks,
		},
		{
			Title:          "Limpar o forno",
			Description:    "Limpar o interior e exterior do forno",
			Points:         15,
			Status:         "pending",
			ScheduledTo:    timePtr(now.Add(14 * 24 * time.Hour)),
			FrequencyValue: 1,
			FrequencyUnit:  domain.UnitMonths,
		},
	}

	var createdTasks []domain.Task
	for _, t := range tasks {
		task := &domain.Task{
			Title:          t.Title,
			Description:    t.Description,
			Points:         t.Points,
			Status:         t.Status,
			ScheduledTo:    t.ScheduledTo,
			ScheduledById:  &userID,
			FrequencyValue: t.FrequencyValue,
			FrequencyUnit:  t.FrequencyUnit,
			Completed:      false,
			CompletedById:  nil,
			TenantID:       tenantID,
			CreatedAt:      now,
			CreatedById:    userID,
			UpdatedAt:      now,
			UpdatedById:    nil,
			DeletedAt:      nil,
		}

		if err := repo.Create(ctx, task); err != nil {
			return nil, err
		}

		createdTasks = append(createdTasks, *task)
	}

	return createdTasks, nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func waitForTables(db *sql.DB) error {
	tables := []string{"tenants", "users", "tasks"}
	maxAttempts := 30
	attempt := 0

	for attempt < maxAttempts {
		allExist := true
		for _, table := range tables {
			var exists bool
			query := `
				SELECT EXISTS (
					SELECT FROM information_schema.tables 
					WHERE table_schema = 'public' 
					AND table_name = $1
				)
			`
			err := db.QueryRow(query, table).Scan(&exists)
			if err != nil || !exists {
				allExist = false
				break
			}
		}

		if allExist {
			return nil
		}

		attempt++
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("tables not found after %d attempts", maxAttempts)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

