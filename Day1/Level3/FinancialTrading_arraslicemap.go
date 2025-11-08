package main

import "fmt"

/*

🟢 Industry Insight

Arrays = fixed instruments, towers, servers

Slices = dynamic data like active users, orders, connections

Maps = fast lookups (tower → signal, stock → price, user → session)

*/

func main() {
	// Fixed trading instruments array
	instruments := [3]string{"INFY", "TCS", "RELIANCE"}
	fmt.Println("[TRADE] Instruments Array:", instruments)

	// Slice — dynamic orders placed
	orders := []string{"Buy INFY", "Sell TCS"}
	orders = append(orders, "Buy RELIANCE")
	fmt.Println("[TRADE] Orders Slice:", orders)

	// Map — symbol → current price
	priceMap := map[string]float64{
		"INFY": 1850.50,
		"TCS":  3600.75,
	}
	priceMap["RELIANCE"] = 2450.80
	fmt.Println("[TRADE] Price Map:", priceMap)
}
