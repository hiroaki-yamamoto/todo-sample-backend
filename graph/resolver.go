package graph

import (
	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/user"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	DB   *gorm.DB
	user user.User
}

func NewResolver(db *gorm.DB, usr user.User) *Resolver {
	return &Resolver{DB: db, user: usr}
}
