package main

import "fmt"

type Junkfood interface {
	taste(string)
}

type Stress interface {
	Workpressure(string)
}

type Healthy interface {
	Stress
	Junkfood
	nutrition()
}

type ITProfessional struct {
	toungstatste []string // Junkfood interface

	protienfood []string // Healthy interface
	fruits      []string // Healthy interface

	bodyfats []string // comon all interfaces

	deadlines  []string // Stress interface
	relaxation []string // Stress interface
}

func (itp *ITProfessional) taste(hungry string) {

	itp.toungstatste = []string{}
	itp.toungstatste = append(itp.toungstatste, hungry)

	fmt.Println("IT professional loves to eat junk food for its taste")
}

func (itp *ITProfessional) nutrition() {

	fmt.Println("IT professional needs healthy food for nutrition")
}

func (itp *ITProfessional) Workpressure(energy string) {

	fmt.Println("IT professional faces work pressure")

	itp.deadlines = []string{}
	itp.deadlines = append(itp.deadlines, energy)

}

func main() {

	// 6 days  only object creation and method calling for interface type Healthy

	var yooung Healthy
	yooung = &ITProfessional{}
	yooung.nutrition()
	yooung.Workpressure("energry low")
	yooung.taste("very hungry")

	// Lib provdies interface we need imoplemnt struct & methods

}
