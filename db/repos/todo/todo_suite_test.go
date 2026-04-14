package todo_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func TestTodo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Todo Suite")
}
