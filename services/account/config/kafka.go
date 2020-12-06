package config

import (
	"time"

	"github.com/segmentio/kafka-go"
)

var dialer = newDialer()

func newDialer() *kafka.Dialer {
	return &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}
}

func NewWriter(topicName string) *kafka.Writer {
	brokerUrls := []string{"localhost:9092"}
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokerUrls,
		Topic:    topicName,
		Balancer: &kafka.LeastBytes{},
		Dialer:   dialer,
	})
	return w
}

func NewReader(topicName string, groupId string) *kafka.Reader {
	brokerUrls := []string{"localhost:9092"}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerUrls,
		GroupID: groupId,
		Topic:   topicName,
		Dialer:  dialer,
	})
	return r
}

func CreateTopic(topicName string, partitions int, replicas int) error {
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.CreateTopics(
		kafka.TopicConfig{
			Topic:             topicName,
			NumPartitions:     partitions,
			ReplicationFactor: replicas,
		},
	)
}

// func init() {
// 	err := CreateTopic("test", 8, 1)
// 	if err != nil {
// 		panic(err)
// 	}
// }
