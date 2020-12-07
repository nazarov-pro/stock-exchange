package implementation

import (
	"context"
	"github.com/go-kit/kit/log"

	"github.com/nazarov-pro/stock-exchange/pkg/util/time"
	"github.com/nazarov-pro/stock-exchange/services/email-sender"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/pb"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/mail"
)

// SimpleService represents default service for email
type SimpleService struct {
	repo   email.Repository
	logger log.Logger
}

// New creates a new instance of Service
func New(repo email.Repository, logger log.Logger) email.Service {
	return &SimpleService{
		repo: repo,
		logger: logger,
	}
}

// Save saving message to the database
func (svc *SimpleService) Save(ctx context.Context, msg *pb.SendEmail, 
	sender string, status email.Status) error {
	message := &email.Message{
		Content: msg.Content,
		Recipients: msg.Recipients,
		Sender: sender,
		CreatedDate: time.Epoch(),
		Subject: msg.Subject,
		Status: status,
	}
	return svc.repo.Insert(ctx, message)
}

// Send sending email message via smtp
func (svc *SimpleService) Send(ctx context.Context, msg *pb.SendEmail) (string, error) {
	return mail.SendEmail(msg)
}