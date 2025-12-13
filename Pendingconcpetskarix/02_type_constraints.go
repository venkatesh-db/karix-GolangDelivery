package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// ===== Built-in Type Constraints =====

// any - accepts any type
func PrintAny[T any](value T) {
	fmt.Printf("%v\n", value)
}

// comparable - types that support == and !=
func Equal[T comparable](a, b T) bool {
	return a == b
}

// ===== Custom Type Constraints =====

// Number constraint using union of types
type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

// Add function using Number constraint
func Add[T Number](a, b T) T {
	return a + b
}

// Multiply function using Number constraint
func Multiply[T Number](a, b T) T {
	return a * b
}

// ===== Using constraints package =====

// Max function using Ordered constraint
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Sum function using constraints package
func Sum[T constraints.Integer | constraints.Float](numbers []T) T {
	var total T
	for _, n := range numbers {
		total += n
	}
	return total
}

// ===== Interface-based Constraints =====

// Stringer constraint - types that have String() method
type Stringer interface {
	String() string
}

func PrintString[T Stringer](value T) {
	fmt.Println(value.String())
}

// Custom type implementing Stringer
type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s (age %d)", p.Name, p.Age)
}

// ===== Method Constraints =====

// Addable constraint - types with Add method
type Addable[T any] interface {
	Add(T) T
}

// Complex number type implementing Addable
type Complex struct {
	Real, Imag float64
}

func (c Complex) Add(other Complex) Complex {
	return Complex{
		Real: c.Real + other.Real,
		Imag: c.Imag + other.Imag,
	}
}

func (c Complex) String() string {
	return fmt.Sprintf("%.2f + %.2fi", c.Real, c.Imag)
}

// Generic function using Addable constraint
func SumAddable[T Addable[T]](values []T) T {
	if len(values) == 0 {
		var zero T
		return zero
	}
	result := values[0]
	for i := 1; i < len(values); i++ {
		result = result.Add(values[i])
	}
	return result
}

// ===== Combining Constraints =====

// StringableNumber - types that are numbers AND have String() method
type StringableNumber interface {
	constraints.Ordered
	String() string
}

// MyInt with String() method
type MyInt int

func (m MyInt) String() string {
	return fmt.Sprintf("MyInt(%d)", m)
}

// ===== Constraint with approximate element =====

// Unsigned constraint using ~
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// IsEven works with any unsigned integer type (including custom types)
func IsEven[T Unsigned](n T) bool {
	return n%2 == 0
}

type CustomUint uint

// ===== Slice Constraint =====

// SliceConstraint - accepts any slice type
type SliceConstraint[T any] interface {
	~[]T
}

func Length[S SliceConstraint[T], T any](slice S) int {
	return len(slice)
}

// ===== Map Constraint =====

// MapConstraint - accepts any map type
type MapConstraint[K comparable, V any] interface {
	~map[K]V
}

func Keys[M MapConstraint[K, V], K comparable, V any](m M) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// ===== Multiple Constraint Example =====

// Numeric operations with constraint
type Numeric interface {
	constraints.Integer | constraints.Float
}

func Average[T Numeric](numbers []T) float64 {
	if len(numbers) == 0 {
		return 0
	}
	var sum T
	for _, n := range numbers {
		sum += n
	}
	return float64(sum) / float64(len(numbers))
}

func main() {
	fmt.Println("=== Type Constraints Demo ===\n")

	// 1. Built-in constraints
	fmt.Println("1. Built-in Constraints:")
	PrintAny(42)
	PrintAny("hello")
	fmt.Printf("Equal(5, 5): %v\n", Equal(5, 5))
	fmt.Printf("Equal(\"a\", \"b\"): %v\n", Equal("a", "b"))
	fmt.Println()

	// 2. Custom Number constraint
	fmt.Println("2. Custom Number Constraint:")
	fmt.Printf("Add(10, 20): %d\n", Add(10, 20))
	fmt.Printf("Add(3.14, 2.71): %.2f\n", Add(3.14, 2.71))
	fmt.Printf("Multiply(5, 7): %d\n", Multiply(5, 7))
	fmt.Println()

	// 3. Constraints package (Ordered)
	fmt.Println("3. Constraints Package:")
	fmt.Printf("Max(10, 20): %d\n", Max(10, 20))
	fmt.Printf("Max(3.14, 2.71): %.2f\n", Max(3.14, 2.71))
	fmt.Printf("Max(\"apple\", \"banana\"): %s\n", Max("apple", "banana"))
	fmt.Println()

	// 4. Sum with Integer/Float constraint
	fmt.Println("4. Sum Function:")
	intSlice := []int{1, 2, 3, 4, 5}
	fmt.Printf("Sum of %v: %d\n", intSlice, Sum(intSlice))
	floatSlice := []float64{1.1, 2.2, 3.3}
	fmt.Printf("Sum of %v: %.2f\n", floatSlice, Sum(floatSlice))
	fmt.Println()

	// 5. Interface-based constraints (Stringer)
	fmt.Println("5. Interface-based Constraints:")
	person := Person{Name: "Alice", Age: 30}
	PrintString(person)
	fmt.Println()

	// 6. Method constraints (Addable)
	fmt.Println("6. Method Constraints:")
	complexNums := []Complex{
		{Real: 1, Imag: 2},
		{Real: 3, Imag: 4},
		{Real: 5, Imag: 6},
	}
	complexSum := SumAddable(complexNums)
	fmt.Printf("Sum of complex numbers: %s\n", complexSum)
	fmt.Println()

	// 7. Unsigned constraint with ~
	fmt.Println("7. Unsigned Constraint (with ~):")
	var normalUint uint = 10
	var customUint CustomUint = 15
	fmt.Printf("IsEven(%d): %v\n", normalUint, IsEven(normalUint))
	fmt.Printf("IsEven(%d): %v\n", customUint, IsEven(customUint))
	fmt.Println()

	// 8. Slice constraint
	fmt.Println("8. Slice Constraint:")
	numbers := []int{1, 2, 3, 4, 5}
	fmt.Printf("Length of slice: %d\n", Length(numbers))
	fmt.Println()

	// 9. Map constraint
	fmt.Println("9. Map Constraint:")
	m := map[string]int{"one": 1, "two": 2, "three": 3}
	fmt.Printf("Keys: %v\n", Keys(m))
	fmt.Println()

	// 10. Average with Numeric constraint
	fmt.Println("10. Average Function:")
	intNums := []int{10, 20, 30, 40, 50}
	fmt.Printf("Average of %v: %.2f\n", intNums, Average(intNums))
	floatNums := []float64{1.5, 2.5, 3.5, 4.5}
	fmt.Printf("Average of %v: %.2f\n", floatNums, Average(floatNums))
}
