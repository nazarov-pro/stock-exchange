package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"

	"github.com/nazarov-pro/stock-exchange/services/email-sender/pkg/conf"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/pkg/domain/pb"
)

// SendEmail sending email
func SendEmail(emailMessage *pb.SendEmail) (string, error) {
	// Connect to the SMTP Server
	config := conf.Config
	servername := config.GetString("mail.smtp.server")

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", config.GetString("mail.smtp.username"),
		config.GetString("mail.smtp.password"), host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return "", err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return "", err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return "", err
	}

	from := mail.Address{
		Name:    config.GetString("mail.sender.name"),
		Address: config.GetString("mail.sender.email"),
	}
	// To && From
	if err = c.Mail(from.Address); err != nil {
		return from.Address, err
	}

	to := ""
	for _, item := range emailMessage.Recipients {
		addr := mail.Address{Address: item}
		if err = c.Rcpt(addr.Address); err != nil {
			return from.Address, err
		}

		to += addr.String() + ";"
	}

	subj := emailMessage.Subject
	body := emailMessage.Content

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to
	headers["Subject"] = subj
	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Data
	w, err := c.Data()
	if err != nil {
		return from.Address, err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return from.Address, err
	}

	err = w.Close()
	if err != nil {
		return from.Address, err
	}

	c.Quit()
	return from.Address, nil
}
