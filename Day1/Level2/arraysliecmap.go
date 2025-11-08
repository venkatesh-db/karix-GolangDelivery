package main

import "fmt"

/*

🟢 Production Practices

Use arrays for fixed resources

Use slices for dynamic workloads

Use maps for lookup tables or key-value configs

*/

func main() {
	// Array — fixed server nodes
	serverNodes := [3]string{"ServerA", "ServerB", "ServerC"}
	fmt.Println("[PROD] Server Nodes Array:", serverNodes)

	// Slice — active users dynamically
	activeUsers := []string{"User101", "User102"}
	activeUsers = append(activeUsers, "User103")
	fmt.Println("[PROD] Active Users Slice:", activeUsers)

	// Map — userID → sessionID
	sessionMap := map[string]string{
		"User101": "SessionA1",
		"User102": "SessionB2",
	}
	sessionMap["User103"] = "SessionC3"
	fmt.Println("[PROD] Session Map:", sessionMap)
}
