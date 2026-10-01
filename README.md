# debtplan

Go 库：按雪球法（先清小额账户）反求在目标月数内清偿全部债务的**最低固定月预算**，并返回逐月摊销表。账户数据与计算结果仅保存在内存。

## 规则

- 金额以整数分表示，年利率以整数基点表示（1 基点 = 0.01%）。
- 每期先按期初余额计息：`余额 × 基点 / 120000`，每账户独立四舍五入到分，半分向上。
- 每期先从预算中满足所有未清账户的最低月还款（不超过含息欠款），预算不足则不可行。
- 剩余预算按剩余余额从小到大（同额按编号字典序）依次清偿，当期可连续清零多个账户。
- 预算每月固定、不结转；目标月数内（含首期）全部清零即达标，提前清偿也算。

由于付款顺序随余额变化，可行性对预算**不是单调**的（更大预算可能失败）。`MinBudget` 从下界向上扫描，返回第一个可行预算，即真正的最低值。

## 用法

```go
accounts := []debtplan.Account{
    {ID: "A", Principal: 15500, AnnualRateBps: 0,     MinPayment: 100},
    {ID: "B", Principal: 25620, AnnualRateBps: 12000, MinPayment: 10000},
}
plan, err := debtplan.MinBudget(accounts, 3)
// plan.Budget == 15326（即 153.26 元），plan.Months 为逐月摊销表
```

- `MinBudget(accounts, months) (Plan, error)`：求最低可行预算及摊销表。
- `Simulate(accounts, budget, months) (Plan, bool)`：模拟指定预算，返回摊销表与是否达标。

## 测试

```sh
go test ./...
```
