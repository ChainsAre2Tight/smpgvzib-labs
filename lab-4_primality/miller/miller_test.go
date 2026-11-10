package miller_test

import (
	"fmt"
	"testing"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-4_primality/miller"
)

func TestMiller(t *testing.T) {
	tt := []struct {
		a   int64
		res bool
	}{
		{1, false},
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{9, false},
		{11, true},
		{13, true},
		{17, true},
		{97, true},
		{7919, true},
		{15, false},
		{21, false},
		{25, false},
		{100, false},
		{561, false},
		{1105, false},
		{1729, false},
		{2465, false},
	}

	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%d is prime: %t", td.a, td.res),
			func(t *testing.T) {
				prime, err := miller.MillerTest(td.a, 1024)
				if err != nil {
					t.Fatalf("Non-nill error: %s", err)
				}
				if prime != td.res {
					t.Fatalf("\nWrong result: Expected %t, got %t", td.res, prime)
				}
			},
		)
	}
}
