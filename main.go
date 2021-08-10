package main

import (
	"fmt"
	"github.com/ymakhloufi/avatarme/identicon"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You need to pass the textual identifier as an argument.")
		os.Exit(1)
	}
	identifier := os.Args[1]

	identiconObj := identicon.New(identifier)

	fmt.Println(identiconObj)
}