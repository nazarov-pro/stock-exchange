package kafka

import(
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/nazarov-pro/stock-exchange/services/account/pkg/conf"

)

// CreateTopics creating topics
func CreateTopics() {
	conf.CreateTopic(conf.Config.GetString("kafka.topics.email-send"), 1, 1)
}

// SendEmail sening email to specific topic
func SendEmail(key string, val []byte) (error) {
	w := conf.NewWriter(conf.Config.GetString("kafka.topics.email-send"))
	defer w.Close()
	err := w.WriteMessages(
		context.Background(),
		kafka.Message {
			Key: []byte (key),
			Value: val,
		},
	)
	if err != nil {
		return err
	}
	return nil
}