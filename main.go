package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	processedFolderName = "ProcessedLogs"
	sentinel            = "\x1E"           // must never appear in real log text
	maxScannerCapacity  = 10 * 1024 * 1024 // 10 MiB per record
)

var (
	// Pattern A: '2024-06-16 10:42:15,123'
	patternA = regexp.MustCompile(`^'(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2},\d{3})'`)
	// Pattern B:  2024-06-16 10:42:15.123
	patternB    = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3})`)
	verboseMode bool
	version     = "dev" // overwritten at build time via -ldflags "-X main.version=<ver>"
)

// -------- verbose helper ----------
func vPrintf(format string, args ...any) {
	if verboseMode {
		log.Printf(format, args...)
	}
}

// ---------------------------------- FLAGS ----------------------------------
func main() {
	start := time.Now() // timer starts immediately

	parent := flag.String("parentFolder", "", "Directory containing log files")
	parentShort := flag.String("p", "", "Directory containing log files (shorthand)")
	keep := flag.Bool("keep", false, "Keep flattened files (skip purge)")
	workers := flag.Int("workers", runtime.NumCPU()*2, "Concurrent file-processing workers (default = 2 × CPU cores)")

	// verbose and help flags
	verboseFlag := flag.Bool("verbose", false, "Verbose output")
	flag.BoolVar(verboseFlag, "v", false, "Verbose output (shorthand)")
	help := flag.Bool("h", false, "Show help")
	versionFlag := flag.Bool("version", false, "Show application version")

	flag.Parse()

	verboseMode = *verboseFlag

	if *help {
		displayHelp()
		return
	}

	if *versionFlag {
		fmt.Println(version)
		return
	}

	// Resolve folder argument (long or short)
	parentFolder := firstNonEmpty(*parent, *parentShort)
	if parentFolder == "" {
		fmt.Fprintln(os.Stderr, "Error: --parentFolder / -p is required")
		displayHelp()
		os.Exit(1)
	}
	if stat, err := os.Stat(parentFolder); err != nil || !stat.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: %q is not a valid directory\n", parentFolder)
		os.Exit(1)
	}

	processFolder, err := createProcessedLogsFolder(parentFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create %s: %v\n", processedFolderName, err)
		os.Exit(1)
	}

	logFiles, err := getAllLogFiles(parentFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
		os.Exit(1)
	}
	if len(logFiles) == 0 {
		fmt.Println("No .log files found.")
		return
	}

	// ---------- 1) FLATTEN EACH FILE USING A WORKER POOL ----------
	processedCh := make(chan string)
	taskCh := make(chan string)

	// 1a. spawn bounded workers
	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for input := range taskCh {
				base := filepath.Base(input)
				outPath := uniqueFileName(filepath.Join(processFolder, base))
				if err := processLogFile(input, outPath); err != nil {
					fmt.Fprintf(os.Stderr, "Processing %s failed: %v\n", input, err)
					continue
				}
				vPrintf("Flattened: %s -> %s", input, outPath)
				processedCh <- outPath
			}
		}()
	}

	// 1b. feed all file paths to the workers
	go func() {
		for _, lf := range logFiles {
			taskCh <- lf
		}
		close(taskCh)
	}()

	// 1c. close processedCh when all workers are finished
	go func() {
		wg.Wait()
		close(processedCh)
	}()

	var processed []string
	for p := range processedCh {
		processed = append(processed, p)
	}

	// ---------- 2) MERGE, ORDER, RESTORE ----------
	merged := filepath.Join(processFolder, "MERGED.log")
	if err := mergeProcessedLogs(processed, merged); err != nil {
		fmt.Fprintf(os.Stderr, "Merge failed: %v\n", err)
		os.Exit(1)
	}

	ordered := filepath.Join(processFolder, "MERGED_ORDERED.log")
	if err := orderByDate(merged, ordered); err != nil {
		fmt.Fprintf(os.Stderr, "Ordering failed: %v\n", err)
		os.Exit(1)
	}

	timestamp := time.Now().Format("20060102_150405")
	finalPath := filepath.Join(processFolder, fmt.Sprintf("%s_FinalMerged.log", timestamp))
	if err := formatSupport(ordered, finalPath); err != nil {
		fmt.Fprintf(os.Stderr, "Re-format failed: %v\n", err)
		os.Exit(1)
	}

	if !*keep {
		if err := cleanupProcessFolder(processFolder, finalPath); err != nil {
			fmt.Fprintf(os.Stderr, "Cleanup warning: %v\n", err)
		}
	} else {
		vPrintf("Cleanup skipped because --keep flag was provided.")
	}
	log.Printf("%d workers finished in %s. Final file: %s",
		*workers,
		time.Since(start).Round(time.Millisecond),
		finalPath)
}

func displayHelp() {
	fmt.Println("GoMergeLog4Net - merge & order log files (Go version)")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  -p, --parentFolder <path>   Folder containing logs to process")
	fmt.Println("  -workers <N>               Number of concurrent workers (default 2×CPU)")
	fmt.Println("  -keep                      Keep flattened files (skip purge)")
	fmt.Println("  -h                         Show help")
}

func createProcessedLogsFolder(parent string) (string, error) {
	path := filepath.Join(parent, processedFolderName)
	err := os.MkdirAll(path, 0755)
	return path, err
}

func getAllLogFiles(parent string) ([]string, error) {
	var out []string
	err := filepath.Walk(parent, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// skip the ProcessedLogs folder
		if info.IsDir() && info.Name() == processedFolderName {
			return filepath.SkipDir
		}
		if info.Mode().IsRegular() && matchesLogPattern(info.Name()) {
			out = append(out, path)
		}
		return nil
	})
	return out, err
}

var logPattern = regexp.MustCompile(`(?i)\.log(\.\d+)?$`)

func matchesLogPattern(name string) bool {
	return logPattern.MatchString(name)
}

func uniqueFileName(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}
	ext := filepath.Ext(path) // ".log"
	name := strings.TrimSuffix(path, ext)
	i := 1
	for {
		candidate := fmt.Sprintf("%s%d%s", name, i, ext)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
		i++
	}
}

// ---------- per-file processing ----------

// newScanner builds a bufio.Scanner whose buffer can grow up to
// maxScannerCapacity, avoiding the default 64 KiB limit.
func newScanner(r io.Reader) *bufio.Scanner {
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 64*1024), maxScannerCapacity)
	return s
}

func processLogFile(input, output string) error {
	file, err := os.Open(input)
	if err != nil {
		return err
	}
	defer file.Close()

	pattern, err := detectPattern(file)
	if err != nil {
		return err
	}
	// reset reader to start
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	scanner := newScanner(file)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	var entry strings.Builder // build the current flattened record

	for scanner.Scan() {
		line := scanner.Text()
		if pattern.MatchString(line) {
			if entry.Len() > 0 { // flush previous entry
				fmt.Fprintln(writer, entry.String())
				entry.Reset()
			}
			entry.WriteString(line) // start new entry
		} else {
			if entry.Len() > 0 { // continuation line
				entry.WriteString(sentinel)
				entry.WriteString(line)
			}
		}
	}
	if entry.Len() > 0 { // flush last entry
		fmt.Fprintln(writer, entry.String())
	}
	return scanner.Err()
}

func detectPattern(r io.Reader) (*regexp.Regexp, error) {
	reader := bufio.NewReader(r)
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}
	switch {
	case patternA.MatchString(line):
		return patternA, nil
	case patternB.MatchString(line):
		return patternB, nil
	default:
		return nil, fmt.Errorf("unrecognised log format")
	}
}

// ---------- merge ----------

func mergeProcessedLogs(files []string, dest string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	defer w.Flush()

	// Re-use one 128 KiB buffer for every copy operation to
	// avoid repeated 32 KiB allocations inside io.Copy.
	buf := make([]byte, 128*1024)

	for _, f := range files {
		in, err := os.Open(f)
		if err != nil {
			return err
		}
		if _, err := io.CopyBuffer(w, in, buf); err != nil { // use shared buffer
			in.Close()
			return err
		}
		in.Close()
	}
	return nil
}

// ---------- order chronologically ----------

type entry struct {
	ts   time.Time
	line string
}

func orderByDate(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	var entries []entry
	scanner := newScanner(in)
	for scanner.Scan() {
		l := scanner.Text()
		pat := patternA
		layout := "2006-01-02 15:04:05,000"
		m := pat.FindStringSubmatch(l)
		if m == nil {
			pat = patternB
			layout = "2006-01-02 15:04:05.000"
			m = pat.FindStringSubmatch(l)
		}
		if m == nil {
			continue // skip if no timestamp found
		}
		t, err := time.ParseInLocation(layout, m[1], time.UTC)
		if err != nil {
			continue
		}
		entries = append(entries, entry{t, l})
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].ts.Before(entries[j].ts) })

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	for _, e := range entries {
		fmt.Fprintln(w, e.line)
	}
	return w.Flush()
}

// ---------- restore multi-line structure ----------

func formatSupport(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	scanner := newScanner(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), sentinel)
		for _, p := range parts {
			fmt.Fprintln(writer, p)
		}
	}
	return scanner.Err()
}

// ---------- cleanup ----------

func cleanupProcessFolder(folder, final string) error {
	files, err := os.ReadDir(folder)
	if err != nil {
		return err
	}
	for _, f := range files {
		p := filepath.Join(folder, f.Name())
		if p == final || strings.Contains(p, "_FinalMerged") {
			continue
		}
		_ = os.Remove(p) // ignore errors, just best effort
	}
	return nil
}

// firstNonEmpty returns the first non-empty string from its arguments.
func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
