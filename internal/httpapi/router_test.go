package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-user/go-service-template/internal/task"
)

func TestHealth(t *testing.T) {
	tests := []struct {
		name       string
		pingerErr  error
		wantStatus int
	}{
		{name: "database available", wantStatus: http.StatusOK},
		{name: "database unavailable", pingerErr: errors.New("connection refused"), wantStatus: http.StatusServiceUnavailable},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := NewRouter(task.NewService(testRepository{}), testPinger{err: test.pingerErr})
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health", nil))
			if recorder.Code != test.wantStatus {
				t.Fatalf("expected status %d, got %d", test.wantStatus, recorder.Code)
			}
		})
	}
}

type testRepository struct{}

func (testRepository) Create(context.Context, *task.Task) error  { return nil }
func (testRepository) List(context.Context) ([]task.Task, error) { return nil, nil }
func (testRepository) Get(context.Context, uint) (task.Task, error) {
	return task.Task{}, task.ErrNotFound
}
func (testRepository) Update(context.Context, *task.Task) error { return nil }
func (testRepository) Delete(context.Context, uint) error       { return nil }
