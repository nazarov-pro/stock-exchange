package domain

import (
	"context"

	"github.com/nazarov-pro/stock-exchange/services/email-sender/pkg/domain/pb"
)

// EmailStatus represents status of the email message
type EmailStatus int8

const (
	// Failed failed to send an email
	Failed EmailStatus = iota
	// Sent message successfully sent
	Sent
)

// EmailMessage represents email message
type EmailMessage struct {
	ID             int64
	Subject        string
	Content        string
	Sender         string
	Status         EmailStatus
	Recipients     []string
	CreatedDate    int64
	LastUpdateDate int64
}

//EmailRepository consists of repository for rdbms
type EmailRepository interface {
	Insert(ctx context.Context, message *EmailMessage) error
}

// EmailService represents business logic of email sending
type EmailService interface {
	Save(ctx context.Context, msg *pb.SendEmail, sender string, status EmailStatus) error

	Send(ctx context.Context, msg *pb.SendEmail) (string, error)
}

// EmailConsumer email consumer interface
type EmailConsumer interface {
	Consume(ctx context.Context) error
}
