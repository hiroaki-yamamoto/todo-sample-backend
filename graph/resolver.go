package graph

import (
	"github.com/hiroaki-yamamoto/todo-sample-backend/db/repos/todo"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	todoRepo todo.ITodoRepo
}

func NewResolver(todoRepo todo.ITodoRepo) *Resolver {
	return &Resolver{todoRepo: todoRepo}
}
