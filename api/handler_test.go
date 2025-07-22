// api/handler_test.go (Full Table-Driven Test Suite)
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"loan-service/core/loan"
)

type mockEmailSender struct {
	sentTo []string
}

func (m *mockEmailSender) SendInvestorNotification(investorID, agreementLink string) error {
	m.sentTo = append(m.sentTo, investorID)
	return nil
}

func setupRouterWithMemoryService() (*gin.Engine, *loan.LoanService) {
	repo := loan.NewInMemoryLoanRepository()
	email := &mockEmailSender{}
	svc := loan.NewLoanService(repo, email)
	handler := NewHandler(svc)
	return SetupRouter(handler), svc
}

func TestLoanHandlers_All(t *testing.T) {
	router, svc := setupRouterWithMemoryService()
	type testCase struct {
		name       string
		method     string
		endpoint   string
		payload    interface{}
		setup      func() string
		expectCode int
		contains   string
	}

	now := time.Now().Format("2006-01-02")

	tests := []testCase{
		{
			name:     "CreateLoan success",
			method:   "POST",
			endpoint: "/loans",
			payload: map[string]interface{}{
				"borrower_id":      "B001",
				"principal_amount": 5000000,
				"rate":             10,
				"roi":              12,
			},
			expectCode: 201,
			contains:   "\"borrower_id\":\"B001\"",
		},
		{
			name:       "CreateLoan invalid JSON",
			method:     "POST",
			endpoint:   "/loans",
			payload:    `invalid json`,
			expectCode: 400,
		},
		{
			name:   "ApproveLoan success",
			method: "POST",
			setup: func() string {
				ln, _ := svc.CreateLoan("B002", 4000000, 10, 10)
				return "/loans/" + ln.ID + "/approve"
			},
			payload: map[string]interface{}{
				"photo_proof_url":    "https://proof",
				"field_validator_id": "EMP001",
				"approval_date":      time.Now().Format("2006-01-02"),
			},
			expectCode: 200,
			contains:   "\"state\":\"approved\"",
		},
		{
			name:   "ApproveLoan missing fields",
			method: "POST",
			setup: func() string {
				ln, _ := svc.CreateLoan("B003", 3000, 10, 10)
				return "/loans/" + ln.ID + "/approve"
			},
			payload: map[string]interface{}{
				"photo_proof_url":    "",
				"field_validator_id": "",
				"approval_date":      "2025-07-22",
			},
			expectCode: 400,
		},
		{
			name:   "ApproveLoan invalid date format",
			method: "POST",
			setup: func() string {
				ln, _ := svc.CreateLoan("B007", 1000, 1, 1)
				return "/loans/" + ln.ID + "/approve"
			},
			payload: map[string]interface{}{
				"photo_proof_url":    "img",
				"field_validator_id": "EMP007",
				"approval_date":      "07/22/2025", // wrong format
			},
			expectCode: 400,
		},
		{
			name:   "DisburseLoan invalid date format",
			method: "POST",
			setup: func() string {
				ln, _ := svc.CreateLoan("B008", 2000, 10, 10)
				svc.ApproveLoan(ln.ID, loan.Approval{
					PhotoProofURL: "url", ValidatorID: "EMP008", ApprovalDate: time.Now(),
				})
				svc.InvestLoan(ln.ID, loan.Investor{ID: "INV008", Amount: 2000})
				return "/loans/" + ln.ID + "/disburse"
			},
			payload: map[string]interface{}{
				"agreement_letter_file": "x.jpg",
				"field_officer_id":      "FO008",
				"disbursement_date":     "not-a-date",
				"agreement_letter_link": "https://fail.com",
			},
			expectCode: 400,
		},
		{
			name:       "CreateLoan bad request (invalid json)",
			method:     "POST",
			endpoint:   "/loans",
			payload:    `{"borrower_id": B001,}`,
			expectCode: 400,
		},
		{
			name:   "ApproveLoan bad request (missing fields)",
			method: "POST",
			setup: func() string {
				ln, _ := svc.CreateLoan("B111", 1000, 10, 10)
				return "/loans/" + ln.ID + "/approve"
			},
			payload: map[string]interface{}{
				"photo_proof_url":    "",
				"field_validator_id": "",
				"approval_date":      now,
			},
			expectCode: 400,
		},
		{
			name:   "InvestLoan malformed JSON",
			method: "POST",
			setup: func() string {
				ln, _ := svc.CreateLoan("B009", 1000, 1, 1)
				svc.ApproveLoan(ln.ID, loan.Approval{
					PhotoProofURL: "proof", ValidatorID: "EMP009", ApprovalDate: time.Now(),
				})
				return "/loans/" + ln.ID + "/invest"
			},
			payload:    `{"investor_id": "INV999",`,
			expectCode: 400,
		},
		{
			name:   "DisburseLoan bad request (missing fields)",
			method: "POST",
			setup: func() string {
				ln, _ := svc.CreateLoan("B333", 3000, 10, 10)
				svc.ApproveLoan(ln.ID, loan.Approval{
					PhotoProofURL: "img", ValidatorID: "VAL2", ApprovalDate: time.Now(),
				})
				svc.InvestLoan(ln.ID, loan.Investor{ID: "INV3", Amount: 3000})
				return "/loans/" + ln.ID + "/disburse"
			},
			payload: map[string]interface{}{
				"agreement_letter_file": "",
				"field_officer_id":      "",
				"disbursement_date":     now,
				"agreement_letter_link": "",
			},
			expectCode: 400,
		},
		{
			name:   "InvestLoan success",
			method: "POST",
			setup: func() string {
				ln, _ := svc.CreateLoan("B004", 3000000, 10, 10)
				svc.ApproveLoan(ln.ID, loan.Approval{
					PhotoProofURL: "proof", ValidatorID: "EMPX", ApprovalDate: time.Now(),
				})
				return "/loans/" + ln.ID + "/invest"
			},
			payload: map[string]interface{}{
				"investor_id": "INV123",
				"amount":      3000000,
			},
			expectCode: 200,
			contains:   "\"state\":\"invested\"",
		},
		{
			name:   "DisburseLoan success",
			method: "POST",
			setup: func() string {
				ln, _ := svc.CreateLoan("B005", 2000000, 10, 10)
				svc.ApproveLoan(ln.ID, loan.Approval{
					PhotoProofURL: "proof", ValidatorID: "EMPY", ApprovalDate: time.Now(),
				})
				svc.InvestLoan(ln.ID, loan.Investor{ID: "INV123", Amount: 2000000})
				return "/loans/" + ln.ID + "/disburse"
			},
			payload: map[string]interface{}{
				"agreement_letter_file": "signed.jpg",
				"field_officer_id":      "FO001",
				"disbursement_date":     time.Now().Format("2006-01-02"),
				"agreement_letter_link": "https://link.com",
			},
			expectCode: 200,
			contains:   "\"state\":\"disbursed\"",
		},
		{
			name:     "GetLoan success",
			method:   "GET",
			endpoint: "/loans/",
			setup: func() string {
				ln, _ := svc.CreateLoan("B000", 1000000, 10, 10)
				return "/loans/" + ln.ID
			},
			expectCode: 200,
		},
		{
			name:       "GetLoan not found",
			method:     "GET",
			endpoint:   "/loans/non-existent-id",
			expectCode: 404,
		},
		{
			name:   "ListLoans success",
			method: "GET",
			setup: func() string {
				svc.CreateLoan("B006", 10000, 10, 10)
				return "/loans"
			},
			expectCode: 200,
			contains:   "\"borrower_id\":\"B006\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := tt.endpoint
			if tt.setup != nil {
				url = tt.setup()
			}
			var req *http.Request
			var err error
			switch payload := tt.payload.(type) {
			case nil:
				req, err = http.NewRequest(tt.method, url, nil)
			case string:
				req, err = http.NewRequest(tt.method, url, bytes.NewBufferString(payload))
			default:
				b, _ := json.Marshal(payload)
				req, err = http.NewRequest(tt.method, url, bytes.NewBuffer(b))
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectCode, w.Code)
			if tt.contains != "" {
				assert.Contains(t, w.Body.String(), tt.contains)
			}
		})
	}
}

type brokenRepoList struct{}

func (r *brokenRepoList) Create(*loan.Loan) error            { return nil }
func (r *brokenRepoList) GetByID(string) (*loan.Loan, error) { return nil, nil }
func (r *brokenRepoList) Update(*loan.Loan) error            { return nil }
func (r *brokenRepoList) List() ([]*loan.Loan, error)        { return nil, errors.New("fail list") }

func TestListLoansInternalError(t *testing.T) {
	email := &mockEmailSender{}
	svc := loan.NewLoanService(&brokenRepoList{}, email)
	router := SetupRouter(NewHandler(svc))

	req, _ := http.NewRequest("GET", "/loans", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
}
