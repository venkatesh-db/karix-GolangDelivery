package main

import "fmt"

/*

Concept:

Function returns two things → grade + pass/fail

Human analogy = multiple outcomes from one action

*/

// Human analogy: Check exam scores and return grade + pass/fail

func checkScore(score int) (string, bool) {
	if score >= 50 {
		return "Pass", true
	}
	return "Fail", false
}

// Human analogy: Family members share their blessings

func totalBlessings(blessings ...int) int {
	sum := 0
	for _, b := range blessings {
		sum += b
	}
	return sum
}

func inline() {
	// Inline / anonymous function
	double := func(x int) int {
		return x * 2
	}
	fmt.Println("Double 5 =", double(5))
}

func factorial(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorial(n-1)
}

func safeDivision(a, b int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	fmt.Println("Result:", a/b)
}

func main() {

	defer fmt.Println("Cleanup after all tasks")

	grade, passed := checkScore(65)
	fmt.Printf("Grade: %s, Passed: %t\n", grade, passed)

	fmt.Println("Total Blessings:", totalBlessings(5, 10, 7))

	fmt.Println("Factorial of 5 =", factorial(5))

	safeDivision(10, 0) // division by zero
	fmt.Println("Program continues safely")

}
