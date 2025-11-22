package compliment

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
	r.Route("/api/v1/compliments", func(r chi.Router) {
		r.Get("/", h.ListCompliments)
		r.Get("/history", h.GetUserComplimentsHistory)
		r.Get("/unviewed", h.GetUnviewedReceivedCompliments)
		r.Get("/last-received", h.GetLastReceivedCompliment)
		r.Post("/mark-viewed", h.MarkComplimentsAsViewed)
		r.Get("/{id}", h.GetCompliment)
		r.Post("/", h.CreateCompliment)
		r.Delete("/{id}", h.DeleteCompliment)
	})
}

func (h *Handler) CreateCompliment(w http.ResponseWriter, r *http.Request) {
	var req CreateComplimentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	compliment, err := h.service.CreateCompliment(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, ErrInvalidPoints) || errors.Is(err, ErrInvalidUser) || errors.Is(err, ErrUserNotFound) {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, compliment)
}

func (h *Handler) GetCompliment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid compliment ID")
		return
	}

	compliment, err := h.service.GetComplimentByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrComplimentNotFound) {
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

	respondWithJSON(w, http.StatusOK, compliment)
}

func (h *Handler) ListCompliments(w http.ResponseWriter, r *http.Request) {
	compliments, err := h.service.ListCompliments(r.Context())
	if err != nil {
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, compliments)
}

func (h *Handler) GetUserComplimentsHistory(w http.ResponseWriter, r *http.Request) {
	compliments, err := h.service.GetUserComplimentsHistory(r.Context())
	if err != nil {
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, compliments)
}

func (h *Handler) GetUnviewedReceivedCompliments(w http.ResponseWriter, r *http.Request) {
	compliments, err := h.service.GetUnviewedReceivedCompliments(r.Context())
	if err != nil {
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, compliments)
}

func (h *Handler) MarkComplimentsAsViewed(w http.ResponseWriter, r *http.Request) {
	var req MarkAsViewedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.MarkComplimentsAsViewed(r.Context(), req.IDs); err != nil {
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetLastReceivedCompliment(w http.ResponseWriter, r *http.Request) {
	compliment, err := h.service.GetLastReceivedCompliment(r.Context())
	if err != nil {
		if errors.Is(err, ErrUserNotAuthenticated) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if compliment == nil {
		respondWithJSON(w, http.StatusOK, nil)
		return
	}

	respondWithJSON(w, http.StatusOK, compliment)
}

func (h *Handler) DeleteCompliment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid compliment ID")
		return
	}

	if err := h.service.DeleteCompliment(r.Context(), id); err != nil {
		if errors.Is(err, ErrComplimentNotFound) {
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

