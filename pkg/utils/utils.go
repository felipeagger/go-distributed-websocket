package utils

import "fmt"

func GetTopicName(userId, origin string) string {
	return fmt.Sprintf("ws-%s-%s", userId, origin)
}

func SetDefaultEnvs(streamName, consumerGroup string) (string, string) {
	if streamName == "" {
		streamName = "all_messages"
	}

	if consumerGroup == "" {
		consumerGroup = "default"
	}

	return streamName, consumerGroup
}