package vrobot_msg

type BootRobotRequest struct {
	RobotId int    `json:"robotId"`
	Memo    string `json:"memo"`
}

type ShutDownRobotRequest struct {
	RobotId int `json:"robotId"`
}

type GetRobotResult struct {
	RobotId        int      `json:"robotId"`
	RobotName      string   `json:"robotName"`
	Memo           string   `json:"memo"`
	State          string   `json:"state"`
	Estop          bool     `json:"estop"`
	Problems       []string `json:"problems"`
	Solutions      []string `json:"solutions"`
	SessionStartAt string   `json:"sessionStartAt"`
	SessionEndAt   string   `json:"sessionEndAt"`
}
