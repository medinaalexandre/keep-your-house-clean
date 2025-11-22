package handlers

import (
	"context"
	"keep-your-house-clean/internal/domain"
	"keep-your-house-clean/internal/events"
	"time"
)

type UserPointsHandler struct {
	userRepo domain.UserRepository
}

func NewUserPointsHandler(userRepo domain.UserRepository) *UserPointsHandler {
	return &UserPointsHandler{
		userRepo: userRepo,
	}
}

func (h *UserPointsHandler) Handle(ctx context.Context, event events.Event) error {
	if event.Type == events.EventTypeTaskCompleted {
		return h.handleTaskCompleted(ctx, event)
	}
	
	if event.Type == events.EventTypeTaskUndone {
		return h.handleTaskUndone(ctx, event)
	}
	
	if event.Type == events.EventTypeComplimentReceived {
		return h.handleComplimentReceived(ctx, event)
	}
	
	return nil
}

func (h *UserPointsHandler) handleTaskCompleted(ctx context.Context, event events.Event) error {
	payload, ok := event.Payload.(events.TaskCompletedPayload)
	if !ok {
		return nil
	}

	if payload.Points <= 0 {
		return nil
	}

	user, err := h.userRepo.GetByID(ctx, payload.CompletedBy)
	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	user.Points += payload.Points
	user.UpdatedAt = time.Now()

	return h.userRepo.Update(ctx, user)
}

func (h *UserPointsHandler) handleTaskUndone(ctx context.Context, event events.Event) error {
	payload, ok := event.Payload.(events.TaskUndonePayload)
	if !ok {
		return nil
	}

	if payload.Points <= 0 {
		return nil
	}

	user, err := h.userRepo.GetByID(ctx, payload.CompletedBy)
	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	user.Points -= payload.Points
	if user.Points < 0 {
		user.Points = 0
	}
	user.UpdatedAt = time.Now()

	return h.userRepo.Update(ctx, user)
}

func (h *UserPointsHandler) handleComplimentReceived(ctx context.Context, event events.Event) error {
	payload, ok := event.Payload.(events.ComplimentReceivedPayload)
	if !ok {
		return nil
	}

	if payload.Points <= 0 {
		return nil
	}

	user, err := h.userRepo.GetByID(ctx, payload.ToUser)
	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	user.Points += payload.Points
	user.UpdatedAt = time.Now()

	return h.userRepo.Update(ctx, user)
}

