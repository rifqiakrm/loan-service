package loan

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// LoanRepository defines the contract for any loan storage mechanism.
type LoanRepository interface {
	Create(loan *Loan) error
	GetByID(id string) (*Loan, error)
	Update(loan *Loan) error
	List() ([]*Loan, error)
}

// InMemoryLoanRepository provides a thread-safe in-memory store for loans.
// It is useful for development, testing, or as a temporary mock.
type InMemoryLoanRepository struct {
	store sync.Map
}

// NewInMemoryLoanRepository creates and returns a new in-memory loan repository instance.
func NewInMemoryLoanRepository() *InMemoryLoanRepository {
	return &InMemoryLoanRepository{}
}

// Create inserts a new loan into the store and assigns it a unique ID.
func (r *InMemoryLoanRepository) Create(loan *Loan) error {
	loan.ID = uuid.NewString()
	now := time.Now()
	loan.CreatedAt = now
	loan.UpdatedAt = now
	loan.State = Proposed
	r.store.Store(loan.ID, loan)
	return nil
}

// GetByID retrieves a loan by its ID. Returns error if not found.
func (r *InMemoryLoanRepository) GetByID(id string) (*Loan, error) {
	if val, ok := r.store.Load(id); ok {
		if loan, valid := val.(*Loan); valid {
			return loan, nil
		}
	}
	return nil, errors.New("loan not found")
}

// Update updates an existing loan in the store.
func (r *InMemoryLoanRepository) Update(loan *Loan) error {
	if _, ok := r.store.Load(loan.ID); !ok {
		return errors.New("loan not found for update")
	}
	loan.UpdatedAt = time.Now()
	r.store.Store(loan.ID, loan)
	return nil
}

// List returns all loans in the store.
func (r *InMemoryLoanRepository) List() ([]*Loan, error) {
	var result []*Loan
	r.store.Range(func(_, val any) bool {
		if loan, ok := val.(*Loan); ok {
			result = append(result, loan)
		}
		return true
	})
	return result, nil
}
