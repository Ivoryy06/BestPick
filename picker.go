package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PickerResult struct {
	Selected    int
	ImagePath   string
	Stats       ImageStats
}

func RunPicker(groups []DuplicateGroup) {
	if len(groups) == 0 {
		fmt.Println("No duplicate groups to pick from.")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== Interactive Image Picker ===")
	fmt.Println("For each group, you'll see the top 2 candidates.")
	fmt.Println("Enter 1 or 2 to select your preferred image, or 's' to skip.\n")

	results := []PickerResult{}

	for i, group := range groups {
		if len(group.Images) < 2 {
			continue
		}

		img1 := group.Images[group.BestPick]
		img2 := group.Images[group.SecondPick]

		fmt.Printf("--- Group %d ---\n\n", i+1)
		printComparison(img1, img2)

		fmt.Print("Your choice (1/2/s): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "s" || input == "skip" {
			fmt.Println("Skipped.\n")
			continue
		}

		choice, err := strconv.Atoi(input)
		if err != nil || (choice != 1 && choice != 2) {
			fmt.Println("Invalid choice, skipping.\n")
			continue
		}

		var selected ImageStats
		if choice == 1 {
			selected = img1
		} else {
			selected = img2
		}

		results = append(results, PickerResult{
			Selected:  choice,
			ImagePath: selected.Path,
			Stats:     selected,
		})

		fmt.Printf("Selected: %s\n\n", selected.Path)
	}

	fmt.Println("\n=== Selection Summary ===")
	for i, r := range results {
		fmt.Printf("%d. %s (Quality: %.2f)\n", i+1, r.ImagePath, r.Stats.Quality)
	}
}

func printComparison(img1, img2 ImageStats) {
	fmt.Println("┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│                        CANDIDATE 1                              │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Printf("│  %-60s  │\n", truncate(img1.Path, 60))
	fmt.Printf("│  Resolution: %-8d  File Size: %-10s                   │\n",
		img1.Resolution, formatFileSize(img1.FileSize))
	fmt.Printf("│  Blur Score: %-8.2f  Quality: %-8.2f/100                 │\n",
		img1.BlurScore, img1.Quality)
	fmt.Printf("│  Noise:      %-8.2f  Dimensions: %dx%d                      │\n",
		img1.NoiseScore, img1.Width, img1.Height)
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│                        CANDIDATE 2                              │")
	fmt.Println("├─────────────────────────────────────────────────────────────────┤")
	fmt.Printf("│  %-60s  │\n", truncate(img2.Path, 60))
	fmt.Printf("│  Resolution: %-8d  File Size: %-10s                   │\n",
		img2.Resolution, formatFileSize(img2.FileSize))
	fmt.Printf("│  Blur Score: %-8.2f  Quality: %-8.2f/100                 │\n",
		img2.BlurScore, img2.Quality)
	fmt.Printf("│  Noise:      %-8.2f  Dimensions: %dx%d                      │\n",
		img2.NoiseScore, img2.Width, img2.Height)
	fmt.Println("└─────────────────────────────────────────────────────────────────┘")
	fmt.Println()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return "..." + s[len(s)-maxLen+3:]
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
