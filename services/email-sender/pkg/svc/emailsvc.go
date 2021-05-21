package svc

import (
	"context"

	"github.com/go-kit/kit/log"

	"github.com/nazarov-pro/stock-exchange/pkg/util/time"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/pkg/domain"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/pkg/domain/pb"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/pkg/mail"
)

// emailService represents default service for email
type emailService struct {
	repo   domain.EmailRepository
	logger log.Logger
}

// New creates a new instance of Service
func New(repo domain.EmailRepository, logger log.Logger) domain.EmailService {
	return &emailService{
		repo:   repo,
		logger: logger,
	}
}

// Save saving message to the database
func (svc *emailService) Save(ctx context.Context, msg *pb.SendEmail,
	sender string, status domain.EmailStatus) error {
	message := &domain.EmailMessage{
		Content:     msg.Content,
		Recipients:  msg.Recipients,
		Sender:      sender,
		CreatedDate: time.Epoch(),
		Subject:     msg.Subject,
		Status:      status,
	}
	return svc.repo.Insert(ctx, message)
}

// Send sending email message via smtp
func (svc *emailService) Send(ctx context.Context, msg *pb.SendEmail) (string, error) {
	return mail.SendEmail(msg)
}
