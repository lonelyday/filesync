package app

import (
	"fmt"
	"os"

	"github.com/lonelyday/filesync/internal/cli"
	"github.com/lonelyday/filesync/internal/config"
	"github.com/lonelyday/filesync/internal/sync"
)

// Run executes the filesync application.
func Run() error {
	// Parse command line arguments
	cfg, err := cli.ParseFlags()
	if err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	// Create logger
	logger := config.NewLogger(cfg.Verbose)

	// Create sync configuration
	syncCfg := sync.Config{
		Source:        cfg.Source,
		Target:        cfg.Target,
		DeleteMissing: cfg.DeleteMissing,
		Logger:        logger,
	}

	// Create and run syncer
	syncer := sync.New(syncCfg)
	stats, err := syncer.Sync()
	if err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	// Print summary
	logger.Successf("Sync completed: %d copied, %d updated, %d deleted, %d skipped, %d errors",
		stats.Copied, stats.Updated, stats.Deleted, stats.Skipped, stats.Errors)

	// Exit with error code if there were errors
	if stats.Errors > 0 {
		os.Exit(1)
	}

	return nil
}
