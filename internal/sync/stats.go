package sync

// Stats holds the results of a sync operation.
type Stats struct {
	Copied  int
	Updated int
	Deleted int
	Skipped int
	Errors  int
}
