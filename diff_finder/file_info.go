package diff_finder

import (
	"fmt"
)

type IFileInformative interface {
	Version() int
	Content() string
	ValidateVersion(version int) error
}

type FileInfo struct {
	fileText string
	version  int
}

func NewFileInfo(fileText string, version int) IFileInformative {
	return &FileInfo{
		fileText: fileText,
		version:  version,
	}
}

func (rhs *FileInfo) Version() int {
	return rhs.version
}

func (rhs *FileInfo) Content() string {
	return rhs.fileText
}

func (rhs *FileInfo) ValidateVersion(version int) error {
	if version <= rhs.version {
		return fmt.Errorf("newer version should be used. Current version is %v", rhs.version)
	}

	if version > rhs.version+1 {
		return fmt.Errorf("the latest version is %v. %s", rhs.version, "Please use incremental number for versioning of file info.")
	}

	return nil
}
