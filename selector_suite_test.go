package selector_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSelector(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Selector Suite")
}
