# BestPick - AI-Powered Image Picker

> *Automatically find and pick the best image from a set of duplicates.*

BestPick analyzes images for quality metrics and presents the top 2 candidates for user selection.

## Features

- **Quality scoring** - blur detection, noise analysis, resolution comparison
- **Duplicate grouping** - groups similar/duplicate images
- **Interactive picker** - compare top 2 candidates with detailed stats
- **Web UI** - generates beautiful HTML report for browser-based selection
- **Local processing** - no API calls, runs entirely on your machine

## Quick Start

### Linux/macOS

```bash
# Build from source
go build -o bestpick .

# Run scan
./bestpick -path ~/Pictures

# Generate HTML report
./bestpick -path ~/Pictures -html report.html

# Interactive picker
./bestpick -path ~/Pictures -pick
```

### Windows

```powershell
# Build from source
go build -o bestpick.exe .

# Run scan
.\bestpick.exe -path "C:\Users\YourName\Pictures"

# Generate HTML report
.\bestpick.exe -path "C:\Users\YourName\Pictures" -html report.html

# Interactive picker
.\bestpick.exe -path "C:\Users\YourName\Pictures" -pick
```

## Usage

```bash
# Scan directory for duplicates
./bestpick -path ~/Pictures

# Adjust perceptual threshold
./bestpick -path ~/Pictures -threshold 8

# Exact duplicates only (faster)
./bestpick -path ~/Pictures -quality-only

# Output JSON for integration
./bestpick -path ~/Pictures -json results.json

# Generate interactive HTML report
./bestpick -path ~/Pictures -html report.html

# Interactive terminal picker
./bestpick -path ~/Pictures -pick
```

## Web UI

The HTML output provides an interactive card-based interface:

- Purple/violet theme with dark mode support
- Shows top 2 candidates per group with quality scores
- Click or use buttons to select your preferred image
- Displays blur score, noise, resolution, and file size

## Quality Metrics

| Metric | Description |
|--------|-------------|
| **Blur Score** | Edge sharpness - higher is better |
| **Noise Score** | Grain level - lower is better |
| **Resolution** | Pixel dimensions |
| **File Size** | Actual file size |
| **Quality** | Combined score (0-100) |

## Project Structure

```
.
├── main.go      # Core analysis engine
├── picker.go    # Interactive picker UI
├── index.html   # Web UI template
├── go.mod      # Go module
├── README.md   # This file
└── LICENSE     # MIT License
```

## License

MIT