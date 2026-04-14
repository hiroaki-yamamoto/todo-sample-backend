package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/hiroaki-yamamoto/todo-sample-backend/auth"
)

var _ = Describe("Middleware", func() {
	It("injects response writer", func() {
		called := false
		handler := auth.InjectResponseWriter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			gotW := auth.GetResponseWriter(r.Context())
			Expect(gotW).To(Equal(w))
		}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		handler.ServeHTTP(w, r)
		Expect(called).To(BeTrue())
	})

	It("returns nil if response writer is not injected", func() {
		gotW := auth.GetResponseWriter(context.Background())
		Expect(gotW).To(BeNil())
	})
})
