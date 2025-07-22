package loan

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockEmailSender simulates an email sender for testing purposes.
type mockEmailSender struct {
	sentTo []string
}

func (m *mockEmailSender) SendInvestorNotification(investorID, agreementLink string) error {
	m.sentTo = append(m.sentTo, investorID)
	return nil
}

// errorRepo mocks repo with update failure
type errorRepo struct{}

func (e *errorRepo) Create(*Loan) error { return nil }
func (e *errorRepo) GetByID(id string) (*Loan, error) {
	return &Loan{
		ID:    id,
		State: Proposed,
	}, nil
}
func (e *errorRepo) Update(*Loan) error     { return errors.New("forced update error") }
func (e *errorRepo) List() ([]*Loan, error) { return nil, nil }

func setupTestService() (*LoanService, *mockEmailSender) {
	repo := NewInMemoryLoanRepository()
	email := &mockEmailSender{}
	svc := NewLoanService(repo, email)
	return svc, email
}

func TestCreateLoan(t *testing.T) {
	svc, _ := setupTestService()

	tests := []struct {
		name        string
		borrowerID  string
		principal   float64
		expectError bool
	}{
		{"Valid loan", "B001", 5000000, false},
		{"Zero principal", "B002", 0, false}, // Allowed in current logic, may change if validated
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ln, err := svc.CreateLoan(tt.borrowerID, tt.principal, 10, 12)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.borrowerID, ln.BorrowerID)
			}
		})
	}
}

func TestApproveLoan(t *testing.T) {
	svc, _ := setupTestService()

	tests := []struct {
		name       string
		approval   Approval
		shouldFail bool
	}{
		{
			"Valid approval",
			Approval{
				PhotoProofURL: "url",
				ValidatorID:   "EMP123",
				ApprovalDate:  time.Now(),
			}, false,
		},
		{
			"Missing fields",
			Approval{
				PhotoProofURL: "",
				ValidatorID:   "",
			}, true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ln, _ := svc.CreateLoan("B003", 4000000, 10, 12)
			_, err := svc.ApproveLoan(ln.ID, tt.approval)
			if tt.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInvestLoan(t *testing.T) {
	svc, email := setupTestService()

	t.Run("Fully funded triggers notification", func(t *testing.T) {
		ln, _ := svc.CreateLoan("B004", 1000000, 10, 10)
		svc.ApproveLoan(ln.ID, Approval{
			PhotoProofURL: "proof",
			ValidatorID:   "EMP001",
			ApprovalDate:  time.Now(),
		})

		_, err := svc.InvestLoan(ln.ID, Investor{ID: "INV001", Amount: 1000000})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(email.sentTo))
	})

	t.Run("Overfund should fail", func(t *testing.T) {
		ln, _ := svc.CreateLoan("B005", 2000000, 10, 10)
		svc.ApproveLoan(ln.ID, Approval{
			PhotoProofURL: "proof",
			ValidatorID:   "EMP002",
			ApprovalDate:  time.Now(),
		})

		_, err := svc.InvestLoan(ln.ID, Investor{ID: "INV999", Amount: 2500000})
		assert.Error(t, err)
	})
}

func TestDisburseLoan(t *testing.T) {
	svc, _ := setupTestService()

	ln, _ := svc.CreateLoan("B006", 1500000, 10, 10)
	svc.ApproveLoan(ln.ID, Approval{
		PhotoProofURL: "proof",
		ValidatorID:   "EMP777",
		ApprovalDate:  time.Now(),
	})
	svc.InvestLoan(ln.ID, Investor{ID: "INV", Amount: 1500000})

	tests := []struct {
		name         string
		disb         Disbursement
		agreementURL string
		shouldFail   bool
	}{
		{
			"Valid disbursement",
			Disbursement{
				AgreementFile:    "signed.jpg",
				FieldOfficerID:   "FO001",
				DisbursementDate: time.Now(),
			},
			"https://link.pdf", false,
		},
		{
			"Missing fields",
			Disbursement{
				AgreementFile:  "",
				FieldOfficerID: "",
			},
			"", true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.DisburseLoan(ln.ID, tt.disb, tt.agreementURL)
			if tt.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoanService_UpdateFailures(t *testing.T) {
	email := &mockEmailSender{}
	svc := NewLoanService(&errorRepo{}, email)

	tests := []struct {
		name   string
		action func() error
	}{
		{
			"ApproveLoan update failure",
			func() error {
				approval := Approval{
					PhotoProofURL: "https://proof.img",
					ValidatorID:   "EMP001",
					ApprovalDate:  time.Now(),
				}
				_, err := svc.ApproveLoan("LOAN001", approval)
				return err
			},
		},
		{
			"InvestLoan update failure",
			func() error {
				inv := Investor{ID: "INV01", Amount: 1000}
				_, err := svc.InvestLoan("LOAN002", inv)
				return err
			},
		},
		{
			"DisburseLoan update failure",
			func() error {
				d := Disbursement{
					AgreementFile:    "img.jpg",
					FieldOfficerID:   "FO123",
					DisbursementDate: time.Now(),
				}
				_, err := svc.DisburseLoan("LOAN003", d, "https://link.pdf")
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.action()
			assert.Error(t, err)
		})
	}
}
