package sync

import (
	"github.com/lonelyday/filesync/internal/file"
)

// NeedsUpdate returns true if the source file is different from the target file.
// Comparison is based on size and modification time.
func NeedsUpdate(src, dst file.FileInfo) bool {
	if src.Size != dst.Size {
		return true
	}
	// Source is newer than target
	return src.ModTime.After(dst.ModTime)
}
