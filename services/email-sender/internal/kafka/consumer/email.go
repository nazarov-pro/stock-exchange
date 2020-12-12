package consumer

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/internal/config"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/domain/pb"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/domain"

)

// ConsumeEmails kafka email consumers
func ConsumeEmails(svc domain.Service, ctx context.Context) {
	r := config.NewReader("email", "email-consumer")
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
			sender, err := svc.Send(ctx, &msg)
			if err != nil {
				fmt.Printf("Error occurred, %v\n", err)
				status = domain.Failed
			}

			err = svc.Save(ctx, &msg, sender, status)
			if err != nil {
				fmt.Printf("Error occurred, %v\n", err)
			}
		}
		fmt.Printf("message at offset %d: key: %s \n", m.Offset, string(m.Key))
	}
}
