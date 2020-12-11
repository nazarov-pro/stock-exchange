package domain

import (
	"context"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/domain/pb"
)

// Status represents status of the email message
type Status int8

const (
	// Failed failed to send an email
	Failed Status = iota
	// Sent message successfully sent
	Sent
)

// Message represents email message
type Message struct {
	ID             int64
	Subject        string
	Content        string
	Sender         string
	Status         Status
	Recipients     []string
	CreatedDate    int64
	LastUpdateDate int64
}

//Repository consists of repository for rdbms
type Repository interface {
	Insert(ctx context.Context, message *Message) error
}

// Service represents business logic
type Service interface {
	Save(ctx context.Context, msg *pb.SendEmail, sender string, status Status) error

	Send(ctx context.Context, msg *pb.SendEmail) (string, error)
}