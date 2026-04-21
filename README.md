# BestPick - AI-Powered Image Picker

> *Automatically find and pick the best image from a set of duplicates.*

BestPick analyzes images for quality metrics and presents the top 2 candidates for user selection.

## Features

- **Quality scoring** - blur detection, noise analysis, resolution comparison
- **Duplicate grouping** - groups similar/duplicate images
- **Interactive picker** - compare top 2 candidates with detailed stats
- **Local processing** - no API calls, runs entirely on your machine

## Quick Start

```bash
go build -o bestpick main.go picker.go
./bestpick -path ~/Pictures
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
```

## Quality Metrics

| Metric | Description |
|--------|-------------|
| **Blur Score** | Edge sharpness - higher is better |
| **Noise Score** | Grain level - lower is better |
| **Resolution** | Pixel dimensions |
| **File Size** | Actual file size |

## Project Structure

```
.
├── main.go      # Core analysis engine
├── picker.go    # Interactive picker UI
├── go.mod      # Go module
└── README.md   # This file
```

## License

MIT