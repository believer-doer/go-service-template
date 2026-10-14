package task

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestPostgresRepositoryCRUD(t *testing.T) {
	repository, mock, closeDatabase := newMockRepository(t)
	defer closeDatabase()
	ctx := context.Background()

	item := Task{Title: "Write tests", Description: "Repository tests", Status: "todo"}
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tasks"`)).
		WithArgs(item.Title, item.Description, item.Status, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	if err := repository.Create(ctx, &item); err != nil {
		t.Fatalf("create: %v", err)
	}
	if item.ID != 1 {
		t.Fatalf("expected created id 1, got %d", item.ID)
	}

	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" ORDER BY id ASC`)).
		WillReturnRows(taskRows().AddRow(1, "Write tests", "Repository tests", "todo", createdAt, createdAt))
	items, err := repository.List(ctx)
	if err != nil || len(items) != 1 || items[0].Title != "Write tests" {
		t.Fatalf("list: items=%+v err=%v", items, err)
	}

	mock.ExpectQuery(`SELECT \* FROM "tasks"`).
		WithArgs(1, 1).
		WillReturnRows(taskRows().AddRow(1, "Write tests", "Repository tests", "todo", createdAt, createdAt))
	got, err := repository.Get(ctx, 1)
	if err != nil || got.ID != 1 {
		t.Fatalf("get: task=%+v err=%v", got, err)
	}

	got.Status = "done"
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tasks"`)).
		WithArgs(got.Title, got.Description, got.Status, got.CreatedAt, sqlmock.AnyArg(), got.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repository.Update(ctx, &got); err != nil {
		t.Fatalf("update: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tasks" WHERE "tasks"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	if err := repository.Delete(ctx, 1); err != nil {
		t.Fatalf("delete: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestPostgresRepositoryReturnsNotFound(t *testing.T) {
	repository, mock, closeDatabase := newMockRepository(t)
	defer closeDatabase()

	mock.ExpectQuery(`SELECT \* FROM "tasks"`).
		WithArgs(99, 1).
		WillReturnRows(taskRows())
	if _, err := repository.Get(context.Background(), 99); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tasks" WHERE "tasks"."id" = $1`)).
		WithArgs(99).
		WillReturnResult(sqlmock.NewResult(0, 0))
	if err := repository.Delete(context.Background(), 99); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func newMockRepository(t *testing.T) (*PostgresRepository, sqlmock.Sqlmock, func()) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		_ = sqlDB.Close()
		t.Fatal(err)
	}
	return NewPostgresRepository(db), mock, func() { _ = sqlDB.Close() }
}

func taskRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "title", "description", "status", "created_at", "updated_at"})
}
