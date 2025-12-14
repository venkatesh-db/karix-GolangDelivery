package main

import (
	"netflix-event-demo/recommendationservice"
)

func main() {

	recommendationservice.InitServices()
	select {}

}
