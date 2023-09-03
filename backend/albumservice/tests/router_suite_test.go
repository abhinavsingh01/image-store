package tests

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAlbum(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Album Suite")
}
