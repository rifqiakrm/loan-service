package loan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryLoanRepository(t *testing.T) {
	repo := NewInMemoryLoanRepository()

	t.Run("Create and GetByID", func(t *testing.T) {
		ln := &Loan{BorrowerID: "B001", PrincipalAmount: 12345}
		err := repo.Create(ln)
		assert.NoError(t, err)
		assert.NotEmpty(t, ln.ID)

		fetched, err := repo.GetByID(ln.ID)
		assert.NoError(t, err)
		assert.Equal(t, ln.ID, fetched.ID)
	})

	t.Run("GetByID non-existent", func(t *testing.T) {
		_, err := repo.GetByID("does-not-exist")
		assert.Error(t, err)
	})

	t.Run("Update existing loan", func(t *testing.T) {
		ln := &Loan{BorrowerID: "B002", PrincipalAmount: 1000}
		_ = repo.Create(ln)
		ln.Rate = 99
		err := repo.Update(ln)
		assert.NoError(t, err)

		updated, _ := repo.GetByID(ln.ID)
		assert.Equal(t, 99.0, updated.Rate)
	})

	t.Run("Update non-existent loan", func(t *testing.T) {
		err := repo.Update(&Loan{ID: "fake-id"})
		assert.Error(t, err)
	})

	t.Run("List all loans", func(t *testing.T) {
		list, err := repo.List()
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(list), 1)
	})
}
