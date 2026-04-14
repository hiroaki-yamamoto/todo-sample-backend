package user

import (
	"context"
	"crypto/subtle"
	"errors"

	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/user"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) Create(ctx context.Context, name string, password string) (*user.User, error) {
	var existing user.User
	if err := r.DB.WithContext(ctx).Where("name = ?", name).First(&existing).Error; err == nil {
		return nil, errors.New("user already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	u := user.New(name, password)
	if err := r.DB.WithContext(ctx).Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Authenticate(ctx context.Context, name string, password string) (*user.User, error) {
	var u user.User
	if err := r.DB.WithContext(ctx).Where(&user.User{Name: name}).First(&u).Error; err != nil {
		return nil, err
	}

	hash := argon2.IDKey([]byte(password), []byte(name), 1, 64*1024, 4, 32)

	if subtle.ConstantTimeCompare(u.Hash, hash) != 1 {
		return nil, errors.New("invalid password")
	}

	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	var u user.User
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

var _ IUserRepo = (*UserRepo)(nil)
