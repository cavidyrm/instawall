package usecase

type EmailSender interface {
	Send(to, subject, body string) error
}
