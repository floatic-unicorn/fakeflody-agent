package message

type BootRobotRequest struct {
	RobotId  int    `json:"robotId"`
	Memo     string `json:"memo"`
	Interval int    `json:"interval"`
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
	Interval       int      `json:"interval"`
	SessionStartAt int64    `json:"sessionStartAt"`
	SessionEndAt   int64    `json:"sessionEndAt"`
}

type GetRobotInfoResult struct {
	RobotID     int    `json:"robotId"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	WarehouseID int    `json:"warehouseId"`
}
