package loan

import (
	"errors"
)

// EmailSender defines the interface for sending email notifications.
// You can implement this using SMTP, external APIs, or mock logs.
type EmailSender interface {
	SendInvestorNotification(investorID, agreementLink string) error
}

// LoanService provides core logic for managing loan lifecycle operations.
type LoanService struct {
	repo  LoanRepository
	email EmailSender
}

// NewLoanService creates a new instance of LoanService.
func NewLoanService(repo LoanRepository, email EmailSender) *LoanService {
	return &LoanService{
		repo:  repo,
		email: email,
	}
}

// CreateLoan creates a new loan with the given parameters.
func (s *LoanService) CreateLoan(borrowerID string, principal float64, rate float64, roi float64) (*Loan, error) {
	loan := &Loan{
		BorrowerID:      borrowerID,
		PrincipalAmount: principal,
		Rate:            rate,
		ROI:             roi,
	}
	if err := s.repo.Create(loan); err != nil {
		return nil, err
	}
	return loan, nil
}

// ApproveLoan moves a loan to Approved state after validating the input data.
func (s *LoanService) ApproveLoan(loanID string, approval Approval) (*Loan, error) {
	loan, err := s.repo.GetByID(loanID)
	if err != nil {
		return nil, err
	}

	if err := ValidateTransition(loan.State, Approved); err != nil {
		return nil, err
	}

	// Ensure required fields are set
	if approval.PhotoProofURL == "" || approval.ValidatorID == "" || approval.ApprovalDate.IsZero() {
		return nil, errors.New("missing approval fields")
	}

	loan.State = Approved
	loan.Approval = &approval

	return s.updateLoan(loan)
}

// InvestLoan adds a new investor to a loan. If fully funded, it moves to Invested state and sends notifications.
func (s *LoanService) InvestLoan(loanID string, investor Investor) (*Loan, error) {
	loan, err := s.repo.GetByID(loanID)
	if err != nil {
		return nil, err
	}

	if loan.State != Approved && loan.State != Invested {
		return nil, errors.New("loan must be in approved or invested state to accept investments")
	}

	// Check if adding this investment exceeds principal
	if loan.TotalInvested+investor.Amount > loan.PrincipalAmount {
		return nil, errors.New("investment exceeds loan principal")
	}

	// Add investor
	loan.Investors = append(loan.Investors, investor)
	loan.TotalInvested += investor.Amount

	// Move to Invested if fully funded
	if loan.TotalInvested == loan.PrincipalAmount {
		if err := ValidateTransition(loan.State, Invested); err != nil {
			return nil, err
		}
		loan.State = Invested

		// Notify all investors
		for _, inv := range loan.Investors {
			_ = s.email.SendInvestorNotification(inv.ID, loan.AgreementLetterURL)
		}
	}

	return s.updateLoan(loan)
}

// DisburseLoan moves a loan to Disbursed state and stores agreement and field officer info.
func (s *LoanService) DisburseLoan(loanID string, disb Disbursement, agreementLink string) (*Loan, error) {
	loan, err := s.repo.GetByID(loanID)
	if err != nil {
		return nil, err
	}

	if err := ValidateTransition(loan.State, Disbursed); err != nil {
		return nil, err
	}

	if disb.AgreementFile == "" || disb.FieldOfficerID == "" || disb.DisbursementDate.IsZero() {
		return nil, errors.New("missing disbursement fields")
	}

	loan.Disbursement = &disb
	loan.AgreementLetterURL = agreementLink
	loan.State = Disbursed

	return s.updateLoan(loan)
}

// GetLoan retrieves a loan by its ID.
func (s *LoanService) GetLoan(id string) (*Loan, error) {
	return s.repo.GetByID(id)
}

// ListLoans returns all loans in the system.
func (s *LoanService) ListLoans() ([]*Loan, error) {
	return s.repo.List()
}

// updateLoan updates the loan and saves it via repository.
func (s *LoanService) updateLoan(loan *Loan) (*Loan, error) {
	if err := s.repo.Update(loan); err != nil {
		return nil, err
	}
	return loan, nil
}
