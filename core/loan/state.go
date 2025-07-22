package loan

import "fmt"

// CanTransition determines if a loan can move from its current state to the desired next state.
//
// Valid transitions:
//   - Proposed  → Approved
//   - Approved  → Invested
//   - Invested  → Disbursed
//
// Backward or invalid transitions are not allowed.
func CanTransition(from, to LoanState) bool {
	switch from {
	case Proposed:
		return to == Approved
	case Approved:
		return to == Invested
	case Invested:
		return to == Disbursed
	default:
		return false
	}
}

// ValidateTransition checks whether a transition from `current` to `next` state is valid.
// Returns an error if the transition is not allowed.
func ValidateTransition(current, next LoanState) error {
	if !CanTransition(current, next) {
		return fmt.Errorf("invalid state transition: cannot move from %s to %s", current, next)
	}
	return nil
}
