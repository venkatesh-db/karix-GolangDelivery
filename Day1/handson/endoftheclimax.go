package main

import (
	"errors"
	"fmt"
)

func sum(numbers ...int) {
	total := 0

	for _, n := range numbers {

		total += n

	}

	fmt.Println("total numbers:", total)
}

func passthrough(sig string) (int, error) {

	if sig == "red" {

		captures := "vch2345"
		policemama := "angyr bird face 1000"

		// if he see vehcile he stops  throw exception
		//  u will fees --> police fees 2000

		fmt.Println("fine & police mamn hug", captures, policemama)
		return 1000, errors.New("u will vehcile will kept in policestation")

	} else if sig == "green" {

	} else {

	}

	fmt.Println("This is passthrough function")
	return 0, nil
}

func signal(sig string) {

	switch sig {
	case "red":
		fine, errros := passthrough("red")
		fmt.Println("fine amount:", fine, "error message:", errros)
	case "green":
		passthrough("green")
	default:
		fmt.Println("human passthrough Signal")
	}

}

func telecom() {

	cloud := map[string]int{
		"TowerA": 78,
		"TowerB": 45,
		"TowerC": 89,
	}

	clous := map[string]int{
		"ec2": 100,
		"s3":  200,
		"rds": 300,
	}

	for k, v := range clous {
		println("Tower:", k, "Load:", v)
	}

	fmt.Println("Cloud Services Loaded", clous, cloud)

}

func main() {

	telecom()
	sum(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	signal("red")

}
