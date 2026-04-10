package todo

import (
	"context"

	gqlModel "github.com/hiroaki-yamamoto/todo-sample-backend/graph/model"
)

type IList interface {
	List(ctx context.Context) ([]*gqlModel.Todo, error)
}

type ICreate interface {
	Create(ctx context.Context, input gqlModel.NewTodo) (*gqlModel.Todo, error)
}

type IUpdate interface {
	Update(ctx context.Context, input gqlModel.UpdateTodo) (*gqlModel.Todo, error)
}

type ITodoRepo interface {
	IList
	ICreate
	IUpdate
}
