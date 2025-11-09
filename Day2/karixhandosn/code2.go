package main

import "fmt"

func jamesbond() {

	var smile int8 = 20     // memory addreess 1 bytes
	var bond *int8 = &smile // memory addreess 4 bytes
	fmt.Println("1 bytes address", &smile)
	fmt.Println("bond", bond, *bond)

}

func songs() {

	var golangsongs []string = []string{"Golang is great", "I love coding in golang", "i love concurrency"}
	fmt.Println(golangsongs[0], golangsongs[1])

	var songhummings []*string = []*string{&golangsongs[0], &golangsongs[1], &golangsongs[2]}
	fmt.Println("songhumming address", songhummings, *songhummings[0], *songhummings[1], *songhummings[2])

}

func songpalyslist() {

	var songslist []string = []string{"F1", "lose my mind", "f1 therme song"}
	var favplaylist []*string = []*string{&songslist[0], &songslist[1], &songslist[2]}

	fmt.Println("favplaylist address", favplaylist, *favplaylist[0], *favplaylist[1], *favplaylist[2])

	songslist = append(songslist, "vadda podda")
	favplaylist = append(favplaylist, &songslist[3])
	fmt.Println(favplaylist)

}

func utubemusic() {

	type favourite struct {
		name          string
		next          *favourite
		times         int8
		mostfavoruite bool
	}

	mysongs := map[string]*[]string{
		"kollwood dance hitylist": &[]string{"F1", "lose my mind", "f1 therme song"},
		"bollywood recharger":     &[]string{"morni", "payal se", "aankh marey"},
	}


	/*

          map                     struct -SLL 

        "f1 music	"    -->      name : "F1" next : <address>
		"goloves"       -->      name : "lose my mind" next : <address>

	*/



	song3 := &favourite{name: "f1 therme song", times: 4, mostfavoruite: true, next: nil}

	song2 := &favourite{name: "lose my mind", times: 3, mostfavoruite: false, next: song3}
	song1 := &favourite{name: "F1", times: 5, mostfavoruite: true, next: song2}

	fmt.Println(song1, song2, song3)

	myfacvoritesongs := map[string]*favourite{
		"f1 music	": song1,
		"goloves":   song2,
		"go climax": song3,
	}

	fmt.Println(mysongs)

	fmt.Println("struct acess", myfacvoritesongs["f1 music	"],myfacvoritesongs["f1 music	"].next, myfacvoritesongs["f1 music	"].next.next)

}

func main() {

	jamesbond()
	songs()
	songpalyslist()
	utubemusic()
}
