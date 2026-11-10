package miller

import (
	"fmt"
	"math/rand"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-4_primality/utils"
)

func MillerTest(N int64, rounds int) (res bool, err error) {
	fail := func(err error) (bool, error) {
		return false, fmt.Errorf("strassen.MillerTest: %s", err)
	}

	if N <= 0 {
		return fail(fmt.Errorf("N must be positive"))
	}

	if rounds == 0 {
		return fail(fmt.Errorf("Rounds should be non-zero"))
	}

	if N < 2 {
		return false, nil
	}
	if N == 2 || N == 3 {
		return true, nil
	}
	if N%2 == 0 {
		return false, nil
	}

	d := N - 1
	s := 0
	for d%2 == 0 {
		d /= 2
		s++
	}

	res = true

	for range rounds {
		a := rand.Int63n(N-3) + 2

		g, _, _ := utils.ExtendedGCD(a, N)
		if g != 1 {
			res = false
			break
		}

		x := int64(utils.ModPow(uint64(a), uint64(d), uint64(N)))
		if x == 1 || x == N-1 {
			continue
		}

		composite := true
		for i := 0; i < s-1; i++ {
			x = int64(utils.ModPow(uint64(x), 2, uint64(N)))
			if x == N-1 {
				composite = false
				break
			}
		}

		if composite {
			res = false
			break
		}
	}

	return res, nil
}
