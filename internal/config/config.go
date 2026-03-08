package config

// Config holds the complete configuration for the filesync application.
type Config struct {
	Source        string
	Target        string
	DeleteMissing bool
	Verbose       bool
	ShowVersion   bool
}
