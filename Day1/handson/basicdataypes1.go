

package main
import "fmt". // fmt --> folder 

/*
model sedan
price 11.16 lakhs
color white
diesel true
*/

func main() {

	var model string = "sedan"    // 1 memory  --> separate
	var price float64 = 110000.16 // 2 mmeory  --> separate
	var color string = "white"    // 3 memory  --> separate
	var dieles bool = true        // 4 memory  --> separate
	
	var cars [4]string = [4]string{"sedan", "110000.16", "white", "true"} // continous 4 memory to store data 
	// No of way storing data in one memory  --> array struct map

	fmt.Println(model, price, color, dieles, cars)

}
