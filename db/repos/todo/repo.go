package todo

import (
	"context"
	"time"

	dbtodo "github.com/hiroaki-yamamoto/todo-sample-backend/db/models/todo"
	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/user"
	gqlModel "github.com/hiroaki-yamamoto/todo-sample-backend/graph/model"
	"gorm.io/gorm"
)

type TodoRepo struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) *TodoRepo {
	return &TodoRepo{DB: db}
}

func (r *TodoRepo) List(ctx context.Context, user user.User) ([]dbtodo.Todo, error) {
	var todos []dbtodo.Todo
	if err := r.DB.WithContext(ctx).Where("user_id = ?", user.Id).Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *TodoRepo) Create(ctx context.Context, user user.User, input gqlModel.NewTodo) (*dbtodo.Todo, error) {
	t := dbtodo.New(input.Text, user)
	if err := r.DB.WithContext(ctx).Create(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TodoRepo) Update(ctx context.Context, input gqlModel.UpdateTodo) (*dbtodo.Todo, error) {
	var t dbtodo.Todo
	if err := r.DB.WithContext(ctx).Preload("User").First(&t, "id = ?", input.ID).Error; err != nil {
		return nil, err
	}
	t.Text = input.Text
	if input.WipAt != nil {
		w, err := time.Parse(time.RFC3339, *input.WipAt)
		if err == nil {
			t.WipAt = &w
		}
	} else {
		t.WipAt = nil
	}
	if input.CompletedAt != nil {
		c, err := time.Parse(time.RFC3339, *input.CompletedAt)
		if err == nil {
			t.CompletedAt = &c
		}
	} else {
		t.CompletedAt = nil
	}
	if err := r.DB.WithContext(ctx).Save(&t).Error; err != nil {
		return nil, err
	}

	return &t, nil
}
