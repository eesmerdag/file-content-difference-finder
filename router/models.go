package router

import (
	"github.com/pkg/errors"

	df "file-diff-finder/diff_finder"
)

type RequestPayload struct {
	Text    *string `json:"text"`
	Version int     `json:"version"`
}

type Response struct {
	Delta          []df.UpdatedIndex `json:"delta"`
	CurrentVersion int               `json:"current_version"`
	UpdatedVersion int               `json:"updated_version"`
}

func (rp *RequestPayload) Valid() error {
	if rp.Version <= 0 {
		return errors.New("Version must be provided and should be positive integer")
	}

	if rp.Text == nil {
		return errors.New("missing text in payload")
	}

	return nil
}
