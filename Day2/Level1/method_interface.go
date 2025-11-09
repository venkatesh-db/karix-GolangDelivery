


package main

import "fmt"

// Struct + method: simple human example
type Person struct {
    Name string
    Age  int
}

// method assigned to struct (value receiver okay here)
func (p Person) Greet() string {
    return fmt.Sprintf("Hello, I'm %s and I'm %d years old.", p.Name, p.Age)
}

// Interface: someone who can greet
type Greeter interface {
    Greet() string
}

func sayHello(g Greeter) {
    // polymorphism: any Greeter works
    fmt.Println(g.Greet())
}

func main() {
    bmw := Person{Name: "venkat", Age: 30}
    skoda := Person{Name: "jamesbond", Age: 25}

    // both implement Greeter via Greet()
    sayHello(bmw)
    sayHello(skoda)
}
