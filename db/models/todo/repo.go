package todo

import "github.com/hiroaki-yamamoto/todo-sample-backend/db/models"

type Repo interface {
	models.Query[Todo]
	// UPDATE @@table SET text=@text WHERE id=@id
	UpdateText(id string, text string) (*Todo, error)
	// UPDATE @@table SET wip_at=now() WHERE id=@id
	StartWip(id string) (*Todo, error)
	// UPDATE @@table SET completed_at=now() WHERE id=@id
	Complete(id string) (*Todo, error)
}
