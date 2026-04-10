package models

type Query[T any] interface {
	// SELECT * from @@table
	All() ([]T, error)
}
