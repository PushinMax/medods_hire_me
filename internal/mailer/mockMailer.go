package mailer

type MockMailer struct {
	sentEmails []Email
}

func newMockMailer() *MockMailer {
	return &MockMailer{
		sentEmails: make([]Email, 0),
	}
}

func (m *MockMailer) Send(to, subject, body string) error {
	m.sentEmails = append(m.sentEmails, Email{
		To:      to,
		Subject: subject,
		Body:    body,
	})
	return nil
}

func (m *MockMailer) GetHistory() []Email {
	return m.sentEmails
}