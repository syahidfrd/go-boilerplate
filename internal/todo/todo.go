package todo

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	// ErrTodoNotFound is returned when a requested todo cannot be found
	ErrTodoNotFound = errors.New("todo not found")
)

// Todo represents a todo item with user association and completion status
type Todo struct {
	ID          int64
	UserID      int64 `gorm:"index"`
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTodo creates a new todo item with the given details
func NewTodo(userID int64, title, description string) *Todo {
	now := time.Now()
	return &Todo{
		UserID:      userID,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MarkAsCompleted marks the todo as completed
func (t *Todo) MarkAsCompleted() {
	t.Completed = true
}

// MarkAsIncomplete marks the todo as incomplete
func (t *Todo) MarkAsIncomplete() {
	t.Completed = false
}

// MarshalJSON implements the json.Marshaler interface for custom JSON serialization
func (t Todo) MarshalJSON() ([]byte, error) {
	var j struct {
		ID          int64  `json:"id"`
		UserID      int64  `json:"user_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}

	j.ID = t.ID
	j.UserID = t.UserID
	j.Title = t.Title
	j.Description = t.Description
	j.Completed = t.Completed
	j.CreatedAt = t.CreatedAt.Format(time.RFC3339)
	j.UpdatedAt = t.UpdatedAt.Format(time.RFC3339)

	return json.Marshal(j)
}
