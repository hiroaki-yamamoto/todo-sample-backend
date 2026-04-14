package user_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var _ = BeforeSuite(func() {
	var err error
	DB, err = gorm.Open(postgres.Open("postgres://postgres:password@localhost:5432/todo_test?sslmode=disable"), &gorm.Config{})
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	DB = nil
})

func TestUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Suite")
}
