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
	todoRepo   domain.TodoRepository
	redisRepo  redis.RedisRepository
	ctxTimeout time.Duration
}

// NewTodoUsecase will create new an todoUsecase object representation of TodoUsecase interface
func NewTodoUsecase(todoRepo domain.TodoRepository, redisRepo redis.RedisRepository, ctxTimeout time.Duration) *todoUsecase {
	return &todoUsecase{
		todoRepo:   todoRepo,
		redisRepo:  redisRepo,
		ctxTimeout: ctxTimeout,
	}
}

func (u *todoUsecase) Create(c context.Context, request *request.CreateTodoReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	err = u.todoRepo.Create(ctx, &domain.Todo{
		Name:      request.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	return
}

func (u *todoUsecase) GetByID(c context.Context, id int64) (todo domain.Todo, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	todo, err = u.todoRepo.GetByID(ctx, id)
	if err != nil && err == sql.ErrNoRows {
		err = utils.NewNotFoundError("todo not found")
		return
	}
	return
}

func (u *todoUsecase) Fetch(c context.Context) (todos []domain.Todo, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	todosCached, _ := u.redisRepo.Get("todos")
	if err = json.Unmarshal([]byte(todosCached), &todos); err == nil {
		return
	}

	todos, err = u.todoRepo.Fetch(ctx)
	if err != nil {
		return
	}

	todosString, _ := json.Marshal(&todos)
	u.redisRepo.Set("todos", todosString, 30*time.Second)
	return
}

func (u *todoUsecase) Update(c context.Context, id int64, request *request.UpdateTodoReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	todo, err := u.todoRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewNotFoundError("todo not found")
			return
		}
		return
	}

	todo.Name = request.Name
	todo.UpdatedAt = time.Now()

	err = u.todoRepo.Update(ctx, &todo)
	return
}

func (u *todoUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	_, err = u.todoRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewNotFoundError("todo not found")
			return
		}
		return
	}

	err = u.todoRepo.Delete(ctx, id)
	return
}
