package thirdparty

import (
	"encoding/json"
	"fakeflody-agent/config"
	"fakeflody-agent/logger"
	"fakeflody-agent/utils/hashids"
	"fmt"
	"io"
	"net/http"
)

type RobotInfoService interface {
	GetRobotInfo(robotId int) *RobotInfoResponse
}

type BelugaRobotInfoService struct {
	baseUrl string
	conf    *config.FakeFlodyConfig
}

func NewRobotInfoService(
	conf *config.FakeFlodyConfig,
) RobotInfoService {
	return &BelugaRobotInfoService{
		baseUrl: conf.BelugaConfig.RobotInfoService.Url,
		conf:    conf,
	}
}

func (svc BelugaRobotInfoService) GetRobotInfo(robotId int) *RobotInfoResponse {
	url := svc.baseUrl + "/v1/robotInfos/" + fmt.Sprint(hashids.ToUid(robotId))
	bearerToken := svc.conf.BelugaConfig.Token

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.WErrorf("Error creating request: %v", err)
		return nil
	}

	req.Header.Add("Authorization", "Bearer "+bearerToken)

	resp, err := client.Do(req)
	if err != nil {
		logger.WErrorf("Error sending request: %v", err)
		return nil
	}
	defer resp.Body.Close()

	// 응답 상태 코드를 출력합니다.

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WErrorf("Error reading response body: %v", err)
		return nil
	}

	var responseData RobotInfoResponse
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		logger.WErrorf("Error unmarshalling response body: %v", err)
		return nil
	}

	return &responseData
}
