package sync_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lonelyday/filesync/internal/config"
	"github.com/lonelyday/filesync/internal/sync"
)

// newLogger returns a silent logger suitable for tests.
func newLogger() *config.Logger {
	return config.NewLogger(false)
}

// writeFile is a test helper that creates a file with given content.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("writeFile mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
}

// readFile is a test helper that reads file content.
func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readFile: %v", err)
	}
	return string(b)
}

// TestSync_CopiesNewFiles verifies that files present in source but not in target are copied.
func TestSync_CopiesNewFiles(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	writeFile(t, filepath.Join(src, "hello.txt"), "hello world")
	writeFile(t, filepath.Join(src, "sub", "nested.txt"), "nested content")

	s := sync.New(sync.Config{
		Source:        src,
		Target:        dst,
		DeleteMissing: false,
		Logger:        newLogger(),
	})

	stats, err := s.Sync()
	if err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	if stats.Copied != 2 {
		t.Errorf("Copied = %d, want 2", stats.Copied)
	}

	if readFile(t, filepath.Join(dst, "hello.txt")) != "hello world" {
		t.Error("hello.txt content mismatch")
	}
	if readFile(t, filepath.Join(dst, "sub", "nested.txt")) != "nested content" {
		t.Error("nested.txt content mismatch")
	}
}

// TestSync_UpdatesChangedFiles verifies that existing but changed files are overwritten.
func TestSync_UpdatesChangedFiles(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	// Create old version in target
	dstFile := filepath.Join(dst, "data.txt")
	writeFile(t, dstFile, "old content")
	// Set an old modification time on the dst file
	old := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(dstFile, old, old); err != nil {
		t.Fatalf("Chtimes: %v", err)
	}

	// Create newer version in source
	writeFile(t, filepath.Join(src, "data.txt"), "new content")

	s := sync.New(sync.Config{
		Source:        src,
		Target:        dst,
		DeleteMissing: false,
		Logger:        newLogger(),
	})

	stats, err := s.Sync()
	if err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	if stats.Updated != 1 {
		t.Errorf("Updated = %d, want 1", stats.Updated)
	}

	if readFile(t, dstFile) != "new content" {
		t.Error("data.txt was not updated to new content")
	}
}

// TestSync_SkipsUpToDateFiles verifies that identical files are not re-copied.
func TestSync_SkipsUpToDateFiles(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	content := "same content"
	srcFile := filepath.Join(src, "same.txt")
	dstFile := filepath.Join(dst, "same.txt")

	writeFile(t, srcFile, content)
	writeFile(t, dstFile, content)

	// Sync mod times so they appear identical
	now := time.Now()
	os.Chtimes(srcFile, now, now)
	os.Chtimes(dstFile, now, now)

	s := sync.New(sync.Config{
		Source:        src,
		Target:        dst,
		DeleteMissing: false,
		Logger:        newLogger(),
	})

	stats, err := s.Sync()
	if err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	if stats.Skipped != 1 {
		t.Errorf("Skipped = %d, want 1", stats.Skipped)
	}
	if stats.Copied != 0 || stats.Updated != 0 {
		t.Errorf("Expected no copies or updates, got Copied=%d Updated=%d", stats.Copied, stats.Updated)
	}
}

// TestSync_PreservesOrphanFiles verifies that extra files in target are kept when delete-missing is off.
func TestSync_PreservesOrphanFiles(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	writeFile(t, filepath.Join(src, "a.txt"), "a")
	writeFile(t, filepath.Join(dst, "orphan.txt"), "should stay")

	s := sync.New(sync.Config{
		Source:        src,
		Target:        dst,
		DeleteMissing: false,
		Logger:        newLogger(),
	})

	_, err := s.Sync()
	if err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(dst, "orphan.txt")); os.IsNotExist(err) {
		t.Error("orphan.txt was deleted but delete-missing is off")
	}
}

// TestSync_DeletesOrphanFiles verifies that extra files in target are removed when delete-missing is on.
func TestSync_DeletesOrphanFiles(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	writeFile(t, filepath.Join(src, "a.txt"), "a")
	orphan := filepath.Join(dst, "orphan.txt")
	writeFile(t, orphan, "delete me")

	s := sync.New(sync.Config{
		Source:        src,
		Target:        dst,
		DeleteMissing: true,
		Logger:        newLogger(),
	})

	stats, err := s.Sync()
	if err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	if stats.Deleted != 1 {
		t.Errorf("Deleted = %d, want 1", stats.Deleted)
	}

	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Error("orphan.txt still exists after delete-missing sync")
	}
}

// TestSync_InvalidSource verifies that Sync returns an error for missing source directory.
func TestSync_InvalidSource(t *testing.T) {
	dst := t.TempDir()

	s := sync.New(sync.Config{
		Source:        "/nonexistent/source/dir",
		Target:        dst,
		DeleteMissing: false,
		Logger:        newLogger(),
	})

	_, err := s.Sync()
	if err == nil {
		t.Error("Sync() expected error for missing source, got nil")
	}
}

// TestSync_CreatesTargetIfMissing verifies that the target directory is created if it doesn't exist.
func TestSync_CreatesTargetIfMissing(t *testing.T) {
	src := t.TempDir()
	dst := filepath.Join(t.TempDir(), "new_target_dir")

	writeFile(t, filepath.Join(src, "file.txt"), "content")

	s := sync.New(sync.Config{
		Source:        src,
		Target:        dst,
		DeleteMissing: false,
		Logger:        newLogger(),
	})

	_, err := s.Sync()
	if err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(dst, "file.txt")); os.IsNotExist(err) {
		t.Error("target dir was not created, file missing")
	}
}

// TestSync_EmptySource verifies syncing from an empty source produces zero stats.
func TestSync_EmptySource(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	s := sync.New(sync.Config{
		Source:        src,
		Target:        dst,
		DeleteMissing: false,
		Logger:        newLogger(),
	})

	stats, err := s.Sync()
	if err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	if stats.Copied != 0 || stats.Updated != 0 || stats.Deleted != 0 || stats.Errors != 0 {
		t.Errorf("Expected all zeros for empty source, got %+v", stats)
	}
}
