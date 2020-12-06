package kafka

import(
	"context"
	"github.com/nazarov-pro/stock-exchange/services/account/config"
	"github.com/segmentio/kafka-go"
)

func SendEmail(key string, val []byte) (error) {
	w := config.NewWriter("test")
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