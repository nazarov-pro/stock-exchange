package mail

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"

	"github.com/nazarov-pro/stock-exchange/services/email-sender/internal/config"
	"github.com/nazarov-pro/stock-exchange/services/email-sender/domain/pb"
)

// SendEmail sending email
func SendEmail(emailMessage *pb.SendEmail) (string, error) {
	// Connect to the SMTP Server
	conf := config.Config
	servername := conf.GetString("mail.smtp.server")

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", conf.GetString("mail.smtp.username"),
		conf.GetString("mail.smtp.password"), host)

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
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	from := mail.Address{
		Name:    conf.GetString("mail.sender.name"),
		Address: conf.GetString("mail.sender.email"),
	}
	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	to := ""
	for _, item := range emailMessage.Recipients {
		addr := mail.Address{Address: item}
		if err = c.Rcpt(addr.Address); err != nil {
			log.Panic(err)
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
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
		return from.Address, err
	}

	c.Quit()
	return from.Address, nil
}
