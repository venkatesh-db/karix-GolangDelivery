package main

import "fmt"

func moveiefansaman() {

	simplebreakafst := []string{"eggs", "bread", "butter", "jam"}

	var myexperience = map[string][]string{
		"kenthomson": {"action", "comedy", "golang"},
		"rob piker":  {"drama", "history", "python"},
	}

	thinktowatchmovie(myexperience, simplebreakafst)

	fmt.Println(myexperience, simplebreakafst, " in main function")
}

func thinktowatchmovie(myexperience map[string][]string, simplebreakafst []string) {

	myexperience["villans"] = []string{"horror", "thriller"}
	fmt.Println(myexperience, "thinktowatchmovie")

	simplebreakafst = append(simplebreakafst, "burger")
	fmt.Println(simplebreakafst, "thinktowatchmovie")

}

func main() {

	moveiefansaman()

}
