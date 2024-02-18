package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pedroegsilva/cw-test/parser"
	"github.com/pedroegsilva/cw-test/reports"
)

func main() {
	// Open the file
	file, err := os.Open("/home/pedro/repos/cw-test/input/qgames.log")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	gameScanner := parser.InitScanner(scanner)

	idx := 0
	for game, ok, err := gameScanner.GetGame(); ok; game, ok, err = gameScanner.GetGame() {
		if err != nil {
			fmt.Println("err", err)
			break
		}
		idx++
		reports.PrintHumanReadableReport(game, fmt.Sprintf("game-%d", idx))
	}
}
