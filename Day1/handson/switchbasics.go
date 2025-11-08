
package main

import "fmt"

/*
movie :"golang by ken thomoson"
theatre: 2d 4d imax 4dx
pricing: 200 400 500 700
*/

func main() {

	var movie string = "golang by ken thomoson"
	var theatre string = "4dx"
	var spend int16

	switch theatre {

	case "2d":
		spend = 200
		fmt.Println(" 2dx exp", spend, movie)

	case "4d":

		spend = 400
		fmt.Println(" 4d exp", spend)

	case "imax":
		spend = 500
		fmt.Println(" imax exp", spend)

	case "4dx":

		spend = 700
		fmt.Println(" 4dx exp", spend)

	default:
		fmt.Println(" golang is booked")

	}

}
