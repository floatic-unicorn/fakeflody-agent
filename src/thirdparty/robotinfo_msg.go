package thirdparty

import (
	"time"
)

type RobotInfoResponse struct {
	RobotID  string `json:"robotId"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	RobotJob struct {
		Type         string    `json:"type"`
		RobotTag     int       `json:"robotTag"`
		LocationCode string    `json:"locationCode"`
		LocationName string    `json:"locationName"`
		AssignedAt   time.Time `json:"assignedAt"`
	} `json:"robotJob"`
	WarehouseID string `json:"warehouseId"`
	Position    struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"position"`
	Orientation struct {
		X int `json:"x"`
		Y int `json:"y"`
		Z int `json:"z"`
		W int `json:"w"`
	} `json:"orientation"`
	BatteryPercent int  `json:"batteryPercent"`
	Online         bool `json:"online"`
}
