package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/lonelyday/filesync/internal/file"
)

// Config holds configuration for a sync operation.
type Config struct {
	Source        string
	Target        string
	DeleteMissing bool
	Logger        Logger
}

// Logger interface for logging operations.
type Logger interface {
	Info(msg string)
	Infof(format string, args ...interface{})
	Debug(msg string)
	Debugf(format string, args ...interface{})
	Success(msg string)
	Successf(format string, args ...interface{})
	Warn(msg string)
	Error(msg string)
	Errorf(format string, args ...interface{})
}

// Syncer performs one-way directory synchronization.
type Syncer struct {
	cfg Config
}

// New creates a new Syncer with the provided configuration.
func New(cfg Config) *Syncer {
	return &Syncer{cfg: cfg}
}

// Sync performs the one-way synchronization from source to target.
// Returns Stats describing what was done, and an error only for fatal issues.
func (s *Syncer) Sync() (Stats, error) {
	var stats Stats
	log := s.cfg.Logger

	// Validate source directory
	if !file.DirExists(s.cfg.Source) {
		return stats, fmt.Errorf("source directory does not exist: %s", s.cfg.Source)
	}

	// Ensure target directory exists
	if err := file.EnsureDir(s.cfg.Target); err != nil {
		return stats, fmt.Errorf("cannot create target directory %s: %w", s.cfg.Target, err)
	}

	log.Infof("Starting sync: %s → %s", s.cfg.Source, s.cfg.Target)
	if s.cfg.DeleteMissing {
		log.Info("Delete-missing mode enabled")
	}

	// Walk source and target
	log.Debug("Scanning source directory...")
	srcFiles, err := file.WalkDir(s.cfg.Source)
	if err != nil {
		return stats, fmt.Errorf("failed to scan source: %w", err)
	}
	log.Debugf("Found %d entries in source", len(srcFiles))

	log.Debug("Scanning target directory...")
	dstFiles, err := file.WalkDir(s.cfg.Target)
	if err != nil {
		return stats, fmt.Errorf("failed to scan target: %w", err)
	}
	log.Debugf("Found %d entries in target", len(dstFiles))

	// Process source files: copy new, update changed
	srcPaths := sortedKeys(srcFiles)
	for _, relPath := range srcPaths {
		srcFile := srcFiles[relPath]

		// Skip directories — they are created automatically during file copy
		if srcFile.IsDir {
			log.Debugf("Skipping directory entry: %s", relPath)
			continue
		}

		dstFile, existsInDst := dstFiles[relPath]

		if !existsInDst {
			// File is new — copy it
			dstPath := targetPath(s.cfg.Target, relPath)
			log.Infof("Copying: %s", relPath)
			if err := file.CopyFile(srcFile.Path, dstPath); err != nil {
				log.Errorf("Failed to copy %s: %v", relPath, err)
				stats.Errors++
			} else {
				log.Successf("Copied: %s", relPath)
				stats.Copied++
			}
		} else if NeedsUpdate(srcFile, dstFile) {
			// File exists but differs — overwrite it
			dstPath := targetPath(s.cfg.Target, relPath)
			log.Infof("Updating: %s", relPath)
			if err := file.CopyFile(srcFile.Path, dstPath); err != nil {
				log.Errorf("Failed to update %s: %v", relPath, err)
				stats.Errors++
			} else {
				log.Successf("Updated: %s", relPath)
				stats.Updated++
			}
		} else {
			// File is identical — skip
			log.Debugf("Skipped (up to date): %s", relPath)
			stats.Skipped++
		}
	}

	// Delete files in target that no longer exist in source
	if s.cfg.DeleteMissing {
		// Process in reverse order so files are deleted before their parent dirs
		dstPaths := sortedKeysReverse(dstFiles)
		for _, relPath := range dstPaths {
			if _, existsInSrc := srcFiles[relPath]; !existsInSrc {
				dstFile := dstFiles[relPath]
				dstPath := targetPath(s.cfg.Target, relPath)

				if dstFile.IsDir {
					// Only remove empty directories
					log.Debugf("Checking orphan directory: %s", relPath)
					if err := removeEmptyDir(dstPath); err != nil {
						log.Debugf("Directory not removed (probably not empty): %s", relPath)
					} else {
						log.Successf("Removed empty directory: %s", relPath)
						stats.Deleted++
					}
				} else {
					log.Infof("Deleting (not in source): %s", relPath)
					if err := file.RemoveFile(dstPath); err != nil {
						log.Errorf("Failed to delete %s: %v", relPath, err)
						stats.Errors++
					} else {
						log.Successf("Deleted: %s", relPath)
						stats.Deleted++
					}
				}
			}
		}
	}

	return stats, nil
}

// targetPath joins the target root with a relative path.
func targetPath(target, relPath string) string {
	return filepath.Join(target, relPath)
}

// sortedKeys returns a sorted slice of map keys.
func sortedKeys(m map[string]file.FileInfo) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// sortedKeysReverse returns a reverse-sorted slice of map keys.
func sortedKeysReverse(m map[string]file.FileInfo) []string {
	keys := sortedKeys(m)
	for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
		keys[i], keys[j] = keys[j], keys[i]
	}
	return keys
}

// removeEmptyDir removes a directory only if it is empty.
func removeEmptyDir(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("cannot read directory %s: %w", path, err)
	}
	if len(entries) != 0 {
		return fmt.Errorf("directory not empty: %s", path)
	}
	return os.Remove(path)
}
