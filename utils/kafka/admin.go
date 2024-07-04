package kafka

import (
	"context"
	"fakeflody-agent/logger"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

func CreateTopicIfNotExists(adminClient *kafka.AdminClient, topics []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, topic := range topics {
		// 토픽이 존재하는지 확인
		metadata, err := adminClient.GetMetadata(&topic, false, 10000)
		if err != nil {
			// 메타데이터를 가져오지 못할 경우 에러 처리
			kafkaErr, ok := err.(kafka.Error)
			if ok && kafkaErr.Code() == kafka.ErrUnknownTopicOrPart {
				logger.Errorf("Topic %s does not exist, will create it.", topic)
			} else {
				return fmt.Errorf("failed to get metadata: %w", err)
			}
		} else if _, exists := metadata.Topics[topic]; exists {
			// 토픽이 존재하면 바로 리턴
			//logger.Infof("Topic %s already exists", topic)
			continue
		}

		// 토픽 생성
		topicSpecification := kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}

		// 토픽 생성 요청
		results, err := adminClient.CreateTopics(ctx, []kafka.TopicSpecification{topicSpecification})
		if err != nil {
			return fmt.Errorf("failed to create topic: %w", err)
		}
		// 결과 확인
		for _, result := range results {
			if result.Error.Code() != kafka.ErrNoError {
				return fmt.Errorf("failed to create topic %s: %v", result.Topic, result.Error)
			}
		}
		logger.Infof("Topic %s created successfully\n", topic)
	}

	return nil
}
