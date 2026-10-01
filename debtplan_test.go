package debtplan

import "testing"

// threeMonthAccounts returns the scenario from the spec: account A with
// 155.00 at 0% and a 1.00 minimum, account B with 256.20 at 120% and a
// 100.00 minimum.
func threeMonthAccounts() []Account {
	return []Account{
		{ID: "A", Principal: 15500, AnnualRateBps: 0, MinPayment: 100},
		{ID: "B", Principal: 25620, AnnualRateBps: 12000, MinPayment: 10000},
	}
}

func TestMinBudgetThreeMonths(t *testing.T) {
	plan, err := MinBudget(threeMonthAccounts(), 3)
	if err != nil {
		t.Fatalf("MinBudget: %v", err)
	}
	if plan.Budget != 15326 {
		t.Fatalf("budget = %d, want 15326", plan.Budget)
	}

	want := []Month{
		{Index: 1, Accounts: []AccountMonth{
			{ID: "A", BeginBalance: 15500, Interest: 0, MinPaid: 100, ExtraPaid: 5226, EndBalance: 10174},
			{ID: "B", BeginBalance: 25620, Interest: 2562, MinPaid: 10000, ExtraPaid: 0, EndBalance: 18182},
		}},
		{Index: 2, Accounts: []AccountMonth{
			{ID: "A", BeginBalance: 10174, Interest: 0, MinPaid: 100, ExtraPaid: 0, EndBalance: 10074},
			{ID: "B", BeginBalance: 18182, Interest: 1818, MinPaid: 10000, ExtraPaid: 5226, EndBalance: 4774},
		}},
		{Index: 3, Accounts: []AccountMonth{
			{ID: "A", BeginBalance: 10074, Interest: 0, MinPaid: 100, ExtraPaid: 9974, EndBalance: 0},
			{ID: "B", BeginBalance: 4774, Interest: 477, MinPaid: 5251, ExtraPaid: 0, EndBalance: 0},
		}},
	}

	if len(plan.Months) != len(want) {
		t.Fatalf("got %d months, want %d", len(plan.Months), len(want))
	}
	for m := range want {
		got, exp := plan.Months[m], want[m]
		if got.Index != exp.Index {
			t.Fatalf("month index = %d, want %d", got.Index, exp.Index)
		}
		for i := range exp.Accounts {
			if got.Accounts[i] != exp.Accounts[i] {
				t.Errorf("month %d account %s:\n got %+v\nwant %+v",
					got.Index, exp.Accounts[i].ID, got.Accounts[i], exp.Accounts[i])
			}
		}
	}
}

// Feasibility is not monotonic: 153.26 works, 154.00 fails, 155.00 works.
func TestFeasibilityNotMonotonic(t *testing.T) {
	accounts := threeMonthAccounts()
	if _, ok := Simulate(accounts, 15325, 3); ok {
		t.Error("budget 15325 should not pay off within 3 months")
	}
	if _, ok := Simulate(accounts, 15326, 3); !ok {
		t.Error("budget 15326 should pay off within 3 months")
	}
	if _, ok := Simulate(accounts, 15400, 3); ok {
		t.Error("budget 15400 should not pay off within 3 months")
	}
	if _, ok := Simulate(accounts, 15500, 3); !ok {
		t.Error("budget 15500 should pay off within 3 months")
	}
}

// When the final month's debt (after interest) is smaller than the minimum
// payment, only the actual debt is paid.
func TestFinalMonthPaysActualDebt(t *testing.T) {
	accounts := []Account{
		{ID: "A", Principal: 10000, AnnualRateBps: 0, MinPayment: 6000},
	}
	plan, err := MinBudget(accounts, 2)
	if err != nil {
		t.Fatalf("MinBudget: %v", err)
	}
	if plan.Budget != 6000 {
		t.Fatalf("budget = %d, want 6000", plan.Budget)
	}
	if len(plan.Months) != 2 {
		t.Fatalf("got %d months, want 2", len(plan.Months))
	}
	last := plan.Months[1].Accounts[0]
	if last.MinPaid != 4000 || last.EndBalance != 0 {
		t.Errorf("final month: got %+v, want MinPaid 4000 and EndBalance 0", last)
	}
}

func TestHalfCentRoundsUp(t *testing.T) {
	// 4775 cents at 12000 bps: 4775 * 12000 / 120000 = 477.5 -> 478.
	if got := monthlyInterest(4775, 12000); got != 478 {
		t.Errorf("monthlyInterest(4775, 12000) = %d, want 478", got)
	}
	if got := monthlyInterest(18182, 12000); got != 1818 {
		t.Errorf("monthlyInterest(18182, 12000) = %d, want 1818", got)
	}
}

func TestInvalidInput(t *testing.T) {
	valid := []Account{{ID: "A", Principal: 100, AnnualRateBps: 0, MinPayment: 10}}
	if _, err := MinBudget(valid, 0); err == nil {
		t.Error("months = 0 should fail")
	}
	dup := []Account{
		{ID: "A", Principal: 100, AnnualRateBps: 0, MinPayment: 10},
		{ID: "A", Principal: 100, AnnualRateBps: 0, MinPayment: 10},
	}
	if _, err := MinBudget(dup, 1); err == nil {
		t.Error("duplicate IDs should fail")
	}
	bad := []Account{{ID: "A", Principal: 100, AnnualRateBps: 0, MinPayment: 0}}
	if _, err := MinBudget(bad, 1); err == nil {
		t.Error("zero minimum payment should fail")
	}
}
