package mailer

type MailApi interface {	
	Send(to, subject, body string) error
	GetHistory() []Email
}

type Mailer struct {
	MailApi
}

func New() *Mailer {
	return &Mailer{
		MailApi: newMockMailer(),
	}
}

type Email struct {
	To      string
	Subject string
	Body    string
}