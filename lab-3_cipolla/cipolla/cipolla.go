package cipolla

import (
	"fmt"
	"math/rand/v2"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-2_legendre/legendre"
)

func FiniteFieldSqrt(a, p int64) (res int64, ok bool, err error) {
	fail := func(err error) (int64, bool, error) {
		return 0, false, fmt.Errorf("cipolla.FiniteFieldSqrt: %s", err)
	}

	if p%2 == 0 || p <= 0 {
		return fail(fmt.Errorf("Expected P to be a prime odd number, got %d", p))
	}

	jakobi, err := legendre.JakobiSymbol(a, p)
	if err != nil {
		return fail(fmt.Errorf("Jakobi symbol (%d | %d): %s", a, p, err))
	}

	if jakobi != 1 {
		return 0, false, nil
	}

	var b int64

	for {
		b = rand.Int64() % p
		val_to_test := b*b - 4*a
		for val_to_test < 0 {
			val_to_test += p
		}
		jak, err := legendre.JakobiSymbol(val_to_test, p)
		if err != nil {
			return fail(fmt.Errorf("Jakobi symbol b^2 - 4a (%d | %d): %s", val_to_test, p, err))
		}

		if jak == -1 {
			break
		}
	}

	y := Poly{1, 0}

	pow := (p + 1) / 2
	resPoly := Pow(y, pow, a, b, p)

	if resPoly.U != 0 {
		panic(fmt.Sprintf("cipolla: poly exp: ay is not zero: %d", resPoly.U))
	}

	return resPoly.V, true, nil
}

type Poly struct {
	U, V int64
}

func Mul(x, y Poly, a, b, p int64) (res Poly) {
	mod := func(x, p int64) int64 {
		x %= p
		if x < 0 {
			x += p
		}
		return x
	}

	return Poly{
		U: mod(x.U*y.V+y.U*x.V+x.U*y.U*b, p),
		V: mod(x.V*y.V-a*x.U*y.U, p),
	}
}

func Pow(base Poly, e int64, a, b, p int64) (res Poly) {
	res = Poly{0, 1}
	for e > 0 {
		if e&1 == 1 {
			res = Mul(res, base, a, b, p)
		}
		base = Mul(base, base, a, b, p)
		e >>= 1
	}

	return res
}
