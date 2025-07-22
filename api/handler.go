package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"loan-service/core/loan"
)

// Handler contains dependencies needed by the HTTP routes.
type Handler struct {
	Service *loan.LoanService
}

// NewHandler creates a new HTTP handler instance.
func NewHandler(service *loan.LoanService) *Handler {
	return &Handler{Service: service}
}

// CreateLoan handles POST /loans to create a new loan.
func (h *Handler) CreateLoan(c *gin.Context) {
	var req struct {
		BorrowerID      string  `json:"borrower_id" binding:"required"`
		PrincipalAmount float64 `json:"principal_amount" binding:"required"`
		Rate            float64 `json:"rate" binding:"required"`
		ROI             float64 `json:"roi" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	loan, err := h.Service.CreateLoan(req.BorrowerID, req.PrincipalAmount, req.Rate, req.ROI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, loan)
}

// ApproveLoan handles POST /loans/:id/approve
func (h *Handler) ApproveLoan(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		PhotoProofURL string `json:"photo_proof_url" binding:"required"`
		ValidatorID   string `json:"field_validator_id" binding:"required"`
		ApprovalDate  string `json:"approval_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	date, err := time.Parse("2006-01-02", req.ApprovalDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (expected YYYY-MM-DD)"})
		return
	}

	approval := loan.Approval{
		PhotoProofURL: req.PhotoProofURL,
		ValidatorID:   req.ValidatorID,
		ApprovalDate:  date,
	}

	ln, err := h.Service.ApproveLoan(id, approval)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ln)
}

// InvestLoan handles POST /loans/:id/invest
func (h *Handler) InvestLoan(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		InvestorID string  `json:"investor_id" binding:"required"`
		Amount     float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	investor := loan.Investor{
		ID:     req.InvestorID,
		Amount: req.Amount,
	}

	ln, err := h.Service.InvestLoan(id, investor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ln)
}

// DisburseLoan handles POST /loans/:id/disburse
func (h *Handler) DisburseLoan(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		AgreementFile    string `json:"agreement_letter_file" binding:"required"`
		FieldOfficerID   string `json:"field_officer_id" binding:"required"`
		DisbursementDate string `json:"disbursement_date" binding:"required"`
		AgreementLink    string `json:"agreement_letter_link" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	date, err := time.Parse("2006-01-02", req.DisbursementDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (expected YYYY-MM-DD)"})
		return
	}

	disb := loan.Disbursement{
		AgreementFile:    req.AgreementFile,
		FieldOfficerID:   req.FieldOfficerID,
		DisbursementDate: date,
	}

	ln, err := h.Service.DisburseLoan(id, disb, req.AgreementLink)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ln)
}

// GetLoan handles GET /loans/:id
func (h *Handler) GetLoan(c *gin.Context) {
	id := c.Param("id")
	ln, err := h.Service.GetLoan(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}
	c.JSON(http.StatusOK, ln)
}

// ListLoans handles GET /loans
func (h *Handler) ListLoans(c *gin.Context) {
	list, err := h.Service.ListLoans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list loans"})
		return
	}
	c.JSON(http.StatusOK, list)
}
