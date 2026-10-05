package main

import (
	"fmt"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-1_montgomery/montgomery"
	"github.com/ChainsAre2Tight/smpvzib-labs/lab-1_montgomery/mpi"
)

func main() {
	fmt.Println("This program performs multiplication of 2 multi-precision integers via Montgomery multiplication")
	fmt.Println("A * B mod N")
	fmt.Println("STD input is supported for radix 10")
	radix := 10
	fmt.Printf("Input A in base-%d:\n", radix)
	var a int64
	fmt.Scanf("%d", &a)
	if a < 0 {
		fmt.Println("Invalid A")
		return
	}
	fmt.Printf("Input B in base-%d:\n", radix)
	var b int64
	fmt.Scanf("%d", &b)
	if b < 0 {
		fmt.Println("Invalid B")
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
	B := mpi.FromZnBaseB(uint64(b), uint32(radix))
	N := mpi.FromZnBaseB(uint64(n), uint32(radix))
	fmt.Printf("A (%d) as a little-endian multi-precision integer in radix %d is\n%v\n", a, radix, A.V)
	fmt.Printf("B (%d) as a little-endian multi-precision integer in radix %d is\n%v\n", b, radix, B.V)
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
	mB, err := mont.Phi(B)
	if err != nil {
		fmt.Printf("Encountered an error during reduction of B: %s", err)
		return
	}
	fmt.Printf("Phi(B) = Redc(b * R) = %vb%d (%d)\n", mB.V, radix, mB.ToZn())
	fmt.Println("\nCalculating MontMul(Phi(a), Phi(b)):")
	mR, err := mont.MontMul(mA, mB)
	if err != nil {
		fmt.Printf("Encountered an error during MontMul(Phi(a)*Phi(b)): %s", err)
		return
	}
	fmt.Printf("Phi(a)*Phi(B) mod N = %vb%d (%d)\n", mR.V, radix, mR.ToZn())
	fmt.Println("\nApplying Phi^-1 to get result as an element of Zn:")
	res, err := mont.InversePhi(mR)
	if err != nil {
		fmt.Printf("Encountered an error during application of InversePhi(A*B): %s", err)
		return
	}
	fmt.Printf("result = Phi^-1(Phi(A)*(Phi(B))) = %vb%d (%d)\n", res.V, radix, res.ToZn())
	fmt.Printf("Expected result was A*BmodN = %d\n", (int64(A.ToZn())*int64(B.ToZn()))%int64(N.ToZn()))
}
