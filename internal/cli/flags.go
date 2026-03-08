package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/lonelyday/filesync/internal/config"
)

const version = "1.0.0"

// ParseFlags parses command line flags and returns the configuration.
// Returns an error if required flags are missing or invalid.
func ParseFlags() (*config.Config, error) {
	return ParseArgs(os.Args[1:])
}

// ParseArgs parses the given arguments as flags and returns the configuration.
// This is a testable version of ParseFlags.
func ParseArgs(args []string) (*config.Config, error) {
	var cfg config.Config

	fs := flag.NewFlagSet("filesync", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	fs.StringVar(&cfg.Source, "source", "", "Path to source directory (required)")
	fs.StringVar(&cfg.Source, "s", "", "Path to source directory (shorthand)")
	fs.StringVar(&cfg.Target, "target", "", "Path to target directory (required)")
	fs.StringVar(&cfg.Target, "t", "", "Path to target directory (shorthand)")
	fs.BoolVar(&cfg.DeleteMissing, "delete-missing", false, "Delete files in target that don't exist in source")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Enable verbose logging")
	fs.BoolVar(&cfg.ShowVersion, "version", false, "Show version information")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "filesync v%s - One-way directory synchronization tool\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  filesync --source <path> --target <path> [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  filesync --source ./docs --target ./backup\n")
		fmt.Fprintf(os.Stderr, "  filesync --source /data/src --target /data/dst --delete-missing\n")
		fmt.Fprintf(os.Stderr, "  filesync -s ./src -t ./dst --verbose\n")
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if cfg.ShowVersion {
		fmt.Printf("filesync v%s\n", version)
		os.Exit(0)
	}

	if cfg.Source == "" || cfg.Target == "" {
		fs.Usage()
		return nil, fmt.Errorf("source and target directories are required")
	}

	return &cfg, nil
}
