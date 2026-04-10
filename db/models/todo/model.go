package todo

import (
	"time"

	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/user"
)

type Todo struct {
	Id          *string
	Text        string
	WipAt       *time.Time
	CompletedAt *time.Time
	User        user.User
}

// New creates a new Todo instance with the given text and user.
// The Id, WipAt, and CompletedAt fields are initialized to nil.
func New(text string, user user.User) Todo {
	return Todo{
		Id:          nil,
		Text:        text,
		WipAt:       nil,
		CompletedAt: nil,
		User:        user,
	}
}
