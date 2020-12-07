package kafka

import(
	"context"
	"github.com/nazarov-pro/stock-exchange/services/account/config"
	"github.com/segmentio/kafka-go"
)

func CreateTopics() {
	config.CreateTopic("email", 1, 1)
}


func SendEmail(key string, val []byte) (error) {
	w := config.NewWriter("email")
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