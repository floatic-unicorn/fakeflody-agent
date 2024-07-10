package thirdparty

import (
	"encoding/json"
	"fakeflody-agent/src/config"
	"fakeflody-agent/src/logger"
	"fakeflody-agent/src/message"
	"fakeflody-agent/src/utils/hashids"
	"fmt"
	"io"
	"net/http"
)

type RobotInfoService interface {
	GetRobotInfo(robotId int) *RobotInfoResponse
	GetRobotInfosByWarehouse(warehouseId int) []*message.GetRobotInfoResult
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
	var responseData RobotInfoResponse
	if err := svc.makeRequest("GET", url, &responseData); err != nil {
		return nil
	}
	return &responseData
}

func (svc BelugaRobotInfoService) GetRobotInfosByWarehouse(warehouseId int) []*message.GetRobotInfoResult {
	url := svc.baseUrl + "/v1/warehouses/" + fmt.Sprint(hashids.ToUid(warehouseId)) + "/robotInfos"
	var responseData []*RobotInfoResponse
	if err := svc.makeRequest("GET", url, &responseData); err != nil {
		return nil
	}

	result := make([]*message.GetRobotInfoResult, len(responseData))
	for i, v := range responseData {
		result[i] = &message.GetRobotInfoResult{
			RobotID: hashids.ToId(v.RobotID),
			Name:    v.Name,
			Status:  v.Status,
		}
	}

	return result
}

func (svc BelugaRobotInfoService) makeRequest(method, url string, responseData interface{}) error {
	bearerToken := svc.conf.BelugaConfig.Token

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		logger.WErrorf("Error creating request: %v", err)
		return err
	}

	req.Header.Add("Authorization", "Bearer "+bearerToken)

	resp, err := client.Do(req)
	if err != nil {
		logger.WErrorf("Error sending request: %v", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.WErrorf("Error reading response body: %v", err)
		return err
	}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		logger.WErrorf("Error unmarshalling response body: %v", err)
		return err
	}

	return nil
}
