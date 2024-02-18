package reports

import (
	"fmt"

	"github.com/pedroegsilva/cw-test/parser"
)

type PlayerKills struct {
	Name      string
	KillCount int
}

type PlayerStatistics struct {
	Name           string
	Score          int
	KillCount      int
	FavoritWeapon  string
	Nemesis        string
	TargetPractice string
	Vulnerability  string
}

type Report struct {
	GameIdentifier    string
	TotalKills        int
	EndingReason      string
	PlayersStatistics []*PlayerStatistics
	WorldEnemy        string
}

func getTop(stat *map[string]int) string {
	top := "-"
	topV := 0
	for k, v := range *stat {
		if v > topV {
			topV = v
			top = k
		}
	}
	return top
}

func getPlayerStatistics(players map[int]*parser.PlayersInfo) []*PlayerStatistics {
	statistics := make([]*PlayerStatistics, len(players))
	idx := 0
	for _, info := range players {
		ps := PlayerStatistics{
			Name:           info.Username,
			Score:          info.Score,
			KillCount:      info.KillCount,
			FavoritWeapon:  getTop(&info.KillCountByMean),
			Nemesis:        getTop(&info.DeathCountBySource),
			TargetPractice: getTop(&info.KillCountByPlayerTag),
			Vulnerability:  getTop(&info.DeathCountByWeapon),
		}
		statistics[idx] = &ps
		idx++
	}
	return statistics
}

func createReportStructure(game *parser.Game, name string) *Report {
	return &Report{
		GameIdentifier:    name,
		TotalKills:        game.TotalKills,
		EndingReason:      game.EndingReason,
		WorldEnemy:        getTop(&game.WorldKillStatus.KillCountByPlayerTag),
		PlayersStatistics: getPlayerStatistics(game.PlayersInfoById),
	}
}

func PrintHumanReadableReport(game *parser.Game, name string) {
	report := createReportStructure(game, name)
	fmt.Printf("-------------------- %s --------------------\n", report.GameIdentifier)
	fmt.Println("Total kills: ", report.TotalKills)
	fmt.Println("Game Ending Event: ", report.TotalKills)
	fmt.Println("World Enemy: ", report.WorldEnemy)
	fmt.Println("Player Statistics: ")
	for _, ps := range report.PlayersStatistics {
		fmt.Println("  ", ps.Name)
		fmt.Println("    Score: ", ps.Score)
		fmt.Println("    Kill Count: ", ps.KillCount)
		fmt.Println("    Nemesis: ", ps.Nemesis)
		fmt.Println("    Target Practice: ", ps.TargetPractice)
		fmt.Println("    Favorit Weapon: ", ps.FavoritWeapon)
		fmt.Println("    Vulnerability: ", ps.Vulnerability)
	}
}
