package cache_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("cache integration", func() {
	AfterEach(func() {
		defer localCache.Clear()
	})

	It("should save new key to cache and get it from cache", func() {
		localCache.Set("test", 100)

		val, exists := localCache.Get("test")
		Expect(exists).To(BeTrue())
		Expect(val).To(BeEquivalentTo(100))
	})

	It("should save new key and delete it", func() {
		localCache.Set("test", 100)
		localCache.Delete("test")

		val, exists := localCache.Get("test")
		Expect(exists).To(BeFalse())
		Expect(val).To(BeNil())
	})

	It("should clear all keys", func() {
		localCache.Set("test", 100)
		localCache.Set("test2", 200)
		localCache.Clear()

		val, exists := localCache.Get("test")
		Expect(exists).To(BeFalse())
		Expect(val).To(BeNil())

		val, exists = localCache.Get("test")
		Expect(exists).To(BeFalse())
		Expect(val).To(BeNil())
	})
})
