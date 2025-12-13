package main

import "fmt"

// Basic generic function - works with any type
func Print[T any](value T) {
	fmt.Printf("Value: %v, Type: %T\n", value, value)
}

// Generic function with multiple type parameters
func Pair[T any, U any](first T, second U) (T, U) {
	return first, second
}

// Generic function to find minimum using comparable types
func Min[T comparable](a, b T) T {
	// For comparable, we need ordered constraint for < operator
	// This is a simplified version
	return a // placeholder
}

// Generic function with ordered constraint (numbers)
// Ordered includes: integers, floats, strings
func MinOrdered[T int | float64 | string](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Generic function to get first element from slice
func First[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[0], true
}

// Generic function to get last element from slice
func Last[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[len(slice)-1], true
}

// Generic function to reverse a slice
func Reverse[T any](slice []T) []T {
	result := make([]T, len(slice))
	for i, v := range slice {
		result[len(slice)-1-i] = v
	}
	return result
}

// Generic function to filter slice based on predicate
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := []T{}
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Generic Map function to transform slice elements
func Map[T any, U any](slice []T, transform func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = transform(v)
	}
	return result
}

// Generic Reduce function
func Reduce[T any, U any](slice []T, initial U, reducer func(U, T) U) U {
	result := initial
	for _, v := range slice {
		result = reducer(result, v)
	}
	return result
}

// Generic function to check if slice contains element
func Contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

// Generic function to remove duplicates
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func main() {
	fmt.Println("=== Generic Functions Demo ===\n")

	// 1. Print function with different types
	fmt.Println("1. Generic Print Function:")
	Print(42)
	Print("Hello, Generics!")
	Print(3.14)
	Print(true)
	fmt.Println()

	// 2. Pair function
	fmt.Println("2. Pair Function:")
	name, age := Pair("Alice", 30)
	fmt.Printf("Name: %s, Age: %d\n", name, age)

	key, value := Pair(1, "one")
	fmt.Printf("Key: %d, Value: %s\n", key, value)
	fmt.Println()

	// 3. Min function with ordered types
	fmt.Println("3. Min Function:")
	fmt.Printf("Min(10, 20): %d\n", MinOrdered(10, 20))
	fmt.Printf("Min(3.14, 2.71): %.2f\n", MinOrdered(3.14, 2.71))
	fmt.Printf("Min(\"apple\", \"banana\"): %s\n", MinOrdered("apple", "banana"))
	fmt.Println()

	// 4. First and Last functions
	fmt.Println("4. First and Last Functions:")
	numbers := []int{1, 2, 3, 4, 5}
	if first, ok := First(numbers); ok {
		fmt.Printf("First: %d\n", first)
	}
	if last, ok := Last(numbers); ok {
		fmt.Printf("Last: %d\n", last)
	}
	fmt.Println()

	// 5. Reverse function
	fmt.Println("5. Reverse Function:")
	fmt.Printf("Original: %v\n", numbers)
	fmt.Printf("Reversed: %v\n", Reverse(numbers))

	words := []string{"Go", "is", "awesome"}
	fmt.Printf("Original: %v\n", words)
	fmt.Printf("Reversed: %v\n", Reverse(words))
	fmt.Println()

	// 6. Filter function
	fmt.Println("6. Filter Function:")
	evenNumbers := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})
	fmt.Printf("Even numbers: %v\n", evenNumbers)

	longWords := Filter(words, func(s string) bool {
		return len(s) > 2
	})
	fmt.Printf("Long words: %v\n", longWords)
	fmt.Println()

	// 7. Map function
	fmt.Println("7. Map Function:")
	squared := Map(numbers, func(n int) int {
		return n * n
	})
	fmt.Printf("Squared: %v\n", squared)

	lengths := Map(words, func(s string) int {
		return len(s)
	})
	fmt.Printf("Word lengths: %v\n", lengths)
	fmt.Println()

	// 8. Reduce function
	fmt.Println("8. Reduce Function:")
	sum := Reduce(numbers, 0, func(acc, n int) int {
		return acc + n
	})
	fmt.Printf("Sum: %d\n", sum)

	product := Reduce(numbers, 1, func(acc, n int) int {
		return acc * n
	})
	fmt.Printf("Product: %d\n", product)
	fmt.Println()

	// 9. Contains function
	fmt.Println("9. Contains Function:")
	fmt.Printf("Contains 3: %v\n", Contains(numbers, 3))
	fmt.Printf("Contains 10: %v\n", Contains(numbers, 10))
	fmt.Printf("Contains \"Go\": %v\n", Contains(words, "Go"))
	fmt.Println()

	// 10. Unique function
	fmt.Println("10. Unique Function:")
	duplicates := []int{1, 2, 2, 3, 3, 3, 4, 5, 5}
	fmt.Printf("With duplicates: %v\n", duplicates)
	fmt.Printf("Unique: %v\n", Unique(duplicates))
}
