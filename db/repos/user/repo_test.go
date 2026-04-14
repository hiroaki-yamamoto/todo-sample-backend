package user_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"github.com/hiroaki-yamamoto/todo-sample-backend/db/repos/user"
)

var _ = Describe("UserRepo", func() {
	var ctx context.Context
	var repo *user.UserRepo

	BeforeEach(func() {
		repo = user.NewRepo(DB)
		ctx = context.Background()
		DB.Exec("TRUNCATE TABLE users CASCADE")
	})

	Describe("Create", func() {
		It("should create a new user successfully", func() {
			u, err := repo.Create(ctx, "testuser", "password123")
			Expect(err).NotTo(HaveOccurred())
			Expect(u).NotTo(BeNil())
			Expect(u.Id).NotTo(BeZero())
			Expect(u.Name).To(Equal("testuser"))
			Expect(u.Hash).NotTo(BeEmpty())
		})

		It("should fail if the user already exists", func() {
			_, err := repo.Create(ctx, "testuser", "password123")
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.Create(ctx, "testuser", "password456")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Authenticate", func() {
		BeforeEach(func() {
			_, err := repo.Create(ctx, "testuser", "password123")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should authenticate successfully with correct credentials", func() {
			u, err := repo.Authenticate(ctx, "testuser", "password123")
			Expect(err).NotTo(HaveOccurred())
			Expect(u).NotTo(BeNil())
			Expect(u.Name).To(Equal("testuser"))
		})

		It("should fail with incorrect password", func() {
			_, err := repo.Authenticate(ctx, "testuser", "wrongpassword")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid password"))
		})

		It("should fail with non-existent user", func() {
			_, err := repo.Authenticate(ctx, "nonexistent", "password123")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(gorm.ErrRecordNotFound))
		})
	})

	Describe("GetByID", func() {
		var createdUserID string

		BeforeEach(func() {
			u, err := repo.Create(ctx, "testuser", "password123")
			Expect(err).NotTo(HaveOccurred())
			createdUserID = *u.Id
		})

		It("should return user for correct ID", func() {
			u, err := repo.GetByID(ctx, createdUserID)
			Expect(err).NotTo(HaveOccurred())
			Expect(u).NotTo(BeNil())
			Expect(*u.Id).To(Equal(createdUserID))
			Expect(u.Name).To(Equal("testuser"))
		})

		It("should fail for non-existent ID", func() {
			_, err := repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(gorm.ErrRecordNotFound))
		})
	})
})
