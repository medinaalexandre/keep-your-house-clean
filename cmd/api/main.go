package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"keep-your-house-clean/internal/auth"
	complimentHandler "keep-your-house-clean/internal/compliment"
	"keep-your-house-clean/internal/events"
	eventHandlers "keep-your-house-clean/internal/events/handlers"
	"keep-your-house-clean/internal/platform/database"
	authMiddleware "keep-your-house-clean/internal/platform/middleware"
	taskHandler "keep-your-house-clean/internal/task"
	tenantHandler "keep-your-house-clean/internal/tenant"
	userHandler "keep-your-house-clean/internal/user"
)

func main() {
	db, err := database.NewPostgresDB(database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "keep_your_house_clean"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	tenantRepo := database.NewTenantRepository(db)
	tenantService := tenantHandler.NewService(tenantRepo)
	tenantHandlerInstance := tenantHandler.NewHandler(tenantService)

	userRepo := database.NewUserRepository(db)
	userService := userHandler.NewService(userRepo)
	userHandlerInstance := userHandler.NewHandler(userService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := events.NewDispatcher(ctx)
	userPointsHandler := eventHandlers.NewUserPointsHandler(userRepo)
	dispatcher.RegisterHandler(events.EventTypeTaskCompleted, userPointsHandler.Handle)
	dispatcher.RegisterHandler(events.EventTypeTaskUndone, userPointsHandler.Handle)
	dispatcher.RegisterHandler(events.EventTypeComplimentReceived, userPointsHandler.Handle)
	dispatcher.Start()

	taskRepo := database.NewTaskRepository(db)
	taskService := taskHandler.NewService(taskRepo, dispatcher)
	taskHandlerInstance := taskHandler.NewHandler(taskService)

	complimentRepo := database.NewComplimentRepository(db)
	complimentService := complimentHandler.NewService(complimentRepo, userRepo, dispatcher)
	complimentHandlerInstance := complimentHandler.NewHandler(complimentService)

	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	authService := auth.NewService(userRepo, tenantRepo, jwtSecret)
	authHandlerInstance := auth.NewHandler(authService)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)

	authHandlerInstance.RegisterRoutes(r)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.JWTAuthMiddleware(jwtSecret))
		tenantHandlerInstance.RegisterRoutes(r)
		userHandlerInstance.RegisterRoutes(r)
		taskHandlerInstance.RegisterRoutes(r)
		complimentHandlerInstance.RegisterRoutes(r)
	})

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		if strings.HasPrefix(path, "/api") {
			http.NotFound(w, req)
			return
		}

		fileServer := http.FileServer(http.Dir("./web/dist"))
		if path == "/" || !strings.Contains(filepath.Base(path), ".") {
			http.ServeFile(w, req, "./web/dist/index.html")
			return
		}

		fullPath := filepath.Join("./web/dist", path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			http.ServeFile(w, req, "./web/dist/index.html")
			return
		}

		req.URL.Path = path
		fileServer.ServeHTTP(w, req)
	})

	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(":"+port, r); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Println("Server started successfully")
	<-sigChan
	log.Println("Shutting down server...")
	cancel()
	dispatcher.Stop()
	log.Println("Server stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
