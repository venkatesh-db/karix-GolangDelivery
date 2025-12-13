package main

import "fmt"

func MasterF1[T any](value T) T {
	return value
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

func Smile[T Number](a, b T) T {
	return a + b
}

func main() {

	fmt.Println(MasterF1[int](10))
	fmt.Println(MasterF1[string]("deepak ring master"))

	var laddu interface{}
	laddu = "sweet"
	laddu = 10000
	fmt.Println(laddu)

	fmt.Println(Smile[int](10, 20))
	fmt.Println(Smile[float64](10.5, 20.3))

}
