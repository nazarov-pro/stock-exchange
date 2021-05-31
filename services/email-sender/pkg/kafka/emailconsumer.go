package kafka

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/pkg/conf"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/pkg/domain/pb"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/pkg/domain"

)

type emailConsumer struct {
	svc domain.EmailService
}

// NewEmailConsumer New instance of email consumer
func NewEmailConsumer(svc domain.EmailService) domain.EmailConsumer {
	return &emailConsumer{
		svc: svc,
	}
}

// Consume kafka email consumers
func (consumer emailConsumer) Consume(ctx context.Context) error {
	r := conf.NewReader("email", "email-consumer")
	defer func() {
		fmt.Printf("Reader closing err: %v \n", r.Close())
	}()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			fmt.Printf("Something went wrong err: %v\n", err)
			break
		}
		var msg pb.SendEmail
		err = proto.Unmarshal(m.Value, &msg)
		if err != nil {
			fmt.Printf("Something went wrong err: %v\n", err)
		} else {
			status := domain.Sent
			sender, err := consumer.svc.Send(ctx, &msg)
			if err != nil {
				fmt.Printf("Something went wrong err: %v\n", err)
				status = domain.Failed
			}

			err = consumer.svc.Save(ctx, &msg, sender, status)
			if err != nil {
				return err
			}
		}
		fmt.Printf("message at offset %d: key: %s \n", m.Offset, string(m.Key))
	}
	return nil
}
