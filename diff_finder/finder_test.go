package diff_finder_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	df "file-diff-finder/diff_finder"
	"file-diff-finder/diff_finder/mocks"
)

var _ = Describe("file context difference finder cases", func() {
	var (
		mockCtrl       *gomock.Controller
		originalFile   *mocks.MockIFileInformative
		fileDiffFinder df.IFileDiffFinder
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		originalFile = mocks.NewMockIFileInformative(mockCtrl)
		fileDiffFinder = df.NewFileDiffFinder(originalFile)
	})

	Context("non timed out cases", func() {
		BeforeEach(func() {
			originalFile.EXPECT().Content().Return("abcabcabcabcabcabcabc").AnyTimes()
		})

		Context("updated file and original file has same length", func() {
			It("should return non-empty delta info", func() {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				affectedIndeces, err := fileDiffFinder.Diff(ctx, "1bcabcabca6cabcab4abc")

				Expect(err).To(BeNil())
				Expect(affectedIndeces[0].OldValue).To(BeEquivalentTo("a"))
				Expect(affectedIndeces[0].NewValue).To(BeEquivalentTo("1"))
				Expect(affectedIndeces[0].Index).To(BeEquivalentTo(0))
				Expect(affectedIndeces[0].Type).To(BeEquivalentTo(df.ChangeTypesUpdated))
				Expect(affectedIndeces[1].OldValue).To(BeEquivalentTo("b"))
				Expect(affectedIndeces[1].NewValue).To(BeEquivalentTo("6"))
				Expect(affectedIndeces[1].Index).To(BeEquivalentTo(10))
				Expect(affectedIndeces[1].Type).To(BeEquivalentTo(df.ChangeTypesUpdated))
				Expect(affectedIndeces[2].OldValue).To(BeEquivalentTo("c"))
				Expect(affectedIndeces[2].NewValue).To(BeEquivalentTo("4"))
				Expect(affectedIndeces[2].Index).To(BeEquivalentTo(17))
				Expect(affectedIndeces[2].Type).To(BeEquivalentTo(df.ChangeTypesUpdated))
			})

			It("should return empty delta info for same file", func() {
				affectedIndeces, err := fileDiffFinder.Diff(context.Background(), "abcabcabcabcabcabcabc")

				Expect(err).To(BeNil())
				Expect(len(affectedIndeces)).To(BeZero())
			})
		})

		Context("updated file has more character set than the original file", func() {
			It("should return non-empty delta info", func() {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				affectedIndeces, err := fileDiffFinder.Diff(ctx, "abca2cabcabcabcabcabc5f")

				Expect(err).To(BeNil())
				Expect(affectedIndeces[0].OldValue).To(BeEquivalentTo("b"))
				Expect(affectedIndeces[0].NewValue).To(BeEquivalentTo("2"))
				Expect(affectedIndeces[0].Index).To(BeEquivalentTo(4))
				Expect(affectedIndeces[0].Type).To(BeEquivalentTo(df.ChangeTypesUpdated))
				Expect(affectedIndeces[1].OldValue).To(BeEquivalentTo(""))
				Expect(affectedIndeces[1].NewValue).To(BeEquivalentTo("5"))
				Expect(affectedIndeces[1].Index).To(BeEquivalentTo(21))
				Expect(affectedIndeces[1].Type).To(BeEquivalentTo(df.ChangeTypesAdded))
				Expect(affectedIndeces[2].OldValue).To(BeEquivalentTo(""))
				Expect(affectedIndeces[2].NewValue).To(BeEquivalentTo("f"))
				Expect(affectedIndeces[2].Index).To(BeEquivalentTo(22))
				Expect(affectedIndeces[2].Type).To(BeEquivalentTo(df.ChangeTypesAdded))
			})
		})

		Context("updated file has less character set than the original file", func() {
			It("should return non-empty delta info", func() {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				affectedIndeces, err := fileDiffFinder.Diff(ctx, "abcabcabcaxcabcabca")

				Expect(err).To(BeNil())
				Expect(affectedIndeces[0].OldValue).To(BeEquivalentTo("b"))
				Expect(affectedIndeces[0].NewValue).To(BeEquivalentTo("x"))
				Expect(affectedIndeces[0].Index).To(BeEquivalentTo(10))
				Expect(affectedIndeces[0].Type).To(BeEquivalentTo(df.ChangeTypesUpdated))
				Expect(affectedIndeces[1].OldValue).To(BeEquivalentTo("b"))
				Expect(affectedIndeces[1].NewValue).To(BeEquivalentTo(""))
				Expect(affectedIndeces[1].Index).To(BeEquivalentTo(19))
				Expect(affectedIndeces[1].Type).To(BeEquivalentTo(df.ChangeTypesRemoved))
				Expect(affectedIndeces[2].OldValue).To(BeEquivalentTo("c"))
				Expect(affectedIndeces[2].NewValue).To(BeEquivalentTo(""))
				Expect(affectedIndeces[2].Index).To(BeEquivalentTo(20))
				Expect(affectedIndeces[2].Type).To(BeEquivalentTo(df.ChangeTypesRemoved))
			})
		})
	})

	Context("content timed out", func() {
		It("should return error", func() {
			originalFile.EXPECT().Content().DoAndReturn(func() string {
				time.Sleep(2 * time.Second)
				return "abcabcabcabcabcabcabc"
			})

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			affectedIndeces, err := fileDiffFinder.Diff(ctx, "abcabcabcaxcabcabca")

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(BeEquivalentTo("context deadline exceeded"))
			Expect(len(affectedIndeces)).To(BeZero())
		})
	})
})
