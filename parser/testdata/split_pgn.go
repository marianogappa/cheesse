package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	inputFile := "pgn_games.pgn"
	outputDir := "games"

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Open input file
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentGame []string
	gameNum := 0
	var outputFile *os.File
	var writer *bufio.Writer

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this is the start of a new game
		if strings.HasPrefix(line, "[Event ") {
			// Close previous game file if open
			if outputFile != nil {
				writer.Flush()
				outputFile.Close()
			}

			// Start new game
			gameNum++
			filename := filepath.Join(outputDir, fmt.Sprintf("game_%05d.pgn", gameNum))
			
			var err error
			outputFile, err = os.Create(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating output file %s: %v\n", filename, err)
				os.Exit(1)
			}
			writer = bufio.NewWriter(outputFile)
			currentGame = []string{line}
		} else if outputFile != nil {
			// Append to current game
			currentGame = append(currentGame, line)
		}

		// Write line to current game file
		if writer != nil {
			fmt.Fprintln(writer, line)
		}
	}

	// Close last game file
	if outputFile != nil {
		writer.Flush()
		outputFile.Close()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Split %d games into %s/\n", gameNum, outputDir)
}
