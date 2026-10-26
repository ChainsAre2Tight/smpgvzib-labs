package cipolla_test

import (
	"fmt"
	"testing"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-3_cipolla/cipolla"
)

func TestCipollaBehaviour(t *testing.T) {
	_, _, err := cipolla.FiniteFieldSqrt(1, 4)
	if err == nil || err.Error() != "cipolla.FiniteFieldSqrt: Expected P to be a prime odd number, got 4" {
		t.Fatalf("Expected to fail on P == 4, got %v", err)
	}
	_, ok, err := cipolla.FiniteFieldSqrt(5, 7)
	if err != nil || ok {
		t.Fatalf("Expected to fail on 5|7 with not ok, got %v|%v", ok, err)
	}
	_, ok, err = cipolla.FiniteFieldSqrt(4, 5)
	if err != nil || !ok {
		t.Fatalf("Expected all ok (4|5), got %v|%v", ok, err)
	}
}

func TestCipolla(t *testing.T) {
	tt := []struct {
		a, p, res1, res2 int64
	}{
		{a: 1, p: 3, res1: 1, res2: 2},
		{a: 4, p: 5, res1: 2, res2: 3},
		{a: 2, p: 7, res1: 3, res2: 4},
		{a: 3, p: 11, res1: 5, res2: 6},
		{a: 10, p: 13, res1: 6, res2: 7},
		{a: 8, p: 17, res1: 5, res2: 12},
		{a: 6, p: 19, res1: 5, res2: 14},
		{a: 8, p: 23, res1: 10, res2: 13},
		{a: 5, p: 29, res1: 11, res2: 18},
		{a: 10, p: 31, res1: 14, res2: 17},
		{a: 12, p: 37, res1: 7, res2: 30},
		{a: 8, p: 41, res1: 7, res2: 34},
		{a: 6, p: 43, res1: 7, res2: 36},
		{a: 8, p: 47, res1: 14, res2: 33},
		{a: 11, p: 53, res1: 8, res2: 45},
		{a: 5, p: 59, res1: 8, res2: 51},
		{a: 3, p: 61, res1: 8, res2: 53},
		{a: 6, p: 67, res1: 26, res2: 41},
		{a: 10, p: 71, res1: 9, res2: 62},
		{a: 8, p: 73, res1: 9, res2: 64},
	}

	for _, td := range tt {
		t.Run(
			fmt.Sprintf("SQRT(%d, %d) == %d|%d", td.a, td.p, td.res1, td.res2),
			func(t *testing.T) {
				v, ok, err := cipolla.FiniteFieldSqrt(td.a, td.p)
				if err != nil {
					t.Fatalf("Non-nill error: %s", err)
				}
				if !ok {
					t.Fatalf("Not ok...")
				}
				if v != td.res1 && v != td.res2 {
					t.Fatalf("\nWrong result: Expected %d or %d, got %d", td.res1, td.res2, v)
				}
			},
		)
	}
}
