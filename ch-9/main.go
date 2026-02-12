package main

import (
	"fmt"
	"errors"
)
// 1
var mySigErr = errors.New("Some error")
// 2
type myStructErr struct {
	Employee string
}

func (m myStructErr) Error() string {
	return fmt.Sprintf("StructErr: %w", m.Employee)
}

type otherStructErr struct {
	Field string
}

func (o otherStructErr) Error() string {
	return fmt.Sprintf("OtherErr: %w", o.Field)
}
// 3
type wrappedErrStruct struct {
	Errors []error
}

func (w wrappedErrStruct) Error() string {
	return "Yay"
}

func (w wrappedErrStruct) Unwrap() []error {
	return w.Errors
}

func handleErr(inp error) string {
	var sE myStructErr
	var sE2 otherStructErr
	if errors.Is(inp, mySigErr) {
		return "A signal error"
	} else if errors.As(inp, &sE) {
		return "Employee Struct error " + sE.Employee
	} else if errors.As(inp, &sE2) {
		return "Other error " + sE2.Field
	}
	return "Unknown error"
}

func main() {
	// 1+2
	sE := mySigErr
	fmt.Println(sE)
	structE := myStructErr{Employee: "Someone"}
	var structE2 otherStructErr
	fmt.Println(sE, structE)
	fmt.Println(errors.Is(sE, mySigErr))
	fmt.Println(errors.As(structE, &structE2))
	// 3
	uE := errors.New("Unknown error")
	wrappedErrors := wrappedErrStruct{Errors: []error{sE, structE}}
	errTable := [5]error{sE, structE, structE2, uE, wrappedErrors}
	accumulator := "" 
	for i := range errTable {
		switch err := errTable[i].(type) {
			case interface { Unwrap() []error }:
				errors := err.Unwrap()
				for j := range errors {
					accumulator += fmt.Sprintf("%s\n", handleErr(errors[j])) // produces extra newline. usage of slice preferred 
				}
			default:
				accumulator += fmt.Sprintf("%s\n", handleErr(err))
		}
	}
	fmt.Println(accumulator)
}
