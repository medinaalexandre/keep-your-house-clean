package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/tasks", func(r chi.Router) {
		r.Get("/", h.ListTasks)
		r.Get("/upcoming", h.GetUpcomingTasks)
		r.Get("/history", h.GetCompletedTasksHistory)
		r.Get("/user/{userId}/completed", h.GetCompletedTasksByUser)
		r.Get("/{id}", h.GetTask)
		r.Post("/", h.CreateTask)
		r.Put("/{id}", h.UpdateTask)
		r.Post("/{id}/complete", h.CompleteTask)
		r.Post("/{id}/undo", h.UndoCompleteTask)
		r.Delete("/{id}", h.DeleteTask)
	})
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task, err := h.service.CreateTask(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.service.GetTaskByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.ListTasks(r.Context())
	if err != nil {
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	task, err := h.service.UpdateTask(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	req := CompleteTaskRequest{}
	if r.Header.Get("Content-Type") == "application/json" && r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
	}

	task, err := h.service.CompleteTask(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, ErrTaskAlreadyCompleted) {
			respondWithError(w, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetUpcomingTasks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 5
	offset := 0

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	tasks, err := h.service.GetUpcomingTasks(r.Context(), limit, offset)
	if err != nil {
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handler) GetCompletedTasksHistory(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	
	limit := 5
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	tasks, err := h.service.GetCompletedTasksHistory(r.Context(), limit)
	if err != nil {
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handler) GetCompletedTasksByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	tasks, err := h.service.GetCompletedTasksByUser(r.Context(), userID, limit, offset)
	if err != nil {
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handler) UndoCompleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.service.UndoCompleteTask(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, ErrTaskNotCompleted) {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

