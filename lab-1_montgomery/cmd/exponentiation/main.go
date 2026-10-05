package main

import (
	"fmt"
	"math"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-1_montgomery/montgomery"
	"github.com/ChainsAre2Tight/smpvzib-labs/lab-1_montgomery/mpi"
)

func main() {
	fmt.Println("This program performs exponentiation of a multi-precision integer via Montgomery multiplication")
	fmt.Println("A ^ c mod N")
	fmt.Println("Beware that large numbers might overflow Expected result field at the end as it uses direct int64 computations")
	fmt.Println("STD input is supported for radix 10")
	radix := 10
	fmt.Printf("Input A in base-%d:\n", radix)
	var a int64
	fmt.Scanf("%d", &a)
	if a < 0 {
		fmt.Println("Invalid A")
		return
	}
	fmt.Printf("Input c in base-%d:\n", radix)
	var c int64
	fmt.Scanf("%d", &c)
	if c < 0 {
		fmt.Println("Invalid c")
		return
	}
	fmt.Printf("Input N in base-%d so that it is co-prime with radix:\n", radix)
	var n int64
	fmt.Scanf("%d", &n)
	if n < 0 {
		fmt.Println("Invalid N")
		return
	}
	A := mpi.FromZnBaseB(uint64(a), uint32(radix))
	N := mpi.FromZnBaseB(uint64(n), uint32(radix))
	fmt.Printf("A (%d) as a little-endian multi-precision integer in radix %d is\n%v\n", a, radix, A.V)
	fmt.Printf("N (%d) as a little-endian multi-precision integer in radix %d is\n%v\n", n, radix, N.V)
	fmt.Println("Setting up Montgomery representation for modulo N:")
	mont, err := montgomery.NewMontgomery(N)
	if err != nil {
		fmt.Printf("Encountered an error during constant calculation: %s", err)
		return
	}
	fmt.Println("Montgomery parameters:")
	fmt.Printf("Modulo N: \t\t%vb%d (%d)\n", mont.N.V, radix, mont.N.ToZn())
	fmt.Printf("R = radix^Ln(N): \t%vb%d (%d)\n", mont.R.V, radix, mont.R.ToZn())
	fmt.Printf("-N^-1 mod R: \t\t%vb%d (%d)\n", mont.Ndashed.V, radix, mont.Ndashed.ToZn())
	fmt.Printf("R^2 mod N: \t\t%vb%d (%d)\n", mont.Rdashed.V, radix, mont.Rdashed.ToZn())
	fmt.Println("\nApplying Phi to A:")
	mA, err := mont.Phi(A)
	if err != nil {
		fmt.Printf("Encountered an error during reduction of A: %s", err)
		return
	}
	fmt.Printf("Phi(A) = Redc(a * R) = %vb%d (%d)\n", mA.V, radix, mA.ToZn())
	fmt.Println("\nCalculating MontExp(Phi(a), c):")
	mR, err := mont.MontExp(mA, uint64(c))
	if err != nil {
		fmt.Printf("Encountered an error during MontExp(Phi(a), c): %s", err)
		return
	}
	fmt.Printf("Phi(a)^c mod N = %vb%d (%d)\n", mR.V, radix, mR.ToZn())
	fmt.Println("\nApplying Phi^-1 to get result as an element of Zn:")
	res, err := mont.InversePhi(mR)
	if err != nil {
		fmt.Printf("Encountered an error during application of InversePhi(A*B): %s", err)
		return
	}
	fmt.Printf("result = Phi^-1(Phi(A)*(Phi(B))) = %vb%d (%d)\n", res.V, radix, res.ToZn())
	fmt.Printf("Expected result was A^c modN = %d\n", (int64(math.Pow(float64(int64(A.ToZn())), float64(c))) % int64(N.ToZn())))
}
