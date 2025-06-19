# Application Requirements Document (ARD)  
**Version:** 1.0  
**Last Updated:** December 2024

---

## 1. Introduction

### 1.1 Purpose  
This document defines the functional and non-functional requirements for GoMergeLog4Net, a cross-platform command-line application designed to efficiently process, merge, and chronologically order multiple Log4Net (or similarly formatted) log files. The application must handle large log files through concurrent processing, support multiple timestamp formats, and provide a streamlined workflow for log analysis and debugging.

### 1.2 Scope  
- **Target Platforms:** Linux, macOS, and Windows (Intel and ARM architectures).
- **Core Functionality:**  
  - Recursive scanning and discovery of log files.
  - Multi-line log entry flattening with concurrent processing.
  - Merging of multiple log files into a single output.
  - Chronological ordering by timestamp extraction and parsing.
  - Restoration of original multi-line log structure.
  - Automated cleanup of intermediate processing files.
- **Performance:**  
  - Concurrent processing using worker pools for optimal throughput.
  - Memory-efficient handling of large log files.
  - Minimal disk I/O through optimized buffering strategies.
- **Extensibility:**  
  - Modular architecture supporting additional timestamp formats.
  - Configurable processing parameters for different use cases.

---

## 2. System Overview

The application is composed of several clearly defined modules:

- **File Discovery Manager:**  
  Recursively scans directories for log files matching specified patterns (`*.log`, `*.log.1`, `*.log.2`, etc.).

- **Pattern Detection Manager:**  
  Automatically detects and validates timestamp patterns in log files (supports multiple formats).

- **Log Processing Manager:**  
  Handles the flattening of multi-line log entries into single-line records for efficient processing.

- **Worker Pool Manager:**  
  Manages concurrent processing of multiple log files using configurable worker pools.

- **Merge Manager:**  
  Combines multiple processed log files into a single consolidated file.

- **Sorting Manager:**  
  Extracts timestamps and chronologically orders all log entries.

- **Format Restoration Manager:**  
  Restores the original multi-line structure of log entries.

- **File Management Manager:**  
  Handles file creation, cleanup, and unique naming to prevent conflicts.

- **Configuration Manager:**  
  Centralizes command-line parameters and processing options.

- **Progress Reporting Manager:**  
  Provides verbose output and progress tracking during processing.

---

## 3. Functional Requirements

### 3.1 File Discovery and Validation  
- **Recursive Directory Scanning:**  
  - Scan specified parent directory and all subdirectories for log files.
  - Support for log file patterns: `*.log`, `*.log.1`, `*.log.2`, etc.
  - Automatic exclusion of the `ProcessedLogs` output directory.
  - Validation of directory existence and accessibility.

### 3.2 Timestamp Pattern Detection  
- **Automatic Format Detection:**  
  - Pattern A: `'2024-06-16 10:42:15,123'` (quoted format with comma)
  - Pattern B: `2024-06-16 10:42:15.123` (unquoted format with period)
  - Automatic detection of the first valid timestamp in each log file.
  - Error handling for unrecognized log formats.

### 3.3 Log Entry Flattening  
- **Multi-line Entry Processing:**  
  - Identify log entries by timestamp patterns at line beginnings.
  - Flatten multi-line log entries into single physical lines.
  - Use sentinel characters (`\x1E`) to separate continuation lines.
  - Handle entries without proper timestamps gracefully.
  - Support for large log entries (up to 10 MiB per record).

### 3.4 Concurrent Processing  
- **Worker Pool Implementation:**  
  - Configurable number of concurrent workers (default: 2 × CPU cores).
  - Concurrent processing of multiple log files.
  - Thread-safe file operations and output generation.
  - Progress tracking and error handling per worker.

### 3.5 File Merging  
- **Consolidation Process:**  
  - Merge all flattened log files into a single output file.
  - Efficient buffered I/O operations (128 KiB buffer per copy).
  - Maintain file integrity during merge operations.
  - Handle large merged files without memory exhaustion.

### 3.6 Chronological Ordering  
- **Timestamp Extraction and Sorting:**  
  - Extract timestamps from flattened log entries.
  - Parse timestamps using detected format patterns.
  - Sort all entries chronologically in ascending order.
  - Handle timezone conversion (UTC parsing).
  - Skip entries with invalid or missing timestamps.

### 3.7 Format Restoration  
- **Multi-line Structure Recovery:**  
  - Restore original multi-line log entry structure.
  - Split flattened entries using sentinel characters.
  - Maintain proper line breaks and formatting.
  - Generate final output with timestamped filename.

### 3.8 File Management  
- **Output Organization:**  
  - Create `ProcessedLogs` subdirectory in parent folder.
  - Generate unique filenames to prevent conflicts.
  - Timestamp-based final output naming (`YYYYMMDD_HHMMSS_FinalMerged.log`).
  - Optional cleanup of intermediate processing files.

### 3.9 Command-Line Interface  
- **Required Parameters:**  
  - `--parentFolder` or `-p`: Directory containing log files to process.
- **Optional Parameters:**  
  - `--workers`: Number of concurrent processing workers.
  - `--keep`: Retain intermediate files (skip cleanup).
  - `--verbose` or `-v`: Enable detailed progress output.
  - `--help` or `-h`: Display usage information.
  - `--version`: Display application version.

---

## 4. Non-Functional Requirements

### 4.1 Performance  
- **Large File Support:** Efficiently process log files of any size through streaming I/O.
- **Concurrent Processing:** Utilize multiple CPU cores for parallel file processing.
- **Memory Efficiency:** Minimize memory footprint through buffered operations.
- **Processing Speed:** Optimize for high-throughput log processing scenarios.

### 4.2 Cross-Platform Compatibility  
- **Supported Platforms:** Linux, macOS, and Windows (Intel and ARM).
- **Static Binaries:** Generate self-contained executables with no external dependencies.
- **CGO Disabled:** Ensure maximum portability and deployment simplicity.
- **Consistent Behavior:** Maintain identical functionality across all platforms.

### 4.3 Reliability and Error Handling  
- **Graceful Degradation:** Continue processing when individual files fail.
- **Error Reporting:** Provide clear error messages for troubleshooting.
- **File Integrity:** Ensure output files are complete and properly formatted.
- **Resource Management:** Proper cleanup of file handles and system resources.

### 4.4 Maintainability and Extensibility  
- **Modular Architecture:** Clear separation of concerns between processing stages.
- **Code Documentation:** Comprehensive comments and clear function purposes.
- **Pattern Extensibility:** Easy addition of new timestamp format patterns.
- **Configuration Flexibility:** Support for customizing processing parameters.

---

## 5. Architectural Diagram

Below is a high-level diagram emphasizing the processing pipeline and data flow:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Flattened     │    │  Worker Pool    │    │    Input Files  │
│   Log Files     │◀───│  Processing     │◀───│    Discovery    │
│  (Individual)   │    │  (Concurrent)   │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         |                                               
         ▼                                               
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   File Merge    │    │  Chronological  │    │  Format         │
│   Process       │───▶│  Ordering       │───▶│  Restoration    │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
                                              ┌─────────────────┐
                                              │   Final Output  │
                                              │   File          │
                                              │                 │
                                              └─────────────────┘
                                                        │
                                                        ▼
                                              ┌─────────────────┐
                                              │   Cleanup       │
                                              │   (Optional)    │
                                              └─────────────────┘
```

```
flowchart TD
    A["Input Files Discovery"] --> B["Worker Pool Processing (Concurrent)"]
    B --> C["Flattened Log Files (Individual)"]
    C --> D["File Merge Process"]
    D --> E["Chronological Ordering"]
    E --> F["Format Restoration"]
    F --> G["Final Output File"]
    G --> H["Cleanup (Optional)"]
```

**Processing Pipeline:**
1. **File Discovery** → Recursive scan for log files
2. **Pattern Detection** → Identify timestamp format per file
3. **Concurrent Flattening** → Worker pool processes files in parallel
4. **File Merging** → Combine all flattened files
5. **Chronological Sorting** → Extract and sort by timestamps
6. **Format Restoration** → Restore multi-line structure
7. **Output Generation** → Create final timestamped file
8. **Cleanup** → Optionally remove intermediate files

---

## 6. Proposed Technology Stack

- **Programming Language:**  
  - Go 1.24.1+ for cross-platform compatibility and concurrent processing.

- **Core Libraries:**  
  - Standard Go libraries: `bufio`, `regexp`, `time`, `sync`, `path/filepath`
  - No external dependencies for maximum portability.

- **Build Tools:**  
  - Go build system with cross-compilation support.
  - PowerShell scripts for Windows builds.
  - Bash scripts for Linux/macOS builds.

- **Deployment:**  
  - Static binaries with CGO disabled.
  - ZIP packaging for distribution.
  - GitHub release integration.

---

## 7. Milestones & Deliverables

1. **Core Processing Engine:**  
   - Implement file discovery and pattern detection.
   - Develop log flattening and worker pool functionality.
   - Create merge and sorting capabilities.

2. **Command-Line Interface:**  
   - Implement all required and optional command-line flags.
   - Add help and version information.
   - Integrate verbose output and progress reporting.

3. **Cross-Platform Build System:**  
   - Create build scripts for all target platforms.
   - Implement ZIP packaging for distribution.
   - Ensure static binary generation.

4. **Testing and Validation:**  
   - Test with various log file formats and sizes.
   - Validate cross-platform compatibility.
   - Performance testing with large log files.

5. **Documentation and Deployment:**  
   - Complete README with usage examples.
   - Generate release packages for all platforms.
   - Publish to GitHub releases.

---

## 8. Usage Examples

### Basic Usage
```bash
# Process all log files in a directory
./GoMergeLog4Net -p /var/log/myapp

# Use custom number of workers
./GoMergeLog4Net -p /var/log/myapp -workers 8

# Keep intermediate files for debugging
./GoMergeLog4Net -p /var/log/myapp --keep

# Enable verbose output
./GoMergeLog4Net -p /var/log/myapp --verbose
```

### Expected Output
```
/var/log/myapp/ProcessedLogs/20241217_143022_FinalMerged.log
```

### Processing Statistics
```
8 workers finished in 2.3s. Final file: /var/log/myapp/ProcessedLogs/20241217_143022_FinalMerged.log
```

---

## 9. Error Handling and Edge Cases

### File Processing Errors
- **Invalid Log Formats:** Skip files with unrecognized timestamp patterns.
- **Corrupted Files:** Continue processing other files when individual files fail.
- **Permission Issues:** Report access denied errors clearly.

### Resource Management
- **Large Files:** Handle files exceeding available memory through streaming.
- **Disk Space:** Check available space before processing large datasets.
- **File Locks:** Handle files that are currently being written to.

### Performance Considerations
- **Memory Usage:** Optimize buffer sizes for different file sizes.
- **CPU Utilization:** Scale worker count based on available cores.
- **I/O Bottlenecks:** Use buffered operations for optimal throughput.

---

## 10. Future Enhancements

### Potential Extensions
- **Additional Timestamp Formats:** Support for more log timestamp patterns.
- **Filtering Capabilities:** Add options to filter log entries by content or time range.
- **Compression Support:** Handle compressed log files (gzip, bzip2).
- **Real-time Processing:** Support for processing live log streams.
- **GUI Interface:** Optional graphical user interface for non-technical users.

### Configuration Options
- **Custom Output Formats:** Allow specification of output file naming patterns.
- **Processing Rules:** Configurable rules for handling different log entry types.
- **Performance Tuning:** Advanced options for memory and CPU usage optimization.

---

## 11. Summary

GoMergeLog4Net is a robust, cross-platform command-line application designed for efficient processing of Log4Net and similarly formatted log files. The application provides:

- **Efficient Processing:** Concurrent file processing with configurable worker pools.
- **Flexible Input:** Support for multiple timestamp formats and file patterns.
- **Reliable Output:** Chronologically ordered, properly formatted log files.
- **Cross-Platform Compatibility:** Static binaries for Linux, macOS, and Windows.
- **Zero Dependencies:** Self-contained application with no external requirements.

The modular architecture ensures maintainability and extensibility, while the concurrent processing design provides optimal performance for large-scale log processing tasks. The application serves as an essential tool for log analysis, debugging, and system monitoring workflows.

--- 