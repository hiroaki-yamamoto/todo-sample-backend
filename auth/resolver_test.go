package auth_test

import (
	"context"
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/hiroaki-yamamoto/gauth/config"
	gauthMw "github.com/hiroaki-yamamoto/gauth/middleware"
	"github.com/hiroaki-yamamoto/todo-sample-backend/auth"
	"github.com/hiroaki-yamamoto/todo-sample-backend/auth/model"
	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/user"
)

type mockUserRepo struct {
	CreateFunc       func(ctx context.Context, name string, password string) (*user.User, error)
	AuthenticateFunc func(ctx context.Context, name string, password string) (*user.User, error)
	GetByIDFunc      func(ctx context.Context, id string) (*user.User, error)
}

func (m *mockUserRepo) Create(ctx context.Context, name string, password string) (*user.User, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, name, password)
	}
	return nil, nil
}

func (m *mockUserRepo) Authenticate(ctx context.Context, name string, password string) (*user.User, error) {
	if m.AuthenticateFunc != nil {
		return m.AuthenticateFunc(ctx, name, password)
	}
	return nil, nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

var _ = Describe("Resolver", func() {
	var (
		resolver *auth.Resolver
		repo     *mockUserRepo
		gcfg     *config.Config
	)

	BeforeEach(func() {
		repo = &mockUserRepo{}
		gcfg = &config.Config{}
		resolver = &auth.Resolver{
			UserRepo:    repo,
			GAuthConfig: gcfg,
		}
	})

	Describe("Mutation CreateUser", func() {
		It("should create a user successfully", func() {
			repo.CreateFunc = func(ctx context.Context, name string, password string) (*user.User, error) {
				return &user.User{Name: name}, nil
			}

			res, err := resolver.Mutation().CreateUser(context.Background(), model.AuthInput{
				Name:     "testuser",
				Password: "password",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.Name).To(Equal("testuser"))
		})

		It("should propagate error on failure", func() {
			repo.CreateFunc = func(ctx context.Context, name string, password string) (*user.User, error) {
				return nil, errors.New("db error")
			}

			res, err := resolver.Mutation().CreateUser(context.Background(), model.AuthInput{
				Name:     "testuser",
				Password: "password",
			})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("db error"))
			Expect(res).To(BeNil())
		})
	})

	Describe("Mutation Login", func() {
		It("should log in user without http response writer (no cookie)", func() {
			repo.AuthenticateFunc = func(ctx context.Context, name string, password string) (*user.User, error) {
				return &user.User{Name: name}, nil
			}

			res, err := resolver.Mutation().Login(context.Background(), model.AuthInput{
				Name:     "testuser",
				Password: "password",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.Name).To(Equal("testuser"))
		})

		It("should fail on bad credentials", func() {
			repo.AuthenticateFunc = func(ctx context.Context, name string, password string) (*user.User, error) {
				return nil, errors.New("invalid credentials")
			}

			res, err := resolver.Mutation().Login(context.Background(), model.AuthInput{
				Name:     "testuser",
				Password: "wrongpassword",
			})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid credentials"))
			Expect(res).To(BeNil())
		})
	})

	Describe("Query Me", func() {
		It("should return unauthenticated without valid context", func() {
			res, err := resolver.Query().Me(context.Background())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unauthenticated"))
			Expect(res).To(BeNil())
		})

		It("should return invalid user context when wrong type in context", func() {
			r := (&http.Request{}).WithContext(context.Background())
			r = gauthMw.SetUser(r, "some-string-instead-of-user")
			res, err := resolver.Query().Me(r.Context())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid user context"))
			Expect(res).To(BeNil())
		})

		It("should return authenticated user when valid user in context", func() {
			u := &user.User{Name: "validuser"}
			r := (&http.Request{}).WithContext(context.Background())
			r = gauthMw.SetUser(r, u)
			res, err := resolver.Query().Me(r.Context())
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.Name).To(Equal("validuser"))
		})
	})
})
