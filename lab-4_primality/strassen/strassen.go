package strassen

import (
	"fmt"
	"math/rand"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-2_legendre/legendre"
	"github.com/ChainsAre2Tight/smpvzib-labs/lab-4_primality/utils"
)

func StrassenTest(N int64, rounds int) (res bool, err error) {
	fail := func(err error) (bool, error) {
		return false, fmt.Errorf("strassen.StrassenTest: %s", err)
	}

	if N <= 0 {
		return fail(fmt.Errorf("N must be positive"))
	}

	if rounds == 0 {
		return fail(fmt.Errorf("Rounds should be non-zero"))
	}

	if N == 1 {
		return false, nil
	}
	if N == 2 || N == 3 {
		return true, nil
	}
	if N%2 == 0 {
		return false, nil
	}

	res = true

	for range rounds {

		a := rand.Int63n(N-3) + 2

		if a == 0 {
			a += 1
		}

		d, _, _ := utils.ExtendedGCD(a, N)
		if d != 1 {
			res = false
			break
		}

		exp := int64(utils.ModPow(uint64(a), uint64((N-1)/2), uint64(N)))

		jakobi, err := legendre.JakobiSymbol(a, N)
		if err != nil {
			return fail(fmt.Errorf("jakobi (%d/%d): %s", a, N, err))
		}

		if jakobi == -1 && exp != N-1 || jakobi == 1 && exp != 1 || jakobi == 0 {
			res = false
			break
		}
	}

	return res, nil
}
