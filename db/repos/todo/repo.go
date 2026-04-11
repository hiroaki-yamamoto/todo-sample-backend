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

func (r *TodoRepo) List(ctx context.Context) ([]*gqlModel.Todo, error) {
	var todos []dbtodo.Todo
	if err := r.DB.WithContext(ctx).Preload("User").Find(&todos).Error; err != nil {
		return nil, err
	}

	var result []*gqlModel.Todo
	for _, t := range todos {
		t := t // capture range variable
		result = append(result, t.ToGraphQL())
	}
	return result, nil
}

func (r *TodoRepo) Create(ctx context.Context, user user.User, input gqlModel.NewTodo) (*gqlModel.Todo, error) {
	t := dbtodo.New(input.Text, user)
	if err := r.DB.WithContext(ctx).Create(&t).Error; err != nil {
		return nil, err
	}
	return t.ToGraphQL(), nil
}

func (r *TodoRepo) Update(ctx context.Context, input gqlModel.UpdateTodo) (*gqlModel.Todo, error) {
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

	return t.ToGraphQL(), nil
}
