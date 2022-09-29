package mail

import (
	"bytes"
	"encoding/json"
	"github.com/mertozler/internal/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Mail struct {
	from     string
	password string
	host     string
	port     int
}

func NewMail(config *config.Mail) *Mail {
	return &Mail{from: config.From, password: config.Password, port: config.Port, host: config.Host}
}

func (m *Mail) SendMail(to []string, subject string, body string) error {
	for _, receiver := range to {
		logrus.Infof("Sending mail to: %s", receiver)
		msg := gomail.NewMessage()
		msg.SetHeader("From", m.from)
		msg.SetHeader("To", receiver)
		msg.SetHeader("Subject", subject)
		msg.SetBody("application/json", jsonPrettyPrint(body))

		dilaer := gomail.NewDialer(m.host, m.port, m.from, m.password)
		if err := dilaer.DialAndSend(msg); err != nil {
			return err
		}
	}
	return nil
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}
