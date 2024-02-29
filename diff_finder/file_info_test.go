package diff_finder_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	df "file-diff-finder/diff_finder"
)

var _ = Describe("file informative actions", func() {
	var originalFile = df.NewFileInfo(
		"abcabcabcabcabcabcabc",
		13,
	)

	It("should return version", func() {
		Expect(originalFile.Version()).To(BeEquivalentTo(13))
	})

	It("should return content", func() {
		Expect(originalFile.Content()).To(BeEquivalentTo("abcabcabcabcabcabcabc"))
	})

	It("should validate newer version number", func() {
		Expect(originalFile.ValidateVersion(14)).To(BeNil())
	})

	It("should not validate older version number for update", func() {
		err := originalFile.ValidateVersion(12)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(BeEquivalentTo("newer version should be used. Current version is 13"))
	})

	It("should not validate unexpected next version number", func() {
		err := originalFile.ValidateVersion(16)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(BeEquivalentTo("the latest version is 13. Please use incremental number for versioning of file info."))
	})
})
