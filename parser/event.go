package parser

import (
	"strconv"
	"strings"
)

// TODO(pedro.silva) normalize data to lower case
var logHeaderByStrHeader = map[string]LogHeader{
	"Item:":                  LHItem,
	"Kill:":                  LHKill,
	"ClientConnect:":         LHClientConnect,
	"InitGame:":              LHInitGame,
	"Exit:":                  LHExit,
	"ShutdownGame:":          LHShutdownGame,
	"ClientUserinfoChanged:": LHClientUserinfoChanged,
	"ClientBegin:":           LHClientBegin,
	"ClientDisconnect:":      LHClientDisconnect,
	"score:":                 LHScore,
	"say:":                   LHSay,
	"------------------------------------------------------------": LHLogDivision,
}

func GetLogHeader(header string) LogHeader {
	if h, ok := logHeaderByStrHeader[header]; ok {
		return h
	} else {
		return LHUnknown
	}
}

func getEvent(line string) (*Event[any], error) {
	words := strings.Fields(line)
	if len(words) > 1 {
		time := words[0]
		logHeader := words[1]

		switch GetLogHeader(logHeader) {
		case LHItem:
			if len(words) != 4 {
				return &Event[any]{}, &SyntaxError{LHItem, "expecting 4 words on log line"}
			}

			itemSplit := strings.Split(words[3], "_")
			if len(itemSplit) < 2 {
				return &Event[any]{}, &SyntaxError{LHItem, "expecting Item type to have at least 2 words"}
			}

			amount, err := strconv.Atoi(words[2])
			if err != nil {
				return &Event[any]{}, &SyntaxError{LHItem, "expecting Item amount to be an integer"}
			}

			item := Item{
				Amount:   amount,
				Category: itemSplit[0],
				Name:     itemSplit[1],
			}

			if len(itemSplit) > 2 {
				item.SubType = itemSplit[2]
			}

			return &Event[any]{
				HeaderType: LHItem,
				Time:       time,
				Data:       item,
			}, nil

		case LHKill:
			if len(words) < 5 {
				return &Event[any]{}, &SyntaxError{LHKill, "expecting more than 5 words on log line"}
			}
			kill := Kill{}
			var buffer []string
			for i := 5; i < len(words); i++ {
				switch words[i] {
				case "by":
					kill.Victim = strings.Join(buffer, " ")
					buffer = nil
				case "killed":
					kill.Killer = strings.Join(buffer, " ")
					buffer = nil
				default:
					buffer = append(buffer, words[i])
				}
			}
			kill.Means = strings.Join(buffer, " ")
			buffer = nil

			if kill.Means == "" {
				return &Event[any]{}, &SyntaxError{LHKill, "missing Means on log line"}
			}
			if kill.Victim == "" {
				return &Event[any]{}, &SyntaxError{LHKill, "missing Victim on log line"}
			}
			if kill.Killer == "" {
				return &Event[any]{}, &SyntaxError{LHKill, "missing Killer on log line"}
			}

			return &Event[any]{
				HeaderType: LHKill,
				Time:       time,
				Data:       kill,
			}, nil

		case LHClientConnect:
			if len(words) != 3 {
				return &Event[any]{}, &SyntaxError{LHClientConnect, "expecting 3 words on on log line"}
			}

			clientId, err := strconv.Atoi(words[2])
			if err != nil {
				return &Event[any]{}, &SyntaxError{LHClientConnect, "expecting clientId to be an integer"}
			}

			clientConnect := ClientConnect{
				ClientId: clientId,
			}

			return &Event[any]{
				HeaderType: LHClientConnect,
				Time:       time,
				Data:       clientConnect,
			}, nil

		case LHInitGame:
			// TODO(pedro.silva) join words from index 2 to last to parse server init data
			return &Event[any]{
				HeaderType: LHInitGame,
				Time:       time,
				Data:       InitGame{},
			}, nil

		case LHExit:
			return &Event[any]{
				HeaderType: LHExit,
				Time:       time,
				Data: Exit{
					Reason: strings.Join(words[2:], " "),
				},
			}, nil

		case LHShutdownGame:
			return &Event[any]{
				HeaderType: LHShutdownGame,
				Time:       time,
				Data:       ShutdownGame{},
			}, nil

		case LHClientUserinfoChanged:
			if len(words) < 4 {
				return &Event[any]{}, &SyntaxError{LHClientUserinfoChanged, "expecting more than 4 words on log line"}
			}

			clientId, err := strconv.Atoi(words[2])
			if err != nil {
				return &Event[any]{}, &SyntaxError{LHClientUserinfoChanged, "expecting clientId to be an integer"}
			}

			// (pedro.silva) refactor to parse function
			username := strings.Split(strings.Join(words[3:], " "), "\\")[1]
			clientUserinfoChanged := ClientUserinfoChanged{
				ClientId: clientId,
				Username: username,
			}

			return &Event[any]{
				HeaderType: LHClientUserinfoChanged,
				Time:       time,
				Data:       clientUserinfoChanged,
			}, nil

		case LHClientBegin:
			if len(words) != 3 {
				return &Event[any]{}, &SyntaxError{LHClientBegin, "expecting 3 words on log line"}
			}

			clientId, err := strconv.Atoi(words[2])
			if err != nil {
				return &Event[any]{}, &SyntaxError{LHClientBegin, "expecting clientId to be an integer"}
			}

			clientBegin := ClientBegin{
				ClientId: clientId,
			}

			return &Event[any]{
				HeaderType: LHClientBegin,
				Time:       time,
				Data:       clientBegin,
			}, nil

		case LHClientDisconnect:
			if len(words) != 3 {
				return &Event[any]{}, &SyntaxError{LHClientDisconnect, "expecting 3 words on log line"}
			}

			clientId, err := strconv.Atoi(words[2])
			if err != nil {
				return &Event[any]{}, &SyntaxError{LHClientDisconnect, "expecting clientId to be an integer"}
			}

			clientDisconnect := ClientDisconnect{
				ClientId: clientId,
			}

			return &Event[any]{
				HeaderType: LHClientDisconnect,
				Time:       time,
				Data:       clientDisconnect,
			}, nil

		case LHLogDivision:
			return &Event[any]{
				HeaderType: LHLogDivision,
				Time:       time,
				Data:       nil,
			}, nil

		case LHScore:
			if len(words) < 8 {
				return &Event[any]{}, &SyntaxError{LHScore, "expecting 8 or more words on log line"}
			}

			score, err := strconv.Atoi(words[2])
			if err != nil {
				return &Event[any]{}, &SyntaxError{LHScore, "expecting points to be an integer"}
			}

			ping, err := strconv.Atoi(words[4])
			if err != nil {
				return &Event[any]{}, &SyntaxError{LHScore, "expecting ping to be an integer"}
			}

			clientId, err := strconv.Atoi(words[6])
			if err != nil {
				return &Event[any]{}, &SyntaxError{LHScore, "expecting clientId to be an integer"}
			}

			username := strings.Join(words[7:], " ")

			s := Score{
				Score:    score,
				Ping:     ping,
				ClientId: clientId,
				Username: username,
			}

			return &Event[any]{
				HeaderType: LHScore,
				Time:       time,
				Data:       s,
			}, nil

		case LHSay:
			return &Event[any]{
				HeaderType: LHSay,
				Time:       time,
				Data:       Say{},
			}, nil

		case LHUnknown:
			return &Event[any]{}, &SyntaxError{LHUnknown, "unknown log header"}
		}
	}

	return &Event[any]{}, &SyntaxError{LHUnknown, " header could not be found"}
}

type LogHeader uint8

const (
	LHUnknown LogHeader = iota
	LHItem
	LHKill
	LHClientConnect
	LHInitGame
	LHExit
	LHShutdownGame
	LHClientUserinfoChanged
	LHClientBegin
	LHClientDisconnect
	LHLogDivision
	LHScore
	LHSay
)

type Event[T any] struct {
	HeaderType LogHeader
	Time       string
	Data       T
}

type Item struct {
	Amount   int
	Category string
	Name     string
	SubType  string
}

type Kill struct {
	Victim string
	Killer string
	Means  string
}

type ClientConnect struct {
	ClientId int
}

type InitGame struct{}

type ShutdownGame struct{}

type ClientBegin struct {
	ClientId int
}

type ClientDisconnect struct {
	ClientId int
}

type Exit struct {
	Reason string
}

type ClientUserinfoChanged struct {
	ClientId int
	Username string
}

type Score struct {
	Score    int
	Ping     int
	ClientId int
	Username string
}

type Say struct{}
