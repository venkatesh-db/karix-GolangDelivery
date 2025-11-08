package main

import "fmt"

func main() {
	// 1️⃣ Array — fixed family members
	family := [3]string{"Father", "Mother", "Daughter"}
	fmt.Println("👨‍👩‍👧 Family Members Array:", family)

	// 2️⃣ Slice — dynamic shopping list
	shoppingList := []string{"Milk", "Eggs", "Biryani"}
	shoppingList = append(shoppingList, "Fruits")
	fmt.Println("🛒 Shopping List Slice:", shoppingList)

	// 3️⃣ Map — person to favorite activity
	favorites := map[string]string{
		"Father":   "Reading",
		"Mother":   "Cooking",
		"Daughter": "Coding",
	}
	fmt.Println("❤️ Favorites Map:", favorites)
}
