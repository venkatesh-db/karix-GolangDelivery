package main

import "log"

/*

🟢 Financial Context

OrderState enum improves safety — no string mismatches.

Constants define system-wide policies (risk %, max orders).

Type alias INR clarifies money handling for precision.

*/

// Order state enum
type OrderState int

const (
    OrderPending OrderState = iota
    OrderExecuted
    OrderCancelled
)

// Trade limits
const (
    MaxDailyOrders int     = 5000
    RiskLimitPct   float64 = 2.5
)

// Custom type alias for prices
type INR float64

func main() {
    var (
        orderID   string     = "ORDX98765"
        state     OrderState = OrderExecuted
        tradePrice INR       = 1856.40
    )

    log.Printf("[TRADE] OrderID=%s | State=%d | Price=%.2f | MaxOrders=%d | RiskLimit=%.2f%%",
        orderID, state, tradePrice, MaxDailyOrders, RiskLimitPct)
}

