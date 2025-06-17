# GoMergeLog4Net

> Cross-platform CLI to flatten, merge, sort and re-expand multiple Log4Net (or similarly formatted) log files.

![build](https://img.shields.io/badge/Go-1.22%2B-blue)

---

## What it does

1. Recursively scans a parent folder for files that match `*.log`, `*.log.1`, `*.log.2`, …
2. Flattens each source log so every multi-line entry becomes **one** physical line (worker-pool, concurrent).
3. Merges all flattened logs into a single file.
4. Orders every record by timestamp (two timestamp patterns supported).
5. Restores the original multi-line structure.
6. Optionally wipes intermediate files, leaving only `<timestamp>_FinalMerged.log`.

The tool is written in pure Go, requires **no** external dependencies and produces fully static binaries (CGO disabled) for Windows, Linux and macOS (intel + arm).

---

## Building

### One-off local build

```bash
# from repo root
cd GoMergeLog4Net

go build -o GoMergeLog4Net .
```

### Cross-platform release builds

* **Windows PowerShell**  – `build_release.ps1`
* **Linux / macOS**       – `build_release.sh`

Both scripts output ready-to-ship binaries in
`Release/<goos>_<goarch>/GoMergeLog4Net[.exe]`.

```bash
# Linux / macOS
chmod +x build_release.sh
./build_release.sh

# Windows
powershell -ExecutionPolicy Bypass -File build_release.ps1
```

---

## Usage

```text
GoMergeLog4Net - merge & order log files (Go version)

Usage:
  -p, --parentFolder <path>   Folder containing logs to process
  -workers <N>               Number of concurrent workers (default 2×CPU)
  -keep                      Keep flattened files (skip purge)
  -v, --verbose              Verbose progress output
  -h                         Show help
```

### Quick example

```bash
./GoMergeLog4Net -p /var/log/myapp -workers 8
```

Result:

```
/var/log/myapp/ProcessedLogs/20240617_103015_FinalMerged.log
```

### Flags explained

* **--parentFolder** *(required)* – Root directory to scan.
* **--workers / ‑workers** – Amount of concurrent file processors. Default = `2 × logical-CPU`.
* **--keep** – Skip cleanup step and keep all intermediate files.
* **--verbose / -v** – Per-file progress messages.

---

## Project structure (after refactor)

```
GoMergeLog4Net/
  main.go                 # CLI + pipeline orchestrator
  build_release.sh        # bash release script
  build_release.ps1       # PowerShell release script
  TestFiles/              # sample logs used by golden tests
  README.md               # this file
```

---

## Testing

```
# unit + golden tests (to be added)
go test ./...
```

### Planned tests

* `detectPattern`, `uniqueFileName`, `newScanner` – table-driven
* Golden tests under `TestFiles/` to assert end-to-end output is stable.

---

## Contributing

1. Open an issue / discussion first if you plan major changes.
2. Ensure `go vet ./...`, `staticcheck ./...`, `golangci-lint run` and
   `go test ./...` all pass before submitting a PR.

---

## License

MIT © Control-Up Support Tools 