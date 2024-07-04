package config

type ReportCommandKeyword string

func (r ReportCommandKeyword) String() string {
	return string(r)
}

const (
	ActionEstop    ReportCommandKeyword = "estop"
	ActionALLEstop ReportCommandKeyword = "all-estop"
	ActionALLClear ReportCommandKeyword = "all-clear"
	ActionClear    ReportCommandKeyword = "clear"
)
