

package main

import (
	"errors"
	"fmt"
)

// Function to do a task and return error if failed
func doTask(task string) error {
	if task == "cook_biryani" {
		return nil // success
	} else if task == "clean_room" {
		return errors.New("room too messy to clean") // failed task
	}
	return nil
}

func main() {
	tasks := []string{"cook_biryani", "clean_room"}

	for _, task := range tasks {
		err := doTask(task)
		if err != nil {
			fmt.Printf("❌ Task '%s' failed: %s\n", task, err)
		} else {
			fmt.Printf("✅ Task '%s' completed successfully\n", task)
		}
	}
}
