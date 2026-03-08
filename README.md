# filesync

A CLI tool for one-way directory synchronization.

Copies files from a **source** directory to a **target** directory — new files are copied, changed files are updated, and (optionally) files missing from source are deleted from target.

---

## Project Structure

```
filesync/
├── cmd/
│   └── filesync/
│       └── main.go             # CLI entry point
├── internal/
│   ├── app/
│   │   └── app.go              # Main application logic
│   ├── cli/
│   │   ├── flags.go            # Command-line argument parsing
│   │   └── flags_test.go       # CLI tests
│   ├── config/
│   │   ├── config.go           # Application configuration
│   │   └── logger.go           # Structured logger
│   ├── sync/
│   │   ├── syncer.go           # Core sync logic
│   │   ├── syncer_test.go      # Sync tests
│   │   ├── stats.go            # Sync statistics
│   │   └── compare.go          # File comparison logic
│   └── file/
│       ├── operations.go       # File operations (copy, walk, etc.)
│       ├── operations_test.go  # File operation tests
│       └── info.go             # File information structures
├── go.mod
└── README.md
```

---

## Requirements

- Go 1.18 or newer
- Linux, macOS, or Windows

---

## Building

```bash
# Clone github.com/lonelyday/filesync
cd filesync

# Build binary
go build -o filesync ./cmd/filesync

# On Windows
go build -o filesync.exe ./cmd/filesync
```

---

## Running

```
filesync --source <path> --target <path> [options]
```

### Flags

| Flag               | Short | Default | Description                                               |
|--------------------|-------|---------|-----------------------------------------------------------|
| `--source`         | `-s`  | —       | **Required.** Path to the source directory                |
| `--target`         | `-t`  | —       | **Required.** Path to the target (destination) directory  |
| `--delete-missing` | —     | false   | Delete files in target that are absent in source          |
| `--verbose`        | —     | false   | Print debug-level output (skipped files, dir scans, etc.) |
| `--version`        | —     | —       | Print version and exit                                    |

---

## Examples

### Basic sync (copy new and changed files)
```bash
./filesync --source ./documents --target ./backup
```

### Sync and remove files from target that no longer exist in source
```bash
./filesync --source ./documents --target ./backup --delete-missing
```

### Verbose output with shorthand flags
```bash
./filesync -s /data/source -t /data/destination --verbose
```

### Windows example
```cmd
filesync.exe --source C:\Users\Alice\Documents --target D:\Backup\Documents --delete-missing
```

---

## Sample Output

```
[2024-03-15 12:00:01] INFO  Starting sync: ./src → ./dst
[2024-03-15 12:00:01] INFO  Copying: report.pdf
[2024-03-15 12:00:01] OK    Copied: report.pdf
[2024-03-15 12:00:01] INFO  Updating: notes/todo.txt
[2024-03-15 12:00:01] OK    Updated: notes/todo.txt
[2024-03-15 12:00:01] INFO  Deleting (not in source): old_draft.docx
[2024-03-15 12:00:01] OK    Deleted: old_draft.docx

Sync complete: 1 copied, 1 updated, 1 deleted, 3 skipped, 0 errors
```

---

## How It Works

1. **Scan** — Both source and target directories are recursively scanned.
2. **Compare** — Each file in source is compared against the target by **size** and **modification time**.
3. **Copy** — Files only in source are copied (parent directories are created automatically).
4. **Update** — Files that differ (newer or different size in source) are overwritten in target.
5. **Delete** *(optional)* — With `--delete-missing`, files in target that have no counterpart in source are removed.
6. **Report** — A summary is printed at the end with counts of each operation.

Errors (e.g., permission denied) are logged to stderr and do **not** abort the sync — the tool continues processing remaining files.

---

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run tests for a specific package
go test -v ./internal/file/
go test -v ./internal/sync/
go test -v ./internal/cli/
```

---

## Exit Codes

| Code | Meaning                                  |
|------|------------------------------------------|
| `0`  | Sync completed (errors are non-fatal)    |
| `1`  | Fatal error (missing source, bad config) |
