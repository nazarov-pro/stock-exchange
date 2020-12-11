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
func ConsumeEmails(svc domain.Service) {
	r := config.NewReader("email", "email-consumer")
	defer r.Close()

	for {
		ctx := context.Background()
		m, err := r.ReadMessage(ctx)
		if err != nil {
			fmt.Printf("Something went wrong err: %v", err)
			break
		}
		var msg pb.SendEmail
		err = proto.Unmarshal(m.Value, &msg)
		if err != nil {
			fmt.Printf("Something went wrong err: %v", err)
		} else {
			status := domain.Sent
			sender, err := svc.Send(ctx, &msg)
			if err != nil {
				fmt.Errorf("Error occurred, %v", err)
				status = domain.Failed
			}

			err = svc.Save(ctx, &msg, sender, status)
			if err != nil {
				fmt.Errorf("Error occurred, %v", err)
			}
		}
		fmt.Printf("message at offset %d: key: %s \n", m.Offset, string(m.Key))
	}
}
