package task

import (
	"context"
	"fmt"
	"strings"
)

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

// ValidationError is safe to return to an API client. Infrastructure errors
// remain opaque at the HTTP boundary.
type ValidationError struct{ Message string }

func (e ValidationError) Error() string { return e.Message }

func (s *Service) Create(ctx context.Context, input CreateInput) (Task, error) {
	item, err := taskFromInput(input)
	if err != nil {
		return Task{}, err
	}
	if err := s.repository.Create(ctx, &item); err != nil {
		return Task{}, fmt.Errorf("create task: %w", err)
	}
	return item, nil
}

func (s *Service) List(ctx context.Context) ([]Task, error)       { return s.repository.List(ctx) }
func (s *Service) Get(ctx context.Context, id uint) (Task, error) { return s.repository.Get(ctx, id) }

func (s *Service) Update(ctx context.Context, id uint, input UpdateInput) (Task, error) {
	item, err := s.repository.Get(ctx, id)
	if err != nil {
		return Task{}, err
	}
	updated, err := taskFromInput(CreateInput(input))
	if err != nil {
		return Task{}, err
	}
	item.Title, item.Description, item.Status = updated.Title, updated.Description, updated.Status
	if err := s.repository.Update(ctx, &item); err != nil {
		return Task{}, err
	}
	return item, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error { return s.repository.Delete(ctx, id) }

func taskFromInput(input CreateInput) (Task, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" || len(title) > 200 {
		return Task{}, ValidationError{Message: "title must be between 1 and 200 characters"}
	}
	if len(input.Description) > 5000 {
		return Task{}, ValidationError{Message: "description must be at most 5000 characters"}
	}
	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "todo"
	}
	if status != "todo" && status != "in_progress" && status != "done" {
		return Task{}, ValidationError{Message: "status must be todo, in_progress, or done"}
	}
	return Task{Title: title, Description: input.Description, Status: status}, nil
}
