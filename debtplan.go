// Package debtplan 计算在目标月数内清偿全部债务的最低固定月预算，
// 并给出该预算下的逐月摊销表。账户数据与计算结果仅保存在内存中。
//
// 金额以整数分表示，年利率以整数基点表示（1 基点 = 0.01%）。
// 付款采用雪球法：每期先满足所有未清账户的最低月还款，
// 剩余预算按剩余余额从小到大（同余额按编号字典序）依次清偿。
package debtplan

import "errors"

// Account 描述一个债务账户。
type Account struct {
	ID            string // 唯一编号
	Principal     int64  // 初始本金，单位：分，>= 0
	AnnualRateBps int64  // 固定年利率，单位：基点，>= 0
	MinPayment    int64  // 固定最低月还款，单位：分，> 0
}

// AccountMonth 记录单个账户在某一期内的明细。
type AccountMonth struct {
	Interest int64 // 本期期初计提的利息，单位：分
	Payment  int64 // 本期实际还款总额（最低还款 + 雪球还款），单位：分
	Balance  int64 // 本期期末余额，单位：分
}

// Month 记录某一期内所有账户的明细，下标与输入账户顺序一致。
type Month struct {
	Accounts []AccountMonth
}

// Plan 是求解结果：最低月预算及该预算下的逐月摊销表。
// 若提前清偿，Months 只包含到最后一个实际还款期为止。
type Plan struct {
	Budget int64
	Months []Month
}

// 输入校验错误。
var (
	ErrNoMonths      = errors.New("debtplan: months must be >= 1")
	ErrEmptyID       = errors.New("debtplan: account ID must not be empty")
	ErrDuplicateID   = errors.New("debtplan: duplicate account ID")
	ErrNegativeValue = errors.New("debtplan: principal and annual rate must be >= 0")
	ErrMinPayment    = errors.New("debtplan: min payment must be > 0")
	ErrInfeasible    = errors.New("debtplan: no feasible budget exists")
)

// validate 校验输入参数。
func validate(accounts []Account, months int) error {
	if months < 1 {
		return ErrNoMonths
	}
	seen := make(map[string]struct{}, len(accounts))
	for _, a := range accounts {
		if a.ID == "" {
			return ErrEmptyID
		}
		if _, ok := seen[a.ID]; ok {
			return ErrDuplicateID
		}
		seen[a.ID] = struct{}{}
		if a.Principal < 0 || a.AnnualRateBps < 0 {
			return ErrNegativeValue
		}
		if a.MinPayment <= 0 {
			return ErrMinPayment
		}
	}
	return nil
}
