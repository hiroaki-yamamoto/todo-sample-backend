package user

import (
	"golang.org/x/crypto/argon2"

	authModel "github.com/hiroaki-yamamoto/todo-sample-backend/auth/model"
)

type User struct {
	Id   *string `gorm:"default:uuidv7()"`
	Name string
	Hash []byte
}

// New creates a new User instance with the given name and password.
// The Id field is initialized to nil, and the Hash field is generated using the Argon2 algorithm.
func New(name string, password string) User {
	hash := argon2.IDKey([]byte(password), []byte(name), 1, 64*1024, 4, 32)
	return User{
		Id:   nil,
		Name: name,
		Hash: hash,
	}
}

// ToGraphQL converts the User instance to a GraphQL User model.
func (me *User) ToGraphQL() *authModel.User {
	return &authModel.User{
		ID:   *me.Id,
		Name: me.Name,
	}
}
