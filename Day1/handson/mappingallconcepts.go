
// wedding

// --Concpets
// datatypes
// conditions
// loops
// continous memory array
// dyamic memory - map

// food        -->          map array
// photooghy   -->         group  duseetti --> family members
// bride groom -->         map prewedding post wedding
// decoration   -->        map
// gifts         --->      map  aman --> smile
// meeting people -->      continous arraay
// girl planning  -->     new prposal  array

package main

import "fmt"

func food(mhall []string) {

	// 10 items we can reepast two ietms

}

func photooghy(photos map[string][]string) {

	fmt.Println(photos)

}

func bridegroom() {

}

func decoration() {

}

func gifts(giftse map[string]interface{}) {

	fmt.Println(giftse)

}

func meetingpeople(meetpeople map[string][]string) {

	fmt.Println(meetpeople)
}

func weddingpair(wedcards []string, wedcarrds map[string]string) {

	fmt.Println(wedcards, wedcarrds)

}

func main() {

	wedcards := []string{"venkatesh", "seetha"}
	wedcarrds := map[string]string{"dusetti": "venkatesh",
		"rama": "seetha"}
	weddingpair(wedcards, wedcarrds)

	meetpeople := map[string][]string{
		"bride": []string{"mamma", "aunty"},
		"groom": []string{"father", "mother"},
	}

	meetingpeople(meetpeople)

	gift := map[string]interface{}{
		"venkatesh": "handsshake",
		"sagar":     "bokke",
		"manager":   5000,
	}

	gifts(gift)

	photos := map[string][]string{

		"freinds":        []string{"man1", "mabn2", "lady1", "lady2"},
		"glassbatch":     []string{"man1", "mabn2", "man2", "man4"},
		"relativefamily": []string{"uncle1", "uncle2", "aunty1", "younggirls"},
	}

	photooghy(photos)

	// array fixed
	var fixedfood [5]string = [5]string{"sweet1", "swewet2", "veg rice", "icecream", "white rice"}

	//slice is append
	var mhall []string = []string{"sweet1", "swewet2", "veg rice", "icecream", "white rice"}

	mhall = append(mhall, "rice")
	mhall = append(mhall, "sweet1")

	fmt.Println(fixedfood, mhall)

	food(mhall)

}


