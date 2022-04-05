package utils

import "fmt"

// Status provides logging during setup
type Status struct {
	str string
}

// LogStatus creates the status object
func LogStatus(input string) Status {
	fmt.Print(input)
	return Status{input}
}

// Exists prints an additional note when called
func (s Status) Exists() {
	fmt.Print(" - already exists")
}

// Check parses and prints accordingly
func (s *Status) Check(e error) error {
	if e != nil {
		fmt.Println(" ❌")
		return e
	}
	fmt.Println(" ✅")
	return e
}
