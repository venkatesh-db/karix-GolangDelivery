package main

import "fmt"

func chcoloates() { // function definition

	// separte memory

	var falvour string = "nuts"
	var weight int16 = 150
	var price int16 = 1000
	var todaydate string = "27-oct-2025"
	var mfgmonth int8 = 2             // 12 months
	var totalgap int8 = 10 - mfgmonth // 8 months

	if mfgmonth < 10 && totalgap < 3 { // 2 month < 10 month

		var wishtobuy int8 = 40 // pay to chocolte
		fmt.Println("wish to buy", wishtobuy)
	} else {

		fmt.Println(" healthy person")
	}
	// continous memory 4 memory. interface{} --> any type of data string int float bool
	var chcolates []interface{} = []interface{}{"nuts", 150, 1000, "27-oct-2025"}
	fmt.Println(chcolates, falvour, weight, price, todaydate)

}

func main() {
	chcoloates() // to execute code  function call
}
