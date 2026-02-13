package main

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed english_rights.txt
var en_rights string

//go:embed french_rights.txt
var fr_rights string

//go:embed german_rights.txt
var ge_rights string

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "en":
			fmt.Println(en_rights)
		case "fr":
			fmt.Println(fr_rights)
		case "ge":
			fmt.Println(ge_rights)
		}
		os.Exit(0)
	}
	fmt.Println("Input the country code [en/fr/ge]")
}
