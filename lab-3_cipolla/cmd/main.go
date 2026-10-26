package main

import (
	"fmt"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-3_cipolla/cipolla"
)

func main() {
	fmt.Println("This program performs calculation of finite firld square root via Cipolla's alogithm")

	fmt.Println("Input A:")
	var a int64
	fmt.Scanf("%d", &a)
	if a < 0 {
		fmt.Println("Invalid A")
		return
	}
	fmt.Println("Input P:")
	var p int64
	fmt.Scanf("%d", &p)
	if p < 0 {
		fmt.Println("Invalid P")
		return
	}

	res, ok, err := cipolla.FiniteFieldSqrt(a, p)
	if err != nil {
		fmt.Printf("Encountered an error during calculation: %s\n", err)
		return
	}

	if !ok {
		fmt.Printf("%d is not a squadratic residue mod %d\n", a, p)
		return
	}

	fmt.Printf("One of the roots x1 is %d, other x2 = p - x1 = %d\n", res, p-res)
}
