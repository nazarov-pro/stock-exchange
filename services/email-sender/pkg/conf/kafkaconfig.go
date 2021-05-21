package conf

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

// NewWriter creates new writer for the kafka producer
func NewWriter(topicName string) *kafka.Writer {
	brokerUrls := Config.GetStringSlice("kafka.producer.brokerUrls")
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokerUrls,
		Topic:    topicName,
		Balancer: &kafka.LeastBytes{},
		Dialer:   dialer,
	})
	return w
}

// NewReader creates new reader for the kafka consumer
func NewReader(topicName string, groupID string) *kafka.Reader {
	brokerUrls := Config.GetStringSlice("kafka.consumer.brokerUrls")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerUrls,
		GroupID: groupID,
		Topic:   topicName,
		Dialer:  dialer,
	})
	return r
}

// CreateTopic creates kafka topics regrading topic name, partitions, replicas
func CreateTopic(topicName string, partitions int, replicas int) error {
	conn, err := kafka.Dial(Config.GetString("kafka.admin.network"),
		Config.GetString("kafka.admin.address"))
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
