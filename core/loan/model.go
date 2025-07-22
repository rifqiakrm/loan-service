package loan

import "time"

// LoanState represents the lifecycle state of a loan.
type LoanState string

const (
	// Proposed is the initial state when a loan is created.
	Proposed LoanState = "proposed"

	// Approved is the state after a loan has been approved by staff.
	Approved LoanState = "approved"

	// Invested is the state after the loan has been fully funded by investors.
	Invested LoanState = "invested"

	// Disbursed is the state after the loan is handed over to the borrower.
	Disbursed LoanState = "disbursed"
)

// Loan represents a loan given to a borrower, along with its current state and data.
type Loan struct {
	ID                 string        `json:"id"`                     // Unique identifier of the loan
	BorrowerID         string        `json:"borrower_id"`            // Identifier of the borrower
	PrincipalAmount    float64       `json:"principal_amount"`       // Total loan principal amount
	Rate               float64       `json:"rate"`                   // Interest rate the borrower must pay (in %)
	ROI                float64       `json:"roi"`                    // Return of investment for investors (in %)
	AgreementLetterURL string        `json:"agreement_letter_link"`  // URL to agreement letter (if generated)
	State              LoanState     `json:"state"`                  // Current lifecycle state of the loan
	Approval           *Approval     `json:"approval,omitempty"`     // Approval information (if approved)
	Disbursement       *Disbursement `json:"disbursement,omitempty"` // Disbursement information (if disbursed)
	Investors          []Investor    `json:"investors"`              // List of investors
	TotalInvested      float64       `json:"total_invested"`         // Total amount invested by all investors
	CreatedAt          time.Time     `json:"created_at"`             // Timestamp when loan was created
	UpdatedAt          time.Time     `json:"updated_at"`             // Timestamp when loan was last updated
}

// Approval holds information regarding the loan approval by a field validator.
type Approval struct {
	PhotoProofURL string    `json:"photo_proof_url"`    // URL of photo proof taken by field validator
	ValidatorID   string    `json:"field_validator_id"` // Employee ID of the field validator
	ApprovalDate  time.Time `json:"approval_date"`      // Date of approval
}

// Disbursement contains details about loan disbursement to the borrower.
type Disbursement struct {
	AgreementFile    string    `json:"agreement_letter_file"` // File or image of signed agreement letter
	FieldOfficerID   string    `json:"field_officer_id"`      // Employee ID of field officer
	DisbursementDate time.Time `json:"disbursement_date"`     // Date of disbursement
}

// Investor represents a single investor and the amount they contributed to the loan.
type Investor struct {
	ID     string  `json:"investor_id"` // Unique identifier of the investor
	Amount float64 `json:"amount"`      // Amount invested
}
