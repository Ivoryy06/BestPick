# BestPick - AI-Powered Image Picker

> *Automatically find and pick the best image from a set of duplicates.*

BestPick analyzes images for quality metrics and presents the top 2 candidates for user selection.

## What it does

- **Quality scoring** - blur detection, noise analysis, resolution comparison
- **Duplicate grouping** - groups similar/duplicate images  
- **Interactive picker** - compare top 2 candidates with detailed stats
- **HTML report** - generates beautiful web interface for browser-based selection
- **Local processing** - no API calls, runs entirely on your machine

## Requirements

### Linux/macOS

**Build:**
- Go 1.21+ (`go version`)

**Runtime:**
- None - completely self-contained

### Windows

**Build:**
- Go 1.21+ (download from [go.dev](https://go.dev/dl/))

**Runtime:**
- None - completely self-contained
- No external dependencies

## Installation

### Linux/macOS

```bash
# Clone the repository
git clone https://github.com/Ivoryy06/bestpick.git ~/bestpick

# Build the binary
cd ~/bestpick
go build -o bestpick .

# Install globally (optional)
sudo cp bestpick /usr/local/bin/
```

### Windows

1. **Download Go** from [go.dev/dl](https://go.dev/dl/) and install

2. **Open Command Prompt** (Win+X → Terminal)

3. **Clone and build:**
```cmd
git clone https://github.com/Ivoryy06/bestpick.git
cd bestpick
go build -o bestpick.exe .
```

4. **Run from any location:**
```cmd
move bestpick.exe C:\Windows\System32\
```

## Running

### Linux/macOS

```bash
# Interactive picker
./bestpick -path ~/Pictures -pick

# Generate HTML report
./bestpick -path ~/Pictures -html report.html

# Scan directory
./bestpick -path ~/Pictures

# Exact duplicates only (faster)
./bestpick -path ~/Pictures -quality-only
```

### Windows

Open **Command Prompt** (not PowerShell required):

```cmd
bestpick.exe -path "C:\Users\YourName\Pictures" -pick

bestpick.exe -path "C:\Users\YourName\Pictures" -html report.html

bestpick.exe -path "C:\Users\YourName\Pictures"
```

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-path` | Directory to scan | Current directory |
| `-threshold` | Perceptual distance (lower = stricter) | 10 |
| `-quality-only` | Exact duplicates only | false |
| `-json` | Output JSON file | No output |
| `-html` | Output HTML report | No output |
| `-pick` | Interactive picker | false |

## Quality Metrics

| Metric | Description |
|--------|-------------|
| **Blur Score** | Edge sharpness - higher is better |
| **Noise Score** | Grain level - lower is better |
| **Resolution** | Pixel dimensions |
| **File Size** | Actual file size |
| **Quality** | Combined score (0-100) |

## Quick Examples

```bash
# Find duplicates in your Pictures folder
./bestpick -path ~/Pictures

# High threshold for very similar images only
./bestpick -path ~/Pictures -threshold 5

# Output to HTML for browser viewing
./bestpick -path ~/Pictures -html duplicates.html

# Interactive selection (recommended)
./bestpick -path ~/Pictures -pick
```

## Troubleshooting

- **No images found?** - Make sure path is correct and contains supported formats (JPEG, PNG, WEBP)
- **Too many groups?** - Increase threshold: `-threshold 15`
- **Not finding similar enough images?** - Lower threshold: `-threshold 5`
- **Slow on large folders?** - Use `-quality-only` flag for exact duplicates only

## Project Structure

```
bestpick/
├── main.go      # Core analysis engine
├── picker.go    # Interactive picker UI  
├── index.html  # Web UI template
├── go.mod      # Go module
├── README.md   # This file
└── LICENSE    # MIT License
```

## License

MIT