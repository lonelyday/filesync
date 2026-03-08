package file

import (
	"time"
)

// FileInfo holds metadata about a file used for comparison.
type FileInfo struct {
	Path    string
	RelPath string
	Size    int64
	ModTime time.Time
	IsDir   bool
}
