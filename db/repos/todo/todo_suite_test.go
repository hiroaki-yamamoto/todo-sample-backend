package todo_test

import (
	"database/sql"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var _ = BeforeSuite(func() {
	dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable"
	sql, err := sql.Open("postgres", dsn)
	Expect(err).NotTo(HaveOccurred())
	_, err = sql.Query("DROP DATABASE IF EXISTS test WITH (FORCE)")
	Expect(err).NotTo(HaveOccurred())
	_, err = sql.Query("CREATE DATABASE test")
	Expect(err).NotTo(HaveOccurred())

	Expect(err).NotTo(HaveOccurred())
	m, err := migrate.New(
		"file://../../migrations",
		"postgres://postgres:password@localhost:5432/test?sslmode=disable")
	Expect(err).NotTo(HaveOccurred())
	err = m.Up()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable"
	sql, err := sql.Open("postgres", dsn)
	Expect(err).NotTo(HaveOccurred())
	m, err := migrate.New(
		"file://../../migrations",
		"postgres://postgres:password@localhost:5432/test?sslmode=disable",
	)
	Expect(err).NotTo(HaveOccurred())
	err = m.Down()
	Expect(err).NotTo(HaveOccurred())

	_, err = sql.Query("DROP DATABASE IF EXISTS test WITH (FORCE)")
	Expect(err).NotTo(HaveOccurred())
})

func TestTodo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Todo Suite")
}
