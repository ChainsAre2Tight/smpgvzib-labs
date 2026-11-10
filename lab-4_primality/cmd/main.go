package main

import (
	"fmt"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-4_primality/fermat"
	"github.com/ChainsAre2Tight/smpvzib-labs/lab-4_primality/miller"
	"github.com/ChainsAre2Tight/smpvzib-labs/lab-4_primality/strassen"
)

func main() {
	fmt.Println("This program performs Fermat, Solovay-Strassen and Miller-Rabin primality tests")

	fmt.Println("Input A:")
	var a int64
	fmt.Scanf("%d", &a)
	if a < 0 {
		fmt.Println("Invalid A")
		return
	}
	fmt.Println("Input number of rounds:")
	var rounds int
	fmt.Scanf("%d", &rounds)
	if rounds < 0 {
		fmt.Println("Invalid number of rounds")
		return
	}

	prime, err := fermat.FermatTest(a, rounds)
	if err != nil {
		fmt.Printf("Encountered an error in FermatTest: %s\n", err)
		return
	}

	if prime {
		fmt.Printf("Number %d is considered to be prime by Fermat Test\n", a)
	} else {
		fmt.Printf("Number %d is considered to be non-prime by Fermat Test\n", a)
	}

	prime, err = strassen.StrassenTest(a, rounds)
	if err != nil {
		fmt.Printf("Encountered an error in StrassenTest: %s\n", err)
		return
	}

	if prime {
		fmt.Printf("Number %d is considered to be prime by Solovay-Strassen Test\n", a)
	} else {
		fmt.Printf("Number %d is considered to be non-prime by Solovay-Strassen Test\n", a)
	}

	prime, err = miller.MillerTest(a, rounds)
	if err != nil {
		fmt.Printf("Encountered an error in MillerTest: %s\n", err)
		return
	}

	if prime {
		fmt.Printf("Number %d is considered to be prime by Miller-Rabin Test\n", a)
	} else {
		fmt.Printf("Number %d is considered to be non-prime by Miller-Rabin Test\n", a)
	}
}
