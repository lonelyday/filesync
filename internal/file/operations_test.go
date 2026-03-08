package file_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lonelyday/filesync/internal/file"
)

// TestCopyFile verifies that file copying works correctly.
func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()
	srcPath := filepath.Join(tmpDir, "source.txt")
	dstPath := filepath.Join(tmpDir, "subdir", "dest.txt")

	content := []byte("hello filesync world")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	if err := file.CopyFile(srcPath, dstPath); err != nil {
		t.Fatalf("CopyFile() error = %v", err)
	}

	got, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}

	if string(got) != string(content) {
		t.Errorf("CopyFile() content = %q, want %q", string(got), string(content))
	}
}

// TestCopyFile_CreatesParentDirs verifies that CopyFile creates missing parent directories.
func TestCopyFile_CreatesParentDirs(t *testing.T) {
	tmpDir := t.TempDir()
	srcPath := filepath.Join(tmpDir, "src.txt")
	dstPath := filepath.Join(tmpDir, "a", "b", "c", "dst.txt")

	if err := os.WriteFile(srcPath, []byte("data"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := file.CopyFile(srcPath, dstPath); err != nil {
		t.Errorf("CopyFile() with deep path error = %v", err)
	}

	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		t.Error("CopyFile() did not create destination file")
	}
}

// TestCopyFile_SourceMissing verifies that CopyFile returns an error for missing source.
func TestCopyFile_SourceMissing(t *testing.T) {
	tmpDir := t.TempDir()
	err := file.CopyFile(filepath.Join(tmpDir, "nonexistent.txt"), filepath.Join(tmpDir, "dst.txt"))
	if err == nil {
		t.Error("CopyFile() expected error for missing source, got nil")
	}
}

// TestWalkDir verifies that WalkDir correctly enumerates files.
func TestWalkDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create structure:
	// tmpDir/
	//   file1.txt
	//   subdir/
	//     file2.txt
	files := map[string]string{
		"file1.txt":        "content1",
		"subdir/file2.txt": "content2",
	}
	for relPath, content := range files {
		fullPath := filepath.Join(tmpDir, filepath.FromSlash(relPath))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("setup: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	result, err := file.WalkDir(tmpDir)
	if err != nil {
		t.Fatalf("WalkDir() error = %v", err)
	}

	// Should contain file1.txt, subdir, subdir/file2.txt
	expectedEntries := []string{
		"file1.txt",
		"subdir",
		filepath.Join("subdir", "file2.txt"),
	}

	for _, entry := range expectedEntries {
		if _, ok := result[entry]; !ok {
			t.Errorf("WalkDir() missing expected entry: %q", entry)
		}
	}
}

// TestWalkDir_EmptyDir verifies that WalkDir handles empty directories.
func TestWalkDir_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	result, err := file.WalkDir(tmpDir)
	if err != nil {
		t.Fatalf("WalkDir() error = %v", err)
	}
	if len(result) != 0 {
		t.Errorf("WalkDir() on empty dir returned %d entries, want 0", len(result))
	}
}

// TestWalkDir_NonExistentDir verifies error handling for missing directories.
func TestWalkDir_NonExistentDir(t *testing.T) {
	_, err := file.WalkDir("/path/that/does/not/exist")
	if err == nil {
		t.Error("WalkDir() expected error for non-existent directory, got nil")
	}
}

// TestDirExists verifies the DirExists helper.
func TestDirExists(t *testing.T) {
	tmpDir := t.TempDir()

	if !file.DirExists(tmpDir) {
		t.Error("DirExists() = false for existing dir, want true")
	}

	if file.DirExists(filepath.Join(tmpDir, "nope")) {
		t.Error("DirExists() = true for missing dir, want false")
	}
}

// TestEnsureDir verifies that EnsureDir creates nested directories.
func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "a", "b", "c")

	if err := file.EnsureDir(newDir); err != nil {
		t.Fatalf("EnsureDir() error = %v", err)
	}

	if !file.DirExists(newDir) {
		t.Error("EnsureDir() did not create the directory")
	}
}

// TestRemoveFile verifies file removal.
func TestRemoveFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "to_delete.txt")

	if err := os.WriteFile(filePath, []byte("bye"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := file.RemoveFile(filePath); err != nil {
		t.Fatalf("RemoveFile() error = %v", err)
	}

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("RemoveFile() file still exists after removal")
	}
}

// TestRemoveFile_Missing verifies error returned for missing file.
func TestRemoveFile_Missing(t *testing.T) {
	err := file.RemoveFile("/tmp/filesync_missing_file_xyz.txt")
	if err == nil {
		t.Error("RemoveFile() expected error for missing file, got nil")
	}
}
