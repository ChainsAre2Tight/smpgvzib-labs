package main

import (
	"fmt"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-2_legendre/legendre"
)

func main() {
	fmt.Println("This program performs calculation of Legendre (Jakobi) symbol (A/N)")

	fmt.Println("Input A:")
	var a int64
	fmt.Scanf("%d", &a)
	if a < 0 {
		fmt.Println("Invalid A")
		return
	}
	fmt.Println("Input N:")
	var n int64
	fmt.Scanf("%d", &n)
	if n < 0 {
		fmt.Println("Invalid N")
		return
	}

	res, err := legendre.JakobiSymbol(a, n)
	if err != nil {
		fmt.Printf("Encountered an error during calculation: %s\n", err)
		return
	}

	fmt.Printf("Legendre (Jakobi) symbol for %d / %d is %d\n", a, n, res)
}
