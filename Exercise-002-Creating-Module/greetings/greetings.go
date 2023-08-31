package greetings

import "fmt"

func Hello(name string) string {
	// %v is a placeholder for the variable name in the format string
	// Sprintf formats and returns a string without printing it anywhere
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}
