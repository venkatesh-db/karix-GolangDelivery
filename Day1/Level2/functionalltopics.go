
func signalStatus(tower string, strength int) (string, bool) {
	if strength > -80 {
		return "Good", true
	}
	return "Weak", false
}

func totalUsers(usersPerTower ...int) int {
	total := 0
	for _, u := range usersPerTower {
		total += u
	}
	return total
}


calcProfit := func(buy, sell float64) float64 {
	return sell - buy
}

func batchOrders(n int) int {
	if n == 0 {
		return 1
	}
	return n * batchOrders(n-1)
}


func openDB() {
	defer fmt.Println("Closing DB connection")
	fmt.Println("DB Connection Opened")
}

func processPayment(amount float64) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from error:", r)
		}
	}()
	if amount < 0 {
		panic("Invalid transaction amount")
	}
	fmt.Println("Payment processed:", amount)
}



func main(){
	
}
