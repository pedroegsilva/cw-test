package parser

import (
	"bufio"
	"fmt"
	"strings"
)

func InitScanner(scanner *bufio.Scanner) *GameScanner {
	return &GameScanner{
		Scanner:            scanner,
		buffer:             nil,
		clientIdByUsername: make(map[string]int),
	}
}

func (gs *GameScanner) GetGame() (*Game, bool, error) {
	var game *Game = nil
	for event, ok := gs.scan(); ok; event, ok = gs.scan() {
		switch event.HeaderType {
		case LHKill:
			if game == nil {
				return nil, true, fmt.Errorf("kill from empty game")
			}
			kill, ok := event.Data.(Kill)
			if !ok {
				return nil, true, fmt.Errorf("bad parse of event. %s", LHKill.String())
			}

			var kInfo *PlayersInfo
			if kId, ok := gs.clientIdByUsername[kill.Killer]; !ok {
				return nil, true, fmt.Errorf("could not find killer id. username: %s", kill.Killer)
			} else {
				if ki, ok := game.PlayersInfoById[kId]; ok {
					kInfo = ki
				} else {
					return nil, true, fmt.Errorf("could not find killer information. id: %d", kId)
				}
			}

			game.KillCountByMeans[kill.Means]++
			game.TotalKills++

			if kill.Killer == kill.Victim {
				kInfo.SuicideCount++
				kInfo.DeathCount++
				kInfo.DeathCountByWeapon[kill.Means]++
			} else {
				var vInfo *PlayersInfo
				if vId, ok := gs.clientIdByUsername[kill.Victim]; ok {
					if vi, ok := game.PlayersInfoById[vId]; ok {
						vInfo = vi
					} else {
						return nil, true, fmt.Errorf("could not find victim information. id: %d", vId)
					}
				} else {
					return nil, true, fmt.Errorf("could not find victim id. username: %s", kill.Victim)
				}

				vInfo.DeathCount++
				vInfo.DeathCountBySource[kill.Killer]++
				vInfo.DeathCountByWeapon[kill.Means]++
				if kInfo.Id == 0 { // world kill
					vInfo.Score--
				}

				kInfo.KillCount++
				kInfo.Score++
				kInfo.KillCountByPlayerTag[kill.Victim]++
				kInfo.KillCountByMean[kill.Means]++
			}

		case LHInitGame:
			if game != nil {
				// put init event back to buffer
				gs.unScan(event)
				game.EndingReason = "SERVER_UNEXPECTED_SHUTDOWN"
				endGame(game)
				return game, true, nil
			}
			game = &Game{
				PlayersInfoById: map[int]*PlayersInfo{
					0: initPlayerInfo(0),
				},
				KillCountByMeans: initKillCountByMeans(),
			}
			gs.clientIdByUsername["<world>"] = 0
			game.PlayersInfoById[0].Username = "<world>"
		case LHShutdownGame:
			if game == nil {
				return nil, true, fmt.Errorf("shutdown from empty game")
			}
			if game.EndingReason == "" {
				game.EndingReason = "SERVER_UNEXPECTED_SHUTDOWN"
			}
			endGame(game)
			return game, true, nil
		case LHExit:
			if game == nil {
				return nil, true, fmt.Errorf("exit from empty game")
			}
			exit, ok := event.Data.(Exit)
			if !ok {
				return nil, true, fmt.Errorf("bad parse of event. %s", LHExit.String())
			}
			game.EndingReason = exit.Reason
		case LHClientConnect:
			if game == nil {
				return nil, true, fmt.Errorf("client connected to empty game")
			}
			cc, ok := event.Data.(ClientConnect)
			if !ok {
				return nil, true, fmt.Errorf("bad parse of event. %s", LHClientConnect.String())
			}
			game.PlayersInfoById[cc.ClientId] = initPlayerInfo(cc.ClientId)
		case LHClientUserinfoChanged:
			if game == nil {
				return nil, true, fmt.Errorf("user info changed from empty game")
			}
			cuic, ok := event.Data.(ClientUserinfoChanged)
			if !ok {
				return nil, true, fmt.Errorf("bad parse of event. %s", LHClientUserinfoChanged.String())
			}

			if pi, ok := game.PlayersInfoById[cuic.ClientId]; ok {
				pi.Username = cuic.Username
				gs.clientIdByUsername[cuic.Username] = cuic.ClientId
			} else {
				return nil, true, fmt.Errorf("inexistent client change information. id: %d", cuic.ClientId)
			}
		case LHClientDisconnect:
			if game == nil {
				return nil, true, fmt.Errorf("client disconnect from empty game")
			}

			cd, ok := event.Data.(ClientDisconnect)
			if !ok {
				return nil, true, fmt.Errorf("bad parse of event. %s", LHClientDisconnect.String())
			}

			if pi, ok := game.PlayersInfoById[cd.ClientId]; ok {
				game.DisconnectedPlayers = append(game.DisconnectedPlayers, pi)
				delete(game.PlayersInfoById, cd.ClientId)
				delete(gs.clientIdByUsername, pi.Username)
			} else {
				return nil, true, fmt.Errorf("inexistent client disconnected. id: %d", cd.ClientId)
			}
		case LHScore:
		case LHClientBegin:
		case LHItem:
		case LHLogDivision:
		case LHSay:
		}
	}

	// EOF
	endGame(game)
	return game, false, nil
}

func (gs *GameScanner) scan() (*Event[any], bool) {

	if gs.buffer == nil {
		if gs.Scanner.Scan() {
			line := strings.TrimSpace(gs.Scanner.Text())
			event, err := getEvent(line)
			if err != nil {
				//TODO(pedro.silva) Log event error and unknown as a warnning
				// fmt.Println("Error", err)
				// if event.HeaderType == LHUnknown {}
				return gs.scan()
			}
			// Check for errors during scan
			if err := gs.Scanner.Err(); err != nil {
				fmt.Println("Error reading file:", err)
			}
			return event, true
		} else {
			return nil, false
		}
	} else {
		ret := gs.buffer
		gs.buffer = nil
		return ret, true
	}
}

func (gs *GameScanner) unScan(event *Event[any]) {
	gs.buffer = event
}

func endGame(game *Game) {
	if game != nil {
		world := game.PlayersInfoById[0]
		game.WorldKillStatus = WorldKillStatus{
			KillCount:            world.KillCount,
			KillCountByMeans:     world.KillCountByMean,
			KillCountByPlayerTag: world.KillCountByPlayerTag,
		}
		delete(game.PlayersInfoById, 0)
	}
}

func initPlayerInfo(id int) *PlayersInfo {
	return &PlayersInfo{
		Id:                   id,
		KillCount:            0,
		DeathCount:           0,
		SuicideCount:         0,
		DeathCountByWeapon:   make(map[string]int),
		DeathCountBySource:   make(map[string]int),
		KillCountByMean:      make(map[string]int),
		KillCountByPlayerTag: make(map[string]int),
	}
}

func initKillCountByMeans() map[string]int {
	return map[string]int{
		"MOD_UNKNOWN":        0,
		"MOD_SHOTGUN":        0,
		"MOD_GAUNTLET":       0,
		"MOD_MACHINEGUN":     0,
		"MOD_GRENADE":        0,
		"MOD_GRENADE_SPLASH": 0,
		"MOD_ROCKET":         0,
		"MOD_ROCKET_SPLASH":  0,
		"MOD_PLASMA":         0,
		"MOD_PLASMA_SPLASH":  0,
		"MOD_RAILGUN":        0,
		"MOD_LIGHTNING":      0,
		"MOD_BFG":            0,
		"MOD_BFG_SPLASH":     0,
		"MOD_WATER":          0,
		"MOD_SLIME":          0,
		"MOD_LAVA":           0,
		"MOD_CRUSH":          0,
		"MOD_TELEFRAG":       0,
		"MOD_FALLING":        0,
		"MOD_SUICIDE":        0,
		"MOD_TARGET_LASER":   0,
		"MOD_TRIGGER_HURT":   0,
		"MOD_NAIL":           0,
		"MOD_CHAINGUN":       0,
		"MOD_PROXIMITY_MINE": 0,
		"MOD_KAMIKAZE":       0,
		"MOD_JUICED":         0,
		"MOD_GRAPPLE":        0,
	}
}

type PlayersInfo struct {
	Id                   int
	Username             string
	Score                int
	KillCount            int
	DeathCount           int
	SuicideCount         int
	DeathCountByWeapon   map[string]int
	DeathCountBySource   map[string]int
	KillCountByMean      map[string]int
	KillCountByPlayerTag map[string]int
}

type WorldKillStatus struct {
	KillCount            int
	KillCountByMeans     map[string]int
	KillCountByPlayerTag map[string]int
}

type Game struct {
	PlayersInfoById     map[int]*PlayersInfo
	EndingReason        string
	DisconnectedPlayers []*PlayersInfo
	WorldKillStatus     WorldKillStatus
	KillCountByMeans    map[string]int
	TotalKills          int
}

type GameScanner struct {
	Scanner            *bufio.Scanner
	buffer             *Event[any]
	clientIdByUsername map[string]int
}
