package mail

import (
	"testing"

	"github.com/nazarov-pro/stock-exchange/services/email-sender/domain/pb"
)

func TestSendEmail(t *testing.T) {
	emailMsg := &pb.SendEmail{
		Content:    "Hello World from go",
		Subject:    "Test",
		Recipients: []string{"me@shahinnazarov.com", "payday@shahinnazarov.com"},
	}
	_, err := SendEmail(emailMsg)
	if err != nil {
		t.Fatalf("Error occured while sending email, %v", err)
	}
}
