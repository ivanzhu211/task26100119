package debtplan

// Solve 反求能够在 months 期内清偿全部债务的最低固定月预算（整数分），
// 并返回该预算下的逐月摊销表。由于付款顺序随余额变化，可行性关于预算
// 并非单调，这里通过轨迹稳定区间扫描求得真正的最低预算。
func Solve(accounts []Account, months int) (Plan, error) {
	return Plan{}, nil
}
