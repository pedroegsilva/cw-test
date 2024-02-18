package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/pedroegsilva/cw-test/parser"
	"github.com/pedroegsilva/cw-test/reports"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	inputPath := flag.String("i", "", "full path of the log file")
	flag.Parse()

	if *inputPath == "" {
		flag.Usage()
		return
	}

	printJ := os.Getenv("OUT_JSON")
	printH := os.Getenv("OUT_HUMAN")

	file, err := os.Open(*inputPath)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error opening file: %s", err))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	gameScanner := parser.InitScanner(scanner)

	idx := 0
	for game, ok, err := gameScanner.GetGame(); ok; game, ok, err = gameScanner.GetGame() {
		idx++
		if err != nil {
			log.Error().Msg(fmt.Sprintf("error on game-%d: %s", idx, err))
			continue
		}

		if printJ != "" {
			reports.Printjson(game, fmt.Sprintf("game-%d", idx))
		}

		if printH != "" {
			reports.PrintHumanReadableReport(game, fmt.Sprintf("game-%d", idx))
		}
	}
}
