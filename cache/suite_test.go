package cache_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"file-diff-finder/cache"
)

var (
	localCache cache.ICache
)

func TestConsumer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cache test suite")
}

var _ = BeforeSuite(func() {
	localCache = cache.InitCache()
})
