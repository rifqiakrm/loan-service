package loan

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoanStateTransitions(t *testing.T) {
	tests := []struct {
		name  string
		from  LoanState
		to    LoanState
		valid bool
	}{
		{"Proposed to Approved", Proposed, Approved, true},
		{"Approved to Invested", Approved, Invested, true},
		{"Invested to Disbursed", Invested, Disbursed, true},
		{"Proposed to Disbursed", Proposed, Disbursed, false},
		{"Disbursed to Proposed", Disbursed, Proposed, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTransition(tt.from, tt.to)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
