package todo

import (
	"context"

	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/user"
	gqlModel "github.com/hiroaki-yamamoto/todo-sample-backend/graph/model"
)

type IList interface {
	List(ctx context.Context) ([]*gqlModel.Todo, error)
}

type ICreate interface {
	Create(ctx context.Context, user user.User, input gqlModel.NewTodo) (*gqlModel.Todo, error)
}

type IUpdate interface {
	Update(ctx context.Context, input gqlModel.UpdateTodo) (*gqlModel.Todo, error)
}

type ITodoRepo interface {
	IList
	ICreate
	IUpdate
}
