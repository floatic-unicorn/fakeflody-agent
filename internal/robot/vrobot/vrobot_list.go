package vrobot

import (
	"strconv"
)

type VRobotList []VRobot

func NewRobots() VRobotList {
	var robots = make([]VRobot, 0)
	return robots
}

func (list VRobotList) GetRobotIds() []string {
	ids := make([]string, 0)
	for _, robot := range list {

		// int to string
		// robot.RobotId
		ids = append(ids, strconv.Itoa(robot.GetRobotId()))
	}
	return ids
}

func (list VRobotList) GetVRobotById(robotId int) VRobot {
	for _, robot := range list {
		if robot.GetRobotId() == robotId {
			return robot
		}
	}
	return nil
}

func (list VRobotList) GetVRobots() []VRobot {
	return list
}
