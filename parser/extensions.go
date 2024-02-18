package parser

func (s LogHeader) String() string {
	switch s {
	case LHUnknown:
		return "Unknown"
	case LHItem:
		return "Item"
	case LHClientConnect:
		return "ClientConnect"
	case LHInitGame:
		return "InitGame"
	case LHKill:
		return "Kill"
	case LHExit:
		return "Exit"
	case LHShutdownGame:
		return "ShutdownGame"
	case LHClientUserinfoChanged:
		return "ClientUserinfoChanged"
	case LHClientBegin:
		return "ClientBegin"
	case LHClientDisconnect:
		return "ClientDisconnect"
	case LHLogDivision:
		return "LogDivision"
	case LHScore:
		return "Score"
	case LHSay:
		return "Say"
	}
	return "Unknown"
}
