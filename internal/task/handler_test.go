package task

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestCreateTask(t *testing.T) {
	repo := &memoryRepository{}
	handler := NewHandler(NewService(repo))
	router := chi.NewRouter()
	router.Post("/tasks", handler.Create)

	request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"Prepare API","status":"todo"}`))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, response.Code, response.Body.String())
	}
	var got Task
	if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.ID != 1 || got.Title != "Prepare API" || got.Status != "todo" {
		t.Fatalf("unexpected task: %+v", got)
	}
}

func TestCreateTaskRejectsInvalidStatus(t *testing.T) {
	handler := NewHandler(NewService(&memoryRepository{}))
	request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"Prepare API","status":"later"}`))
	response := httptest.NewRecorder()
	handler.Create(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d: %s", http.StatusBadRequest, response.Code, response.Body.String())
	}
}

func TestTaskHandlerCRUD(t *testing.T) {
	handler := NewHandler(NewService(&memoryRepository{}))
	router := chi.NewRouter()
	router.Post("/tasks", handler.Create)
	router.Get("/tasks", handler.List)
	router.Get("/tasks/{id}", handler.Get)
	router.Put("/tasks/{id}", handler.Update)
	router.Delete("/tasks/{id}", handler.Delete)

	create := httptest.NewRecorder()
	router.ServeHTTP(create, httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"Draft"}`)))
	if create.Code != http.StatusCreated {
		t.Fatalf("create failed: %d %s", create.Code, create.Body.String())
	}

	list := httptest.NewRecorder()
	router.ServeHTTP(list, httptest.NewRequest(http.MethodGet, "/tasks", nil))
	if list.Code != http.StatusOK || !strings.Contains(list.Body.String(), "Draft") {
		t.Fatalf("list failed: %d %s", list.Code, list.Body.String())
	}

	update := httptest.NewRecorder()
	router.ServeHTTP(update, httptest.NewRequest(http.MethodPut, "/tasks/1", strings.NewReader(`{"title":"Final","status":"done"}`)))
	if update.Code != http.StatusOK || !strings.Contains(update.Body.String(), "Final") {
		t.Fatalf("update failed: %d %s", update.Code, update.Body.String())
	}

	get := httptest.NewRecorder()
	router.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/tasks/1", nil))
	if get.Code != http.StatusOK || !strings.Contains(get.Body.String(), "done") {
		t.Fatalf("get failed: %d %s", get.Code, get.Body.String())
	}

	deleteResponse := httptest.NewRecorder()
	router.ServeHTTP(deleteResponse, httptest.NewRequest(http.MethodDelete, "/tasks/1", nil))
	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("delete failed: %d %s", deleteResponse.Code, deleteResponse.Body.String())
	}

	notFound := httptest.NewRecorder()
	router.ServeHTTP(notFound, httptest.NewRequest(http.MethodGet, "/tasks/1", nil))
	if notFound.Code != http.StatusNotFound {
		t.Fatalf("expected missing task status 404, got %d", notFound.Code)
	}
}

func TestTaskHandlerRejectsBadRequest(t *testing.T) {
	handler := NewHandler(NewService(&memoryRepository{}))
	router := chi.NewRouter()
	router.Post("/tasks", handler.Create)
	router.Get("/tasks/{id}", handler.Get)

	unknownField := httptest.NewRecorder()
	router.ServeHTTP(unknownField, httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"Draft","unknown":true}`)))
	if unknownField.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request for unknown field, got %d", unknownField.Code)
	}

	invalidID := httptest.NewRecorder()
	router.ServeHTTP(invalidID, httptest.NewRequest(http.MethodGet, "/tasks/not-a-number", nil))
	if invalidID.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request for invalid id, got %d", invalidID.Code)
	}
}

type memoryRepository struct {
	items []Task
}

func (m *memoryRepository) Create(_ context.Context, item *Task) error {
	item.ID = uint(len(m.items) + 1)
	item.CreatedAt = time.Now()
	m.items = append(m.items, *item)
	return nil
}

func (m *memoryRepository) List(_ context.Context) ([]Task, error) { return m.items, nil }

func (m *memoryRepository) Get(_ context.Context, id uint) (Task, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return Task{}, ErrNotFound
}

func (m *memoryRepository) Update(_ context.Context, item *Task) error {
	for index, stored := range m.items {
		if stored.ID == item.ID {
			m.items[index] = *item
			return nil
		}
	}
	return ErrNotFound
}

func (m *memoryRepository) Delete(_ context.Context, id uint) error {
	for index, item := range m.items {
		if item.ID == id {
			m.items = append(m.items[:index], m.items[index+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
