package main

import (
	"crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/disintegration/gift"
	"github.com/sirupsen/logrus"
)

type ImageStats struct {
	Path       string  `json:"path"`
	Resolution int     `json:"resolution"`
	BlurScore  float64 `json:"blur_score"`
	NoiseScore float64 `json:"noise_score"`
	Quality    float64 `json:"quality"`
	FileSize   int64   `json:"file_size"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
}

type DuplicateGroup struct {
	Hash       string       `json:"hash"`
	Images     []ImageStats `json:"images"`
	BestPick   int          `json:"best_pick"`
	SecondPick int          `json:"second_pick"`
}

var (
	path        string
	threshold   int
	verbose     bool
	outputJSON  string
	qualityOnly bool
	pickMode    bool
)

func init() {
	flag.StringVar(&path, "path", ".", "Directory to scan")
	flag.IntVar(&threshold, "threshold", 10, "Perceptual distance threshold")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.StringVar(&outputJSON, "json", "", "Output results as JSON")
	flag.BoolVar(&qualityOnly, "quality-only", false, "Only find exact duplicates by hash")
	flag.BoolVar(&pickMode, "pick", false, "Interactive picker mode")
}

func main() {
	flag.Parse()

	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Infof("Scanning directory: %s", path)

	files, err := scanImages(path)
	if err != nil {
		logrus.Fatalf("Failed to scan directory: %v", err)
	}

	logrus.Infof("Found %d images", len(files))

	var wg sync.WaitGroup
	statsChan := make(chan ImageStats, len(files))
	sem := make(chan struct{}, 8)

	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			stats, err := analyzeImage(f)
			if err != nil {
				logrus.Warnf("Failed to analyze %s: %v", f, err)
				return
			}
			statsChan <- stats
		}(file)
	}

	go func() {
		wg.Wait()
		close(statsChan)
	}()

	var allStats []ImageStats
	for stats := range statsChan {
		allStats = append(allStats, stats)
	}

	logrus.Infof("Analyzed %d images", len(allStats))

	groups := findDuplicateGroups(allStats, threshold, qualityOnly)

	if pickMode {
		RunPicker(groups)
		return
	}

	if outputJSON != "" {
		data, _ := json.MarshalIndent(groups, "", "  ")
		os.WriteFile(outputJSON, data, 0644)
		logrus.Infof("Results saved to %s", outputJSON)
		return
	}

	printResults(groups)
}

func scanImages(root string) ([]string, error) {
	var files []string
	exts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".webp": true, ".bmp": true, ".tiff": true, ".tif": true,
	}

	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(p))
		if exts[ext] {
			files = append(files, p)
		}
		return nil
	})
	return files, err
}

func analyzeImage(path string) (ImageStats, error) {
	file, err := os.Open(path)
	if err != nil {
		return ImageStats{}, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return ImageStats{}, err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return ImageStats{}, err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	resolution := width * height

	blurScore := calculateBlur(img)
	noiseScore := calculateNoise(img)

	quality := calculateQualityScore(resolution, blurScore, noiseScore, info.Size())

	return ImageStats{
		Path:       path,
		Resolution: resolution,
		BlurScore:  blurScore,
		NoiseScore: noiseScore,
		Quality:    quality,
		FileSize:   info.Size(),
		Width:      width,
		Height:     height,
	}, nil
}

func calculateBlur(img image.Image) float64 {
	g := gift.New(gift.Grayscale())
	gray := image.NewGray(img.Bounds())
	g.Draw(gray, img)

	kernel := []float32{
		-1, -1, -1,
		-1,  8, -1,
		-1, -1, -1,
	}
	conv := gift.New(gift.Convolution(kernel, false, false, false, 0))
	result := image.NewGray(gray.Bounds())
	conv.Draw(result, gray)

	var sum float64
	var count int
	for y := result.Bounds().Min.Y; y < result.Bounds().Max.Y; y++ {
		for x := result.Bounds().Min.X; x < result.Bounds().Max.X; x++ {
			sum += math.Abs(float64(result.GrayAt(x, y).Y) - 128)
			count++
		}
	}

	normalized := sum / float64(count)
	return normalized
}

func calculateNoise(img image.Image) float64 {
	g := gift.New(gift.Grayscale())
	gray := image.NewGray(img.Bounds())
	g.Draw(gray, img)

	bounds := gray.Bounds()
	var variance float64
	var count int

	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			curr := float64(gray.GrayAt(x, y).Y)
			left := float64(gray.GrayAt(x-1, y).Y)
			right := float64(gray.GrayAt(x+1, y).Y)
			up := float64(gray.GrayAt(x, y-1).Y)
			down := float64(gray.GrayAt(x, y+1).Y)

			diff := math.Abs(curr-left) + math.Abs(curr-right) + math.Abs(curr-up) + math.Abs(curr-down)
			variance += diff
			count++
		}
	}

	noise := variance / float64(count)
	return noise
}

func calculateQualityScore(resolution int, blur, noise float64, fileSize int64) float64 {
	resolutionScore := math.Min(float64(resolution)/4000000.0*100, 100)
	blurScore := math.Min(blur*10, 100)
	noiseScore := math.Min((1/(noise+1))*50, 100)
	sizeScore := math.Min(float64(fileSize)/5000000.0*100, 100)

	quality := resolutionScore*0.3 + blurScore*0.3 + noiseScore*0.2 + sizeScore*0.2

	return math.Round(quality*100) / 100
}

func computeHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha1.Sum(data)
	return fmt.Sprintf("%x", hash), nil
}

func findDuplicateGroups(stats []ImageStats, threshold int, qualityOnly bool) []DuplicateGroup {
	groups := make(map[string][]ImageStats)

	if qualityOnly {
		for _, s := range stats {
			hash, err := computeHash(s.Path)
			if err != nil {
				continue
			}
			groups[hash] = append(groups[hash], s)
		}
	} else {
		for i := range stats {
			for j := i + 1; j < len(stats); j++ {
				dist := perceptualDistance(stats[i].Path, stats[j].Path)
				if dist <= threshold {
					key := fmt.Sprintf("%s-%d", stats[i].Path, dist)
					groups[key] = append(groups[key], stats[i], stats[j])
				}
			}
		}
	}

	var result []DuplicateGroup
	for _, images := range groups {
		if len(images) < 2 {
			continue
		}

		seen := make(map[string]bool)
		var unique []ImageStats
		for _, img := range images {
			if !seen[img.Path] {
				seen[img.Path] = true
				unique = append(unique, img)
			}
		}

		if len(unique) < 2 {
			continue
		}

		sort.Slice(unique, func(i, j int) bool {
			return unique[i].Quality > unique[j].Quality
		})

		group := DuplicateGroup{
			Images:     unique,
			BestPick:   0,
			SecondPick: 1,
		}
		result = append(result, group)
	}

	return result
}

func perceptualDistance(path1, path2 string) int {
	hash1, _ := computeHash(path1)
	hash2, _ := computeHash(path2)

	if hash1 == hash2 {
		return 0
	}

	minLen := len(hash1)
	if len(hash2) < minLen {
		minLen = len(hash2)
	}

	distance := 0
	for i := 0; i < minLen; i++ {
		if hash1[i] != hash2[i] {
			distance++
		}
	}

	return distance
}

func printResults(groups []DuplicateGroup) {
	if len(groups) == 0 {
		fmt.Println("\nNo duplicate groups found.")
		return
	}

	fmt.Printf("\n=== Found %d duplicate groups ===\n\n", len(groups))

	for i, group := range groups {
		fmt.Printf("--- Group %d ---\n", i+1)
		fmt.Println("Candidates (sorted by quality):")
		fmt.Println()

		for j, img := range group.Images {
			marker := "  "
			if j == group.BestPick {
				marker = "★ "
			} else if j == group.SecondPick {
				marker = "☆ "
			}

			fmt.Printf("  %s[%d] %s\n", marker, j+1, filepath.Base(img.Path))
			fmt.Printf("      Resolution: %dx%d\n", img.Width, img.Height)
			fmt.Printf("      File Size:  %.2f MB\n", float64(img.FileSize)/1000000)
			fmt.Printf("      Blur Score: %.2f\n", img.BlurScore)
			fmt.Printf("      Noise Score: %.2f\n", img.NoiseScore)
			fmt.Printf("      Quality:    %.2f/100\n", img.Quality)
			fmt.Println()
		}

		fmt.Println("Pick the best: run with -pick group_index")
		fmt.Println()
	}
}
