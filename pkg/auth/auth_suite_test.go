package auth_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMongoose(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongoose Auth Suite")
}
