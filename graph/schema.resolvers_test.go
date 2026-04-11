package graph_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/user"
	"github.com/hiroaki-yamamoto/todo-sample-backend/db/repos/todo"
	"github.com/hiroaki-yamamoto/todo-sample-backend/graph"
	"github.com/hiroaki-yamamoto/todo-sample-backend/graph/model"
)

var _ = Describe("Schema.Resolvers", func() {
	var (
		ctrl     *gomock.Controller
		mockRepo *todo.MockITodoRepo
		resolver *graph.Resolver
		ctx      context.Context
		usr      user.User
	)

	BeforeEach(func() {
		usr = user.New("testuser", "password")
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = todo.NewMockITodoRepo(ctrl)
		resolver = graph.NewResolver(usr, mockRepo)
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("MutationResolver", func() {
		var mutResolver graph.MutationResolver

		BeforeEach(func() {
			mutResolver = resolver.Mutation()
		})

		Describe("CreateTodo", func() {
			It("returns the created todo", func() {
				input := model.NewTodo{Text: "Test Todo"}
				expectedTodo := &model.Todo{ID: "todo123", Text: "Test Todo"}

				mockRepo.EXPECT().Create(ctx, usr, input).Return(expectedTodo, nil)

				res, err := mutResolver.CreateTodo(ctx, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(expectedTodo))
			})

			It("returns error when repo fails", func() {
				input := model.NewTodo{Text: "Test Todo"}
				expectedErr := errors.New("creation failed")

				mockRepo.EXPECT().Create(ctx, usr, input).Return(nil, expectedErr)

				res, err := mutResolver.CreateTodo(ctx, input)
				Expect(err).To(MatchError(expectedErr))
				Expect(res).To(BeNil())
			})
		})

		Describe("UpdateTodo", func() {
			It("returns the updated todo", func() {
				input := model.UpdateTodo{ID: "todo123"}
				expectedTodo := &model.Todo{ID: "todo123"}

				mockRepo.EXPECT().Update(ctx, input).Return(expectedTodo, nil)

				res, err := mutResolver.UpdateTodo(ctx, input)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(expectedTodo))
			})

			It("returns error when repo fails", func() {
				input := model.UpdateTodo{ID: "todo123"}
				expectedErr := errors.New("update failed")

				mockRepo.EXPECT().Update(ctx, input).Return(nil, expectedErr)

				res, err := mutResolver.UpdateTodo(ctx, input)
				Expect(err).To(MatchError(expectedErr))
				Expect(res).To(BeNil())
			})
		})
	})

	Describe("QueryResolver", func() {
		var queryResolver graph.QueryResolver

		BeforeEach(func() {
			queryResolver = resolver.Query()
		})

		Describe("Todos", func() {
			It("returns a list of todos", func() {
				expectedTodos := []*model.Todo{
					{ID: "todo1"},
					{ID: "todo2"},
				}

				mockRepo.EXPECT().List(ctx).Return(expectedTodos, nil)

				res, err := queryResolver.Todos(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(Equal(expectedTodos))
			})

			It("returns error when repo fails", func() {
				expectedErr := errors.New("list failed")

				mockRepo.EXPECT().List(ctx).Return(nil, expectedErr)

				res, err := queryResolver.Todos(ctx)
				Expect(err).To(MatchError(expectedErr))
				Expect(res).To(BeNil())
			})
		})
	})
})
