// Package debtplan computes the minimum fixed monthly budget that pays off
// a set of debt accounts within a target number of months, using the
// snowball method (smallest remaining balance first).
//
// All amounts are integer cents. Annual rates are integer basis points
// (1 bp = 0.01%). Everything lives in memory; nothing is persisted.
package debtplan

// Account describes a single debt account.
type Account struct {
	ID            string // unique account identifier
	Principal     int64  // initial balance, in cents
	AnnualRateBps int64  // fixed annual interest rate, in basis points
	MinPayment    int64  // fixed minimum monthly payment, in cents (must be > 0)
}

// AccountMonth records one account's activity within a single month.
type AccountMonth struct {
	ID           string
	BeginBalance int64 // balance at the start of the month, before interest
	Interest     int64 // interest charged this month
	MinPaid      int64 // minimum payment applied
	ExtraPaid    int64 // snowball (extra) payment applied
	EndBalance   int64 // balance after all payments this month
}

// Month is one row of the amortization schedule.
type Month struct {
	Index    int            // 1-based month number
	Accounts []AccountMonth // per-account activity, in input account order
}

// Plan is the result of budgeting: a fixed monthly budget and the
// amortization schedule produced under that budget.
type Plan struct {
	Budget int64   // fixed monthly budget, in cents
	Months []Month // schedule stops at the month everything is paid off
}
