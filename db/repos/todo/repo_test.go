package todo_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	dbtodo "github.com/hiroaki-yamamoto/todo-sample-backend/db/models/todo"
	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/user"
	"github.com/hiroaki-yamamoto/todo-sample-backend/db/repos/todo"
	gqlModel "github.com/hiroaki-yamamoto/todo-sample-backend/graph/model"
)

var _ = Describe("Repo", func() {
	var db *gorm.DB
	var repo *todo.TodoRepo
	var ctx context.Context

	BeforeEach(func() {
		dsn := "host=localhost user=postgres password=password dbname=todo_test port=5432 sslmode=disable"
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())

		repo = todo.NewRepo(db)
		ctx = context.Background()
	})

	AfterEach(func() {
		// Clean up the todos table after each test
		err := db.Exec("TRUNCATE TABLE todos, users CASCADE").Error
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Create", func() {
		var u user.User

		BeforeEach(func() {
			u = user.New("testuser", "password")
			err := db.Create(&u).Error
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create a new todo", func() {
			input := gqlModel.NewTodo{
				Text: "My new task",
			}
			result, err := repo.Create(ctx, u, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(result.Text).To(Equal("My new task"))
			Expect(*result.UserId).To(Equal(*u.Id))
			Expect(result.WipAt).To(BeNil())
			Expect(result.CompletedAt).To(BeNil())

			// Verify in DB
			var count int64
			db.Model(&dbtodo.Todo{}).Count(&count)
			Expect(count).To(Equal(int64(1)))
		})

		It("should return an error if user does not exist", func() {
			input := gqlModel.NewTodo{
				Text: "My new task",
			}
			ghost_uid := "00000000-0000-0000-0000-000000000000"
			user := user.User{Id: &ghost_uid, Name: "ghost", Hash: nil}
			_, err := repo.Create(ctx, user, input)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("List", func() {
		var u user.User

		BeforeEach(func() {
			u = user.New("testuser2", "password")
			err := db.Create(&u).Error
			Expect(err).NotTo(HaveOccurred())

			t1 := dbtodo.New("Task 1", u)
			t2 := dbtodo.New("Task 2", u)
			err = db.Create(&t1).Error
			Expect(err).NotTo(HaveOccurred())
			err = db.Create(&t2).Error
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return all todos with user preloaded", func() {
			results, err := repo.List(ctx, u)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(2))
			// Order may vary slightly depending on primary key, so we check existence
			texts := []string{results[0].Text, results[1].Text}
			Expect(texts).To(ContainElement("Task 1"))
			Expect(texts).To(ContainElement("Task 2"))
			Expect(*results[0].UserId).To(Equal(*u.Id))
		})
	})

	Describe("Update", func() {
		var u user.User
		var t dbtodo.Todo

		BeforeEach(func() {
			u = user.New("testuser3", "password")
			err := db.Create(&u).Error
			Expect(err).NotTo(HaveOccurred())

			t = dbtodo.New("Update Task", u)
			err = db.Create(&t).Error
			Expect(err).NotTo(HaveOccurred())
		})

		It("should update the text of the todo", func() {
			input := gqlModel.UpdateTodo{
				ID:   *t.Id,
				Text: "Updated Task",
			}
			result, err := repo.Update(ctx, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(result.Text).To(Equal("Updated Task"))

			var dbTodo dbtodo.Todo
			db.First(&dbTodo, "id = ?", *t.Id)
			Expect(dbTodo.Text).To(Equal("Updated Task"))
		})

		It("should update WipAt and CompletedAt fields", func() {
			wipTime := time.Now().Truncate(time.Second)
			completedTime := wipTime.Add(time.Hour)
			wipStr := wipTime.Format(time.RFC3339)
			compStr := completedTime.Format(time.RFC3339)

			input := gqlModel.UpdateTodo{
				ID:          *t.Id,
				Text:        "Status Updated",
				WipAt:       &wipStr,
				CompletedAt: &compStr,
			}
			result, err := repo.Update(ctx, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Text).To(Equal("Status Updated"))
			Expect((*result.WipAt).Unix()).To(Equal(wipTime.Unix()))
			Expect((*result.CompletedAt).Unix()).To(Equal(completedTime.Unix()))

			var dbTodo dbtodo.Todo
			db.First(&dbTodo, "id = ?", *t.Id)
			Expect(dbTodo.WipAt.Unix()).To(Equal(wipTime.Unix()))
			Expect(dbTodo.CompletedAt.Unix()).To(Equal(completedTime.Unix()))
		})

		It("should return an error for a non-existent todo", func() {
			input := gqlModel.UpdateTodo{
				ID:   "00000000-0000-0000-0000-000000000000",
				Text: "Ghost Task",
			}
			_, err := repo.Update(ctx, input)
			Expect(err).To(HaveOccurred())
		})
	})
})
