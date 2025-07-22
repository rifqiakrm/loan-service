package email

import (
	"fmt"
	"log"
)

// MockEmailSender simulates sending emails by printing to stdout.
// Useful for testing and development without real SMTP or API services.
type MockEmailSender struct{}

// NewMockEmailSender creates a new instance of MockEmailSender.
func NewMockEmailSender() *MockEmailSender {
	return &MockEmailSender{}
}

// SendInvestorNotification logs an email sent to an investor containing the agreement link.
//
// Example log:
//
//	[EMAIL] Sent an agreement link to investor INV001: https://agreement-link.com/doc.pdf
func (m *MockEmailSender) SendInvestorNotification(investorID, agreementLink string) error {
	log.Printf("[EMAIL] Sent agreement link to investor %s: %s", investorID, agreementLink)
	fmt.Printf("[EMAIL] Sent agreement link to investor %s: %s\n", investorID, agreementLink)
	return nil
}
