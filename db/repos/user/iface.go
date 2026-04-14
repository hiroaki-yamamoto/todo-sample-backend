package user

import (
	"context"

	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/user"
)

type ICreate interface {
	Create(ctx context.Context, name string, password string) (*user.User, error)
}

type IAuthenticate interface {
	Authenticate(ctx context.Context, name string, password string) (*user.User, error)
}

type IGetByID interface {
	GetByID(ctx context.Context, id string) (*user.User, error)
}

type IUserRepo interface {
	ICreate
	IAuthenticate
	IGetByID
}
