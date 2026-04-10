package user

import "golang.org/x/crypto/argon2"

type User struct {
	Id   *string
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
