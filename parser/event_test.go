package parser

import (
	"reflect"
	"testing"
)

func TestGetLogHeader(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected LogHeader
	}{
		"LHItem": {
			input:    "Item:",
			expected: LHItem,
		},
		"LHKill": {
			input:    "Kill:",
			expected: LHKill,
		},
		"LHClientConnect": {
			input:    "ClientConnect:",
			expected: LHClientConnect,
		},
		"LHInitGame": {
			input:    "InitGame:",
			expected: LHInitGame,
		},
		"LHExit": {
			input:    "Exit:",
			expected: LHExit,
		},
		"LHShutdownGame": {
			input:    "ShutdownGame:",
			expected: LHShutdownGame,
		},
		"LHClientUserinfoChanged": {
			input:    "ClientUserinfoChanged:",
			expected: LHClientUserinfoChanged,
		},
		"LHClientBegin": {
			input:    "ClientBegin:",
			expected: LHClientBegin,
		},
		"LHClientDisconnect": {
			input:    "ClientDisconnect:",
			expected: LHClientDisconnect,
		},
		"LHLogDivision": {
			input:    "------------------------------------------------------------",
			expected: LHLogDivision,
		},
		"LHScore": {
			input:    "score:",
			expected: LHScore,
		},
		"LHSay": {
			input:    "say:",
			expected: LHSay,
		},
		"LHUnknown": {
			input:    "UnknownHeader",
			expected: LHUnknown,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := GetLogHeader(test.input)
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestGetEvent(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected *Event[any]
		err      error
	}{
		"ValidItemEvent": {
			input: "  20:42 Item: 2 item_armor_body",
			expected: &Event[any]{
				HeaderType: LHItem,
				Time:       "20:42",
				Data: Item{
					Amount:   2,
					Category: "item",
					Name:     "armor",
					SubType:  "body",
				},
			},
			err: nil,
		},
		"InvalidItemEventAmount": {
			input:    "10: Item: aaa invalid_format",
			expected: &Event[any]{},
			err:      &SyntaxError{LHItem, "expecting Item amount to be an integer"},
		},
		"InvalidItemEventCount": {
			input:    "10: Item: invalidFormat",
			expected: &Event[any]{},
			err:      &SyntaxError{LHItem, "expecting 4 words on log line"},
		},
		"InvalidItemEventFormat": {
			input:    "10: Item: 3 invalidFormat",
			expected: &Event[any]{},
			err:      &SyntaxError{LHItem, "expecting Item type to have at least 2 words"},
		},
		"ValidKillEvent": {
			input: " 20:54 Kill: 1022 2 22: killer killed victim by means",
			expected: &Event[any]{
				HeaderType: LHKill,
				Time:       "20:54",
				Data: Kill{
					Killer: "killer",
					Victim: "victim",
					Means:  "means",
				},
			},
			err: nil,
		},
		"InvalidKillEventCount": {
			input:    "15: Kill: invalid_format",
			expected: &Event[any]{},
			err:      &SyntaxError{LHKill, "expecting more than 5 words on log line"},
		},
		"InvalidKillEventMissV": {
			input:    "20:54 Kill: 1022 2 22: killer victim by means",
			expected: &Event[any]{},
			err:      &SyntaxError{LHKill, "missing Victim on log line"},
		},
		"InvalidKillEventMissK": {
			input:    "20:54 Kill: 1022 2 22: killer killed victim means",
			expected: &Event[any]{},
			err:      &SyntaxError{LHKill, "missing Killer on log line"},
		},
		"InvalidKillEventMissM": {
			input:    "20:54 Kill: 1022 2 22: killer killed victim by",
			expected: &Event[any]{},
			err:      &SyntaxError{LHKill, "missing Means on log line"},
		},
		"ValidClientConnectEvent": {
			input: " 20:34 ClientConnect: 2",
			expected: &Event[any]{
				HeaderType: LHClientConnect,
				Time:       "20:34",
				Data: ClientConnect{
					ClientId: 2,
				},
			},
			err: nil,
		},
		"InvalidClientConnectEventId": {
			input:    "20:34 ClientConnect: aa",
			expected: &Event[any]{},
			err:      &SyntaxError{LHClientConnect, "expecting clientId to be an integer"},
		},
		"InvalidClientConnectEventCount": {
			input:    "20:34 ClientConnect: aa aa aa",
			expected: &Event[any]{},
			err:      &SyntaxError{LHClientConnect, "expecting 3 words on on log line"},
		},
		"ValidShutdownGameEvent": {
			input: "25: ShutdownGame:",
			expected: &Event[any]{
				HeaderType: LHShutdownGame,
				Time:       "25:",
				Data:       ShutdownGame{},
			},
			err: nil,
		},
		"ValidClientBeginEvent": {
			input: "30: ClientBegin: 2",
			expected: &Event[any]{
				HeaderType: LHClientBegin,
				Time:       "30:",
				Data: ClientBegin{
					ClientId: 2,
				},
			},
			err: nil,
		},
		"ValidClientDisconnectEvent": {
			input: "35: ClientDisconnect: 3",
			expected: &Event[any]{
				HeaderType: LHClientDisconnect,
				Time:       "35:",
				Data: ClientDisconnect{
					ClientId: 3,
				},
			},
			err: nil,
		},
		"ValidLogDivisionEvent": {
			input: "40: ------------------------------------------------------------",
			expected: &Event[any]{
				HeaderType: LHLogDivision,
				Time:       "40:",
				Data:       nil,
			},
			err: nil,
		},
		"ValidScoreEvent": {
			input: " 11:57 score: 10 ping: 50 client: 4 username with spaces",
			expected: &Event[any]{
				HeaderType: LHScore,
				Time:       "11:57",
				Data: Score{
					Score:    10,
					Ping:     50,
					ClientId: 4,
					Username: "username with spaces",
				},
			},
			err: nil,
		},
		"InvalidScoreEvent": {
			input:    " 11:57 score: invalid ping: 50 client: 4 username with spaces",
			expected: &Event[any]{},
			err:      &SyntaxError{LHScore, "expecting score to be an integer"},
		},
		"ValidSayEvent": {
			input: "11:57 say: asdasd",
			expected: &Event[any]{
				HeaderType: LHSay,
				Time:       "11:57",
				Data:       Say{},
			},
			err: nil,
		},
		"UnknownLogHeader": {
			input:    "55: UnknownHeader: some data",
			expected: &Event[any]{},
			err:      &SyntaxError{LHUnknown, "unknown log header"},
		},
		"HeaderNotFound": {
			input:    "60: NonExistentHeader: some data",
			expected: &Event[any]{},
			err:      &SyntaxError{LHUnknown, " header could not be found"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := getEvent(test.input)
			if err != nil && test.err == nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if err == nil && test.err != nil {
				t.Errorf("Expected error, but got nil")
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %+v, but got %+v", test.expected, result)
			}
		})
	}
}
