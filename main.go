package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pedroegsilva/cw-test/parser"
	"github.com/pedroegsilva/cw-test/reports"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	file, err := os.Open("/home/pedro/repos/cw-test/input/qgames.log")
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
		reports.PrintHumanReadableReport(game, fmt.Sprintf("game-%d", idx))
		reports.Printjson(game, fmt.Sprintf("game-%d", idx))
	}
}
