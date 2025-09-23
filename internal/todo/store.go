package todo

import (
	"context"

	"github.com/syahidfrd/go-boilerplate/internal/pkg/db"
	"gorm.io/gorm"
)

// store implements todo data persistence using GORM
type store struct {
	dbConn *gorm.DB
}

// NewStore creates a new todo store with the provided database connection
func NewStore(dbConn *gorm.DB) *store {
	return &store{dbConn: dbConn}
}

// Save persists a todo to the database (create or update)
func (s *store) Save(ctx context.Context, todo *Todo, options ...db.Option) error {
	dbConn := s.dbConn

	opts := &db.Options{}
	for _, opt := range options {
		opt(opts)
	}

	if opts.Tx != nil {
		dbConn = opts.Tx
	}

	return dbConn.WithContext(ctx).Save(todo).Error
}

// GetByID retrieves a todo by its ID from the database
func (s *store) GetByID(ctx context.Context, id int64) (*Todo, error) {
	var todo Todo
	if err := s.dbConn.WithContext(ctx).First(&todo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTodoNotFound
		}
		return nil, err
	}
	return &todo, nil
}

// GetByUserID retrieves all todos for a specific user from the database
func (s *store) GetByUserID(ctx context.Context, userID int64) ([]Todo, error) {
	var todos []Todo
	if err := s.dbConn.WithContext(ctx).Where("user_id = ?", userID).Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

// Delete removes a todo from the database by its ID
func (s *store) Delete(ctx context.Context, id int64) error {
	return s.dbConn.WithContext(ctx).Delete(&Todo{}, id).Error
}
