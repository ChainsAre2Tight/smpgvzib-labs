package fermat

import (
	"fmt"
	"math/rand"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-4_primality/utils"
)

func FermatTest(N int64, rounds int) (res bool, err error) {
	fail := func(err error) (bool, error) {
		return false, fmt.Errorf("fermat.FermatTest: %s", err)
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

		d, _, _ := utils.ExtendedGCD(a, N)
		if d != 1 {
			res = false
			break
		}

		ferm := utils.ModPow(uint64(a), uint64(N-1), uint64(N))
		if ferm != 1 {
			res = false
			break
		}
	}

	return res, nil
}
