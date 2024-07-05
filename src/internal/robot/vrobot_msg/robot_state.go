package vrobot_msg

// RobotCommand enum
type RobotCommand string

func (r RobotCommand) String() string {
	return string(r)
}

// RobotCommand enum values
const (
	UNLOADING             RobotCommand = "UNLOADING"
	PICKING               RobotCommand = "PICKING"
	TRAVELING             RobotCommand = "TRAVELING"
	WAITING               RobotCommand = "WAITING"
	WAITING_FOR_UNLOADING RobotCommand = "WAITING_FOR_UNLOADING"
)

type RobotOperation string

func (r RobotOperation) String() string {
	return string(r)
}

const (
	UNPAUSED RobotOperation = "UNPAUSED"
)

type RobotReportState string

func (r RobotReportState) String() string {
	return string(r)
}

const (
	STARTED_UNLOADING             RobotReportState = "STARTED_UNLOADING"
	END_UNLOADING                 RobotReportState = "UNLOADING"
	STARTED_PICKING               RobotReportState = "STARTED_PICKING"
	END_PICKING                   RobotReportState = "PICKING"
	STARTED_TRAVELING             RobotReportState = "STARTED_TRAVELING"
	END_TRAVELING                 RobotReportState = "ARRIVED_AT_POINT"
	STARTED_WAITING               RobotReportState = "STARTED_WAITING"
	END_WAITING                   RobotReportState = "WAITING"
	STARTED_WAITING_FOR_UNLOADING RobotReportState = "STARTED_WAITING_FOR_UNLOADING"
	END_WAITING_FOR_UNLOADING     RobotReportState = "WAITING_FOR_UNLOADING"

	EMERGENCY_STOPPED             RobotReportState = "EMERGENCY_STOPPED"
	RECOVERED_FROM_EMERGENCY_STOP RobotReportState = "RECOVERED_FROM_EMERGENCY_STOP"
	FAILED_TO_UNPAUSE             RobotReportState = "FAILED_TO_UNPAUSE"
)

func NextReports(state RobotCommand) []RobotReportState {
	switch state {
	case UNLOADING:
		return []RobotReportState{STARTED_UNLOADING, END_UNLOADING}
	case PICKING:
		return []RobotReportState{STARTED_PICKING, END_PICKING}
	case TRAVELING:
		return []RobotReportState{STARTED_TRAVELING, END_TRAVELING}
	case WAITING:
		return []RobotReportState{STARTED_WAITING, END_WAITING}
	case WAITING_FOR_UNLOADING:
		return []RobotReportState{STARTED_WAITING_FOR_UNLOADING, END_WAITING_FOR_UNLOADING}
	default:
		return []RobotReportState{}
	}
}
