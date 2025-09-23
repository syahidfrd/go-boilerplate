package todo

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTodo(t *testing.T) {
	userID := int64(123)
	title := "Test Todo"
	description := "Test Description"

	todo := NewTodo(userID, title, description)

	assert.Equal(t, userID, todo.UserID)
	assert.Equal(t, title, todo.Title)
	assert.Equal(t, description, todo.Description)
	assert.False(t, todo.Completed)
	assert.WithinDuration(t, time.Now(), todo.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), todo.UpdatedAt, time.Second)
}

func TestTodo_MarkAsCompleted(t *testing.T) {
	todo := NewTodo(123, "Test", "Description")

	todo.MarkAsCompleted()

	assert.True(t, todo.Completed)
}

func TestTodo_MarkAsIncomplete(t *testing.T) {
	todo := NewTodo(123, "Test", "Description")
	todo.Completed = true

	todo.MarkAsIncomplete()

	assert.False(t, todo.Completed)
}

func TestTodo_MarkAsCompleted_AlreadyCompleted(t *testing.T) {
	todo := NewTodo(123, "Test", "Description")
	todo.Completed = true

	todo.MarkAsCompleted()

	assert.True(t, todo.Completed)
}

func TestTodo_MarkAsIncomplete_AlreadyIncomplete(t *testing.T) {
	todo := NewTodo(123, "Test", "Description")
	todo.Completed = false

	todo.MarkAsIncomplete()

	assert.False(t, todo.Completed)
}

func TestTodo_MarshalJSON(t *testing.T) {
	now := time.Now()
	todo := &Todo{
		ID:          1,
		UserID:      123,
		Title:       "Test Todo",
		Description: "Test Description",
		Completed:   true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	jsonBytes, err := json.Marshal(todo)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), result["id"])
	assert.Equal(t, float64(123), result["user_id"])
	assert.Equal(t, "Test Todo", result["title"])
	assert.Equal(t, "Test Description", result["description"])
	assert.Equal(t, true, result["completed"])
	assert.Equal(t, now.Format(time.RFC3339), result["created_at"])
	assert.Equal(t, now.Format(time.RFC3339), result["updated_at"])
}

func TestTodo_MarshalJSON_EmptyDescription(t *testing.T) {
	todo := &Todo{
		ID:          1,
		UserID:      123,
		Title:       "Test Todo",
		Description: "",
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	jsonBytes, err := json.Marshal(todo)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	assert.NoError(t, err)

	assert.Equal(t, "", result["description"])
	assert.Equal(t, false, result["completed"])
}

func TestTodo_MarshalJSON_TimeFormatting(t *testing.T) {
	// Test with specific time to verify RFC3339 formatting
	specificTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	todo := &Todo{
		ID:        1,
		UserID:    123,
		Title:     "Test",
		CreatedAt: specificTime,
		UpdatedAt: specificTime,
	}

	jsonBytes, err := json.Marshal(todo)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	assert.NoError(t, err)

	expectedTime := "2023-12-25T15:30:45Z"
	assert.Equal(t, expectedTime, result["created_at"])
	assert.Equal(t, expectedTime, result["updated_at"])
}
