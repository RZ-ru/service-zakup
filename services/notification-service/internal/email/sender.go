package email

import "log"

type Sender struct{}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) Send(to, subject, body string) error {
	log.Printf("EMAIL TO=%s SUBJECT=%s BODY=%s", to, subject, body)
	return nil
}
