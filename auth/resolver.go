package auth

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"
	"errors"

	"github.com/hiroaki-yamamoto/gauth/config"
	"github.com/hiroaki-yamamoto/gauth/core"
	gauthMw "github.com/hiroaki-yamamoto/gauth/middleware"
	"github.com/hiroaki-yamamoto/todo-sample-backend/auth/model"
	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/user"
	userRepo "github.com/hiroaki-yamamoto/todo-sample-backend/db/repos/user"
)

type Resolver struct {
	UserRepo    userRepo.IUserRepo
	GAuthConfig *config.Config
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.AuthInput) (*model.User, error) {
	u, err := r.UserRepo.Authenticate(ctx, input.Name, input.Password)
	if err != nil {
		return nil, err
	}
	w := GetResponseWriter(ctx)
	if w != nil {
		err = core.Login(w, r.GAuthConfig, u)
		if err != nil {
			return nil, err
		}
	}
	return u.ToGraphQL(), nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.AuthInput) (*model.User, error) {
	u, err := r.UserRepo.Create(ctx, input.Name, input.Password)
	if err != nil {
		return nil, err
	}
	return u.ToGraphQL(), nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	v := gauthMw.GetUser(ctx)
	if v == nil {
		return nil, errors.New("unauthenticated")
	}
	u, ok := v.(*user.User)
	if !ok {
		return nil, errors.New("invalid user context")
	}
	return u.ToGraphQL(), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	type Resolver struct{}
*/
