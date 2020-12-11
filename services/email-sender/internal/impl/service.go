package impl

import (
	"context"
	"github.com/go-kit/kit/log"

	"github.com/nazarov-pro/stock-exchange/pkg/util/time"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/domain"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/domain/pb"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/internal/mail"
)

// SimpleService represents default service for email
type SimpleService struct {
	repo   domain.Repository
	logger log.Logger
}

// New creates a new instance of Service
func New(repo domain.Repository, logger log.Logger) domain.Service {
	return &SimpleService{
		repo: repo,
		logger: logger,
	}
}

// Save saving message to the database
func (svc *SimpleService) Save(ctx context.Context, msg *pb.SendEmail, 
	sender string, status domain.Status) error {
	message := &domain.Message{
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