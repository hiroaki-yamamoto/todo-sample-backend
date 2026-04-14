package auth

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"

	"github.com/hiroaki-yamamoto/todo-sample-backend/auth/model"
	"github.com/hiroaki-yamamoto/todo-sample-backend/db/repos/user"
)

type Resolver struct {
	UserRepo user.IUserRepo
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.AuthInput) (*model.User, error) {
	panic("not implemented")
}

// CrateUser is the resolver for the crateUser field.
func (r *mutationResolver) CrateUser(ctx context.Context, input model.AuthInput) (*model.User, error) {
	panic("not implemented")
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	panic("not implemented")
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
