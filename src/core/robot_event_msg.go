package core

type ReportedEvent struct {
	Header  Header                 `json:"header"`
	Payload map[string]interface{} `json:"payload"`
}

type DesiredEvent struct {
	Header  Header                 `json:"header"`
	Payload map[string]interface{} `json:"payload"`
}

type Header struct {
	CommandId *string `json:"commandId"`
	RobotId   string  `json:"robotId"`
	Type      string  `json:"type"`
	TimeStamp int64   `json:"timeStamp"`
}

type HeaderType string

func (t HeaderType) String() string {
	return string(t)
}

const (
	REQUEST  HeaderType = "REQUEST"
	OPERATE  HeaderType = "OPERATE"
	RESPONSE HeaderType = "RESPONSE"
	REPORT   HeaderType = "REPORT"
)
