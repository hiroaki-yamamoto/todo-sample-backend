package todo

import (
	"context"

	model "github.com/hiroaki-yamamoto/todo-sample-backend/db/models/todo"
	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/user"
	gqlModel "github.com/hiroaki-yamamoto/todo-sample-backend/graph/model"
)

type IList interface {
	List(ctx context.Context, user user.User) ([]model.Todo, error)
}

type ICreate interface {
	Create(ctx context.Context, user user.User, input gqlModel.NewTodo) (*model.Todo, error)
}

type IUpdate interface {
	Update(ctx context.Context, user user.User, input gqlModel.UpdateTodo) (*model.Todo, error)
}

type ITodoRepo interface {
	IList
	ICreate
	IUpdate
}
