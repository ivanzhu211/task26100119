package debtplan

// affine 表示关于预算增量 δ 的仿射量：实际值 = val + slope·δ。
// 在轨迹（各期事件序列）保持不变的预算区间内，所有余额与还款额
// 都是预算的仿射函数。
type affine struct {
	val   int64
	slope int64
}

// runResult 是一次仿射模拟的结果。
type runResult struct {
	feasible  bool  // 该预算是否可行
	deltaStar int64 // 轨迹发生变化的最小预算增量（>= 1）
	months    []Month
}

// roundInterest 计算一期利息：余额乘年利率基点数再除以 120000，
// 四舍五入到分，恰好半分时向上取整。
func roundInterest(balance, annualRateBps int64) int64 {
	return 0
}

// run 以预算 budget 进行模拟，返回可行性、轨迹稳定区间长度及摊销表。
func run(accounts []Account, months int, budget int64) runResult {
	return runResult{}
}

// Simulate 以固定月预算 budget 模拟 months 期还款过程，
// 返回逐月摊销表及该预算是否能在期限内清偿全部债务。
func Simulate(accounts []Account, months int, budget int64) (Plan, bool) {
	return Plan{}, false
}
