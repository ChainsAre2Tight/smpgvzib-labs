package legendre

import "fmt"

func JakobiSymbol(a, b int64) (res int64, err error) {

	fail := func(err error) (int64, error) {
		return 0, fmt.Errorf("JS(%d, %d): %s", a, b, err)
	}

	if b <= 0 {
		return fail(fmt.Errorf("B should be a positive integer"))
	}

	if b%2 == 0 {
		return fail(fmt.Errorf("B should be an odd integer"))
	}

	res = 1

	for {
		// -1 ^ ((b-1)/2)
		if a == -1 {
			bmod4 := b & 0b11
			if bmod4 == 3 {
				return -res, nil
			} else {
				return res, nil
			}
		}

		if a == 1 {
			return res, nil
		}

		if a == 0 {
			if b == 1 {
				return res, nil
			}
			return 0, nil
		}

		if a < 0 {
			js, err := JakobiSymbol(-1, b)
			if err != nil {
				return fail(err)
			}
			res *= js
		}

		if a < 0 {
			a = -a
			continue
		}

		if a == 2 {
			bmod8 := b & 0b111
			neg := bmod8 == 3 || bmod8 == 5
			if neg {
				return -res, nil
			} else {
				return res, nil
			}
		}

		if a > b {
			a = a % b
			continue
		}

		power2 := 0
		for i := range 64 {
			bit := a&(1<<i) != 0
			if bit {
				break
			}
			power2++
		}

		if power2&1 != 0 {
			js, err := JakobiSymbol(2, b)
			if err != nil {
				return fail(err)
			}
			res *= js
		}

		a >>= int64(power2)
		if power2 != 0 {
			continue
		}

		invert := 1 - 2*((((a-1)*(b-1))&0b100)>>2)
		a, b = b, a
		res *= invert
		continue
	}
}
