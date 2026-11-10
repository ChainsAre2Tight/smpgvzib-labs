package utils

import "math/bits"

func ExtendedGCD(a, b int64) (d, x, y int64) {
	signA, signB := int64(1), int64(1)
	if a < 0 {
		a, signA = -a, -1
	}
	if b < 0 {
		b, signB = -b, -1
	}

	oldR, r := a, b
	oldS, s := int64(1), int64(0)
	oldT, t := int64(0), int64(1)

	for r != 0 {
		q := oldR / r

		oldR, r = r, oldR-q*r
		oldS, s = s, oldS-q*s
		oldT, t = t, oldT-q*t
	}

	return oldR, signA * oldS, signB * oldT
}

func ModMul(a, b, m uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	_, rem := bits.Div64(hi, lo, m)
	return rem
}

func ModPow(base, exp, m uint64) uint64 {
	if m == 1 {
		return 0
	}
	base %= m
	result := uint64(1)
	for exp > 0 {
		if exp&1 == 1 {
			result = ModMul(result, base, m)
		}
		base = ModMul(base, base, m)
		exp >>= 1
	}
	return result
}
