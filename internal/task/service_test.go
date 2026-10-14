package task

import (
	"context"
	"errors"
	"testing"
)

func TestServiceCreateNormalizesInput(t *testing.T) {
	repository := &memoryRepository{}
	service := NewService(repository)

	item, err := service.Create(context.Background(), CreateInput{Title: "  Write tests  "})
	if err != nil {
		t.Fatal(err)
	}
	if item.Title != "Write tests" || item.Status != "todo" || item.ID != 1 {
		t.Fatalf("unexpected task: %+v", item)
	}
}

func TestServiceRejectsInvalidInput(t *testing.T) {
	tests := []CreateInput{
		{Title: ""},
		{Title: "valid", Status: "later"},
		{Title: "valid", Description: string(make([]byte, 5001))},
	}
	service := NewService(&memoryRepository{})
	for _, input := range tests {
		if _, err := service.Create(context.Background(), input); err == nil {
			t.Fatalf("expected invalid input %+v to fail", input)
		} else {
			var validationError ValidationError
			if !errors.As(err, &validationError) {
				t.Fatalf("expected ValidationError, got %T", err)
			}
		}
	}
}

func TestServiceUpdateAndDelete(t *testing.T) {
	service := NewService(&memoryRepository{})
	created, err := service.Create(context.Background(), CreateInput{Title: "Draft"})
	if err != nil {
		t.Fatal(err)
	}
	updated, err := service.Update(context.Background(), created.ID, UpdateInput{Title: "Published", Status: "done"})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "Published" || updated.Status != "done" {
		t.Fatalf("unexpected update: %+v", updated)
	}
	if err := service.Delete(context.Background(), created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Get(context.Background(), created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected a missing task after delete, got %v", err)
	}
}
