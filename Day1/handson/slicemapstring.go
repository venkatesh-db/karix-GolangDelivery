
package main

import (
	"fmt"
	"strings"
)

func main() {

	latestnews := []string{"germany", "tech-golang", "1cr", "company sponserd"}
	slice1 := latestnews[1:3]

	fmt.Println(latestnews, slice1)
	copys := make([]string, 4) // heap
	copy(copys, latestnews)
	fmt.Println(copys)

	latestnews = append(latestnews, "bmw car")
	fmt.Println(latestnews)

	str := "smiling learning"

	fmt.Println(len(str))
	fmt.Println(string(str[0]))

	fmt.Println(strings.ToUpper(str))
	fmt.Println(strings.Contains("today greatest learning", "learning"))
	fmt.Println(strings.Index("today greatest learning", "l"))
	fmt.Println(strings.Replace("go go go", "go", "run", 3))

}
