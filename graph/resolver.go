package graph

import (
	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/user"
	"github.com/hiroaki-yamamoto/todo-sample-backend/db/repos/todo"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	todoRepo todo.ITodoRepo
	user     user.User
}

func NewResolver(usr user.User, todoRepo todo.ITodoRepo) *Resolver {
	return &Resolver{todoRepo: todoRepo, user: usr}
}
