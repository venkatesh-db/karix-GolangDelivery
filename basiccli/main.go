package main

import (
	"flag"
	"fmt"
)

// cli code to run commands with flags
// go run main.go --name=Venkatesh --age=23 camerboy



func main() {

	name := flag.String("name", "guest", "your name")
	age := flag.Int("age", 0, "your age")

	flag.Parse()

	fmt.Printf("Hello %s, you are %d years old.\n", *name, *age)

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("No additional arguments provided.")
		return
	}

	cmd := args[0]
	
	switch cmd {
	case "break":
		fmt.Printf("Greetings, %s! You are %d years old.\n", *name, *age)
	case "camerboy":
		fmt.Printf("Goodbye, %s!\n", *name)
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
	}

}
