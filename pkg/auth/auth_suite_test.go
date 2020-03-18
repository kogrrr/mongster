package auth_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMongster(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongster Auth Suite")
}
