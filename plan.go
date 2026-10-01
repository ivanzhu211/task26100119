package debtplan

import (
	"errors"
	"fmt"
)

// ErrInvalidInput is returned when the account list or horizon is invalid.
var ErrInvalidInput = errors.New("debtplan: invalid input")

// MinBudget finds the smallest fixed monthly budget, in integer cents,
// that pays off every account within months months, and returns the
// amortization schedule under that budget.
//
// Feasibility is not monotonic in the budget: the snowball payment order
// changes with remaining balances, so a larger budget can fail where a
// smaller one succeeds. The search therefore scans upward from a lower
// bound and returns the first feasible budget, which is the true minimum.
func MinBudget(accounts []Account, months int) (Plan, error) {
	if err := validate(accounts, months); err != nil {
		return Plan{}, err
	}

	// Lower bound: the budget must at least cover the first month's
	// minimum payments. Upper bound: enough to clear every account
	// (including first-month interest) in a single month, which is
	// always feasible.
	var lo, hi int64
	for _, a := range accounts {
		debt := a.Principal + monthlyInterest(a.Principal, a.AnnualRateBps)
		lo += minInt64(a.MinPayment, debt)
		hi += debt
	}

	for budget := lo; budget <= hi; budget++ {
		if plan, ok := Simulate(accounts, budget, months); ok {
			return plan, nil
		}
	}
	// Unreachable: the upper bound is always feasible.
	return Plan{}, fmt.Errorf("%w: no feasible budget found", ErrInvalidInput)
}

func validate(accounts []Account, months int) error {
	if months < 1 {
		return fmt.Errorf("%w: months must be >= 1, got %d", ErrInvalidInput, months)
	}
	seen := make(map[string]bool, len(accounts))
	for _, a := range accounts {
		switch {
		case a.ID == "":
			return fmt.Errorf("%w: account ID must not be empty", ErrInvalidInput)
		case seen[a.ID]:
			return fmt.Errorf("%w: duplicate account ID %q", ErrInvalidInput, a.ID)
		case a.Principal < 0:
			return fmt.Errorf("%w: account %q has negative principal", ErrInvalidInput, a.ID)
		case a.AnnualRateBps < 0:
			return fmt.Errorf("%w: account %q has negative rate", ErrInvalidInput, a.ID)
		case a.MinPayment <= 0:
			return fmt.Errorf("%w: account %q has non-positive minimum payment", ErrInvalidInput, a.ID)
		}
		seen[a.ID] = true
	}
	return nil
}
