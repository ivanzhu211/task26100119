package debtplan

import "sort"

// Simulate runs the snowball amortization for the given fixed monthly
// budget over at most months months. It returns the schedule and whether
// all accounts were fully paid off within the horizon.
//
// Each month, in order:
//  1. Interest accrues on every account's beginning-of-month balance and
//     is added to its debt.
//  2. Every unpaid account's minimum payment (capped at its debt) is
//     covered from the budget. If the budget cannot cover them all, the
//     budget is infeasible and simulation stops.
//  3. The remaining budget is paid to unpaid accounts in ascending
//     balance order (ties broken by ID); once an account reaches zero,
//     payment continues to the next account in the same month.
//
// The schedule ends early once every account is paid off.
func Simulate(accounts []Account, budget int64, months int) (Plan, bool) {
	balances := make([]int64, len(accounts))
	for i, a := range accounts {
		balances[i] = a.Principal
	}

	plan := Plan{Budget: budget}
	for m := 1; m <= months; m++ {
		month := Month{Index: m, Accounts: make([]AccountMonth, len(accounts))}

		// 1. Interest on the beginning-of-month balance.
		for i, a := range accounts {
			entry := AccountMonth{ID: a.ID, BeginBalance: balances[i]}
			if balances[i] > 0 {
				entry.Interest = monthlyInterest(balances[i], a.AnnualRateBps)
				balances[i] += entry.Interest
			}
			month.Accounts[i] = entry
		}

		// 2. Minimum payments; the budget must cover all of them.
		var minTotal int64
		for i, a := range accounts {
			if balances[i] > 0 {
				minTotal += minInt64(a.MinPayment, balances[i])
			}
		}
		if minTotal > budget {
			plan.Months = append(plan.Months, month)
			return plan, false
		}
		remaining := budget - minTotal
		for i, a := range accounts {
			if balances[i] > 0 {
				pay := minInt64(a.MinPayment, balances[i])
				balances[i] -= pay
				month.Accounts[i].MinPaid = pay
			}
		}

		// 3. Snowball: smallest remaining balance first, ties by ID.
		order := make([]int, 0, len(accounts))
		for i := range accounts {
			if balances[i] > 0 {
				order = append(order, i)
			}
		}
		sort.Slice(order, func(x, y int) bool {
			if balances[order[x]] != balances[order[y]] {
				return balances[order[x]] < balances[order[y]]
			}
			return accounts[order[x]].ID < accounts[order[y]].ID
		})
		for _, i := range order {
			if remaining == 0 {
				break
			}
			pay := minInt64(remaining, balances[i])
			balances[i] -= pay
			remaining -= pay
			month.Accounts[i].ExtraPaid = pay
		}

		done := true
		for i := range accounts {
			month.Accounts[i].EndBalance = balances[i]
			if balances[i] > 0 {
				done = false
			}
		}
		plan.Months = append(plan.Months, month)
		if done {
			return plan, true
		}
	}
	return plan, false
}

// monthlyInterest returns the interest for one month on balance cents at
// annualRateBps basis points: balance * bps / 120000, rounded to the
// nearest cent with halves rounded up.
func monthlyInterest(balance, annualRateBps int64) int64 {
	// round(x/d + 1/2) for non-negative values == (2x + d) / (2d).
	return (2*balance*annualRateBps + 120000) / 240000
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
