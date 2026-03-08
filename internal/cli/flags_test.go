package cli_test

import (
	"testing"

	"github.com/lonelyday/filesync/internal/cli"
	"github.com/lonelyday/filesync/internal/config"
)

func TestParseArgs_ValidFlags(t *testing.T) {
	args := []string{
		"--source", "/tmp/src",
		"--target", "/tmp/dst",
		"--delete-missing",
		"--verbose",
	}

	cfg, err := cli.ParseArgs(args)
	if err != nil {
		t.Fatalf("ParseArgs() error = %v", err)
	}

	expected := &config.Config{
		Source:        "/tmp/src",
		Target:        "/tmp/dst",
		DeleteMissing: true,
		Verbose:       true,
		ShowVersion:   false,
	}

	if *cfg != *expected {
		t.Errorf("ParseArgs() = %+v, want %+v", cfg, expected)
	}
}

func TestParseArgs_MissingRequired(t *testing.T) {
	args := []string{
		"--source", "/tmp/src",
	}

	_, err := cli.ParseArgs(args)
	if err == nil {
		t.Error("ParseArgs() expected error for missing target, got nil")
	}
}

func TestParseArgs_ShorthandFlags(t *testing.T) {
	args := []string{
		"-s", "/tmp/src",
		"-t", "/tmp/dst",
	}

	cfg, err := cli.ParseArgs(args)
	if err != nil {
		t.Fatalf("ParseArgs() error = %v", err)
	}

	if cfg.Source != "/tmp/src" || cfg.Target != "/tmp/dst" {
		t.Errorf("ParseArgs() shorthand flags not parsed correctly: source=%s, target=%s", cfg.Source, cfg.Target)
	}
}
