package todo_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/todo"
	"github.com/hiroaki-yamamoto/todo-sample-backend/db/models/user"
)

var _ = Describe("Model", func() {
	Describe("ToGraphQL", func() {
		var u user.User
		var t todo.Todo
		var id string
		var uid string

		BeforeEach(func() {
			id = "todo-id-1"
			uid = "user-id-1"
			u = user.User{
				Id:   &uid,
				Name: "testuser",
			}
			t = todo.Todo{
				Id:   &id,
				Text: "Sample Todo",
				User: u,
			}
		})

		It("should map fields correctly when optional dates are nil", func() {
			res := t.ToGraphQL()
			Expect(res).NotTo(BeNil())
			Expect(res.ID).To(Equal(id))
			Expect(res.Text).To(Equal("Sample Todo"))
			Expect(res.WipAt).To(BeNil())
			Expect(res.CompletedAt).To(BeNil())
			Expect(res.User).NotTo(BeNil())
			Expect(res.User.ID).To(Equal(uid))
			Expect(res.User.Name).To(Equal("testuser"))
		})

		Context("with non-nil dates", func() {
			var wip, completed time.Time
			BeforeEach(func() {
				wip = time.Now().Truncate(time.Second)
				completed = wip.Add(1 * time.Hour)
				t.WipAt = &wip
				t.CompletedAt = &completed
			})

			It("should format dates into RFC3339 strings", func() {
				res := t.ToGraphQL()
				wipStr := wip.Format(time.RFC3339)
				completedStr := completed.Format(time.RFC3339)

				Expect(res.WipAt).NotTo(BeNil())
				Expect(*res.WipAt).To(Equal(wipStr))

				Expect(res.CompletedAt).NotTo(BeNil())
				Expect(*res.CompletedAt).To(Equal(completedStr))
			})
		})
	})
})
