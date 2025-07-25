package domain

type MockEmailServer struct {
}

func (m *MockEmailServer) SendWelcomeEmail(name, to, subject, html string) error {
	return nil
}
func (m *MockEmailServer) SendSuccessfullyDeletedEmail(to, email, subject, html string) error {
	return nil
}
func (m *MockEmailServer) SendPasswordResetEmail(email, subject, html string) error {
	return nil
}
func (m *MockEmailServer) SendVerificationEmail(name, to, subject, html string) error {
	return nil
}
