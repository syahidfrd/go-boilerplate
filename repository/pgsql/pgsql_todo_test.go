package pgsql_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syahidfrd/go-boilerplate/domain"
	"github.com/syahidfrd/go-boilerplate/repository/pgsql"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestTodoRepo_Create(t *testing.T) {
	todo := &domain.Todo{
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	query := "INSERT INTO todos"
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(todo.Name, todo.CreatedAt, todo.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	todoRepo := pgsql.NewPgsqlTodoRepository(db)
	err = todoRepo.Create(context.TODO(), todo)
	assert.NoError(t, err)
}

func TestTodoRepo_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	todoMock := domain.Todo{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(todoMock.ID, todoMock.Name, todoMock.CreatedAt, todoMock.UpdatedAt)

	mock.ExpectBegin()

	query := "SELECT id, name, created_at, updated_at FROM todos WHERE id = $1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(1).
		WillReturnRows(rows)

	mock.ExpectCommit()

	todoRepo := pgsql.NewPgsqlTodoRepository(db)
	todo, err := todoRepo.GetByID(context.TODO(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, todoMock.ID, todo.ID)
}

func TestTodoRepo_Fetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockTodos := []domain.Todo{
		{ID: 1, Name: "name", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Name: "name 2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(mockTodos[0].ID, mockTodos[0].Name, mockTodos[0].CreatedAt, mockTodos[0].UpdatedAt).
		AddRow(mockTodos[1].ID, mockTodos[1].Name, mockTodos[1].CreatedAt, mockTodos[1].UpdatedAt)

	mock.ExpectBegin()

	query := "SELECT id, name, created_at, updated_at FROM todos"
	mock.ExpectQuery(query).WillReturnRows(rows)

	mock.ExpectCommit()

	todoRepo := pgsql.NewPgsqlTodoRepository(db)
	todos, err := todoRepo.Fetch(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, todos, 2)
}

func TestTodoRepo_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	todo := &domain.Todo{
		ID:        1,
		Name:      "name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()

	query := "UPDATE todos SET name = $1, updated_at = $2 WHERE id = $3"
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(todo.Name, todo.UpdatedAt, todo.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	todoRepo := pgsql.NewPgsqlTodoRepository(db)
	err = todoRepo.Update(context.TODO(), todo)
	assert.NoError(t, err)
}

func TestTodoRepo_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	query := "DELETE FROM todos WHERE id = $1"
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	todoRepo := pgsql.NewPgsqlTodoRepository(db)
	err = todoRepo.Delete(context.TODO(), 1)
	assert.NoError(t, err)
}
