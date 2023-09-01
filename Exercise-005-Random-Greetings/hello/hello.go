package main

import (
	"fmt"
	"log"

	"example.com/greetings"

)

func main() {
	log.SetPrefix("Greetings: ")
	log.SetFlags(0)

	names := []string{"Eris", "Safi", "Kev", "Juju"}

	messages, err := greetings.Hellos(names)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(messages)
}
