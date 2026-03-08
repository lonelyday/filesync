package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies a file from src to dst, preserving modification time.
// It creates any missing parent directories in the destination path.
func CopyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("failed to create directories for %s: %w", dst, err)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file %s: %w", src, err)
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy data from %s to %s: %w", src, dst, err)
	}

	// Preserve modification time
	if err := os.Chtimes(dst, srcInfo.ModTime(), srcInfo.ModTime()); err != nil {
		// Non-fatal: log as warning but don't fail the copy
		return fmt.Errorf("failed to set modification time on %s: %w", dst, err)
	}

	return nil
}

// WalkDir recursively walks a directory and returns a map of relative path → FileInfo.
// Returns an error only if the root directory itself cannot be read.
func WalkDir(root string) (map[string]FileInfo, error) {
	files := make(map[string]FileInfo)

	err := walkDirRecursive(root, root, files)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func walkDirRecursive(root, current string, files map[string]FileInfo) error {
	entries, err := os.ReadDir(current)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", current, err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(current, entry.Name())
		relPath, err := filepath.Rel(root, fullPath)
		if err != nil {
			return fmt.Errorf("failed to compute relative path for %s: %w", fullPath, err)
		}

		info, err := entry.Info()
		if err != nil {
			// Skip files we can't stat, caller can log this
			continue
		}

		fi := FileInfo{
			Path:    fullPath,
			RelPath: relPath,
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		}
		files[relPath] = fi

		if entry.IsDir() {
			if err := walkDirRecursive(root, fullPath, files); err != nil {
				// Continue walking, don't abort on single dir error
				continue
			}
		}
	}

	return nil
}

// EnsureDir creates a directory and all necessary parents.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// RemoveFile removes a file and returns an error if it fails.
func RemoveFile(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove %s: %w", path, err)
	}
	return nil
}

// DirExists returns true if the path exists and is a directory.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
