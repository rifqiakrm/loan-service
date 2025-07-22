package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockEmailSender_SendInvestorNotification(t *testing.T) {
	sender := NewMockEmailSender()

	investorID := "INV001"
	agreementLink := "AGREEMENT"

	err := sender.SendInvestorNotification(investorID, agreementLink)
	assert.NoError(t, err)
}
