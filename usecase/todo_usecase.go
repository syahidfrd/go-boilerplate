package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/syahidfrd/go-boilerplate/domain"
	"github.com/syahidfrd/go-boilerplate/repository/redis"
	"github.com/syahidfrd/go-boilerplate/transport/request"
	"github.com/syahidfrd/go-boilerplate/utils"
)

type todoUsecase struct {
	db         domain.Database
	todoRepo   domain.TodoRepository
	redisRepo  redis.RedisRepository
	ctxTimeout time.Duration
}

// NewTodoUsecase will create new an todoUsecase object representation of TodoUsecase interface
func NewTodoUsecase(db domain.Database, todoRepo domain.TodoRepository, redisRepo redis.RedisRepository, ctxTimeout time.Duration) *todoUsecase {
	return &todoUsecase{
		db:         db,
		todoRepo:   todoRepo,
		redisRepo:  redisRepo,
		ctxTimeout: ctxTimeout,
	}
}

func (u *todoUsecase) Create(c context.Context, request *request.CreateTodoReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	tx, err := u.db.BeginTx(ctx)
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	err = u.todoRepo.Create(ctx, tx, &domain.Todo{
		Name:      request.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err = tx.Commit(); err != nil {
		return
	}

	return
}

func (u *todoUsecase) GetByID(c context.Context, id int64) (todo domain.Todo, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	tx, err := u.db.BeginTx(ctx)
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	todo, err = u.todoRepo.GetByID(ctx, tx, id)
	if err != nil && err == sql.ErrNoRows {
		err = utils.NewNotFoundError("todo not found")
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}

func (u *todoUsecase) Fetch(c context.Context) (todos []domain.Todo, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	tx, err := u.db.BeginTx(ctx)
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	todosCached, _ := u.redisRepo.Get("todos")
	if err = json.Unmarshal([]byte(todosCached), &todos); err == nil {
		return
	}

	todos, err = u.todoRepo.Fetch(ctx, tx)
	if err != nil {
		return
	}

	todosString, _ := json.Marshal(&todos)
	u.redisRepo.Set("todos", todosString, 30*time.Second)

	if err = tx.Commit(); err != nil {
		return
	}

	return
}

func (u *todoUsecase) Update(c context.Context, id int64, request *request.UpdateTodoReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	tx, err := u.db.BeginTx(ctx)
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	todo, err := u.todoRepo.GetByID(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewNotFoundError("todo not found")
			return
		}
		return
	}

	todo.Name = request.Name
	todo.UpdatedAt = time.Now()

	err = u.todoRepo.Update(ctx, tx, &todo)

	if err = tx.Commit(); err != nil {
		return
	}

	return
}

func (u *todoUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	tx, err := u.db.BeginTx(ctx)
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	_, err = u.todoRepo.GetByID(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewNotFoundError("todo not found")
			return
		}
		return
	}

	err = u.todoRepo.Delete(ctx, tx, id)

	if err = tx.Commit(); err != nil {
		return
	}

	return
}
