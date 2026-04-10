package todo

import (
	"time"

	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/user"

	graph "github.com/hiroaki-yamamoto/todo-sample-backend/graph/model"
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

func (me *Todo) ToGraphQL() *graph.Todo {
	var wipAt, completedAt *string
	if me.WipAt != nil {
		w := me.WipAt.Format(time.RFC3339)
		wipAt = &w
	}
	if me.CompletedAt != nil {
		c := me.CompletedAt.Format(time.RFC3339)
		completedAt = &c
	}
	return &graph.Todo{
		ID:          *me.Id,
		Text:        me.Text,
		WipAt:       wipAt,
		CompletedAt: completedAt,
		User:        me.User.ToGraphQL(),
	}
}
