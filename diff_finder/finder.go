package diff_finder

import (
	"context"
	"fmt"
)

const (
	base              = 256
	mod               = 1000000007
	widerWindowRange  = 8
	narrowWindowRange = 2
)

type IFileDiffFinder interface {
	Diff(ctx context.Context, updatedFileText string) ([]UpdatedIndex, error)
	IFileInformative
}

type FileDiffFinder struct {
	originalFileInfo IFileInformative
}

func NewFileDiffFinder(fileInfo IFileInformative) IFileDiffFinder {
	return &FileDiffFinder{
		originalFileInfo: fileInfo,
	}
}

func (rhs *FileDiffFinder) Diff(ctx context.Context, updatedFileText string) ([]UpdatedIndex, error) {
	c := make(chan struct{})
	var comparedIndexes []UpdatedIndex
	go func() {
		defer close(c)
		if updatedFileText == rhs.originalFileInfo.Content() {
			// do nothing
		} else if len(updatedFileText) < widerWindowRange || len(rhs.originalFileInfo.Content()) < widerWindowRange {
			comparedIndexes = rhs.compareBruteForce(updatedFileText)
		} else if len(updatedFileText) == len(rhs.originalFileInfo.Content()) {
			ranges := rhs.findMinimalRanges(rhs.originalFileInfo.Content(), updatedFileText)
			comparedIndexes = rhs.compareIndexes(ranges, updatedFileText)
		} else if len(updatedFileText) > len(rhs.originalFileInfo.Content()) {
			ranges := rhs.findMinimalRanges(rhs.originalFileInfo.Content(), updatedFileText[:len(rhs.originalFileInfo.Content())])
			comparedIndexes = rhs.compareIndexes(ranges, updatedFileText[:len(rhs.originalFileInfo.Content())])

			for ind := len(rhs.originalFileInfo.Content()); ind < len(updatedFileText); ind++ {
				comparedIndexes = append(comparedIndexes, UpdatedIndex{
					NewValue: fmt.Sprintf("%c", updatedFileText[ind]),
					Index:    ind,
					Type:     ChangeTypesAdded,
				})
			}
		} else {
			ranges := rhs.findMinimalRanges(rhs.originalFileInfo.Content()[:len(updatedFileText)], updatedFileText)
			comparedIndexes = rhs.compareIndexes(ranges, updatedFileText)

			for ind := len(updatedFileText); ind < len(rhs.originalFileInfo.Content()); ind++ {
				comparedIndexes = append(comparedIndexes, UpdatedIndex{
					OldValue: fmt.Sprintf("%c", rhs.originalFileInfo.Content()[ind]),
					Index:    ind,
					Type:     ChangeTypesRemoved,
				})
			}
		}
	}()

	select {
	case <-c:
		// completed in time
		return comparedIndexes, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// applyRollingHash :returns the hash values of each every window
func (rhs *FileDiffFinder) applyRollingHash(s string, windowSize int) []int {
	n := len(s)
	power := make([]int, n+1)
	hashValues := make([]int, 0)

	power[0] = 1
	for i := 1; i <= n; i++ {
		power[i] = (power[i-1] * base) % mod
	}

	// find first hash of first window
	currentHash := rhs.hash(s[:windowSize], windowSize)
	hashValues = append(hashValues, currentHash)

	for i := 1; i <= n-windowSize; i++ {
		currentHash = (currentHash - power[windowSize-1]*int(s[i-1])) % mod
		currentHash = (currentHash*base + int(s[i+windowSize-1])) % mod
		hashValues = append(hashValues, currentHash)
	}
	return hashValues
}

func (rhs *FileDiffFinder) hash(str string, windowSize int) int {
	hash := 0
	for i := 0; i < windowSize; i++ {
		hash = (hash*base + int(str[i])) % mod
	}

	return hash
}

func (rhs *FileDiffFinder) findChangedRanges(originalText, updatedText string, windowSize int) []Range {
	hashValuesOfCurrentFile := rhs.applyRollingHash(originalText, windowSize)
	hashValuesOfUpdatedFile := rhs.applyRollingHash(updatedText, windowSize)

	var rangesToCheck []Range
	var changedIndeces int
	for i, val := range hashValuesOfCurrentFile {
		if val != hashValuesOfUpdatedFile[i] {
			if len(rangesToCheck) == 0 || (len(rangesToCheck) > 0 && rangesToCheck[changedIndeces-1].end < i) {
				rangesToCheck = append(rangesToCheck, Range{start: i, end: i + windowSize})
			} else {
				rangesToCheck[changedIndeces-1].end = i + windowSize
				continue
			}
			changedIndeces++
		}
	}

	return rangesToCheck
}

func (rhs *FileDiffFinder) findMinimalRanges(originalText, updatedText string) []Range {
	// find wider ranges by large window size
	var widerRanges = rhs.findChangedRanges(originalText, updatedText, widerWindowRange)

	// lower ranges of wider ranges
	var ranges []Range
	for _, val := range widerRanges {
		var lowerRanges = rhs.findChangedRanges(rhs.originalFileInfo.Content()[val.start:val.end], updatedText[val.start:val.end], narrowWindowRange)

		for _, val2 := range lowerRanges {
			ranges = append(ranges, Range{start: val.start + val2.start, end: val.start + val2.end})
		}
	}

	return ranges
}

func (rhs *FileDiffFinder) compareIndexes(ranges []Range, updatedFileText string) []UpdatedIndex {
	updatedIndeces := make([]UpdatedIndex, 0)

	// loop ranges and detect updates
	for _, val := range ranges {
		for ind := val.start; ind < val.end; ind++ {
			if rhs.originalFileInfo.Content()[ind] != updatedFileText[ind] {
				updatedIndeces = append(updatedIndeces, UpdatedIndex{
					OldValue: fmt.Sprintf("%c", rhs.originalFileInfo.Content()[ind]),
					NewValue: fmt.Sprintf("%c", updatedFileText[ind]),
					Index:    ind,
					Type:     ChangeTypesUpdated,
				})
			}
		}
	}

	return updatedIndeces
}

func (rhs *FileDiffFinder) compareBruteForce(updatedFileText string) []UpdatedIndex {
	updatedIndeces := make([]UpdatedIndex, 0)
	lower := len(updatedFileText)
	if len(updatedFileText) > len(rhs.originalFileInfo.Content()) {
		lower = len(rhs.originalFileInfo.Content())
	}

	for ind := 0; ind < lower; ind++ {
		if rhs.originalFileInfo.Content()[ind] != updatedFileText[ind] {
			updatedIndeces = append(updatedIndeces, UpdatedIndex{
				OldValue: fmt.Sprintf("%c", rhs.originalFileInfo.Content()[ind]),
				NewValue: fmt.Sprintf("%c", updatedFileText[ind]),
				Index:    ind,
				Type:     ChangeTypesUpdated,
			})
		}
	}

	for ind := lower; ind < len(updatedFileText); ind++ {
		updatedIndeces = append(updatedIndeces, UpdatedIndex{
			NewValue: fmt.Sprintf("%c", updatedFileText[ind]),
			Index:    ind,
			Type:     ChangeTypesAdded,
		})
	}

	for ind := lower; ind < len(rhs.originalFileInfo.Content()); ind++ {
		updatedIndeces = append(updatedIndeces, UpdatedIndex{
			OldValue: fmt.Sprintf("%c", rhs.originalFileInfo.Content()[ind]),
			Index:    ind,
			Type:     ChangeTypesRemoved,
		})
	}

	return updatedIndeces
}

func (rhs *FileDiffFinder) Version() int {
	return rhs.originalFileInfo.Version()
}

func (rhs *FileDiffFinder) Content() string {
	return rhs.originalFileInfo.Content()
}

func (rhs *FileDiffFinder) ValidateVersion(version int) error {
	return rhs.originalFileInfo.ValidateVersion(version)
}
