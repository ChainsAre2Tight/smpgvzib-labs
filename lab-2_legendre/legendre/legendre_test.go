package legendre_test

import (
	"fmt"
	"testing"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-2_legendre/legendre"
)

func TestJakobiBehaviour(t *testing.T) {
	_, err := legendre.JakobiSymbol(1, -1)
	if err == nil || err.Error() != "JS(1, -1): B should be a positive integer" {
		t.Fatalf("Expected \"B should be a positive integer\" error, got: \n%s", err)
	}
	_, err = legendre.JakobiSymbol(1, 2)
	if err == nil || err.Error() != "JS(1, 2): B should be an odd integer" {
		t.Fatalf("Expected \"B should be an odd integer\" error, got: \n%s", err)
	}
}

func TestJakobiTwo(t *testing.T) {
	tt := []struct {
		b, res int64
	}{
		{1, 1}, {7, 1},
		{3, -1}, {5, -1},
		{9, 1}, {15, 1},
		{11, -1}, {13, -1},
	}

	for _, td := range tt {
		t.Run(
			fmt.Sprintf("JS(2, %d) == %d", td.b, td.res),
			func(t *testing.T) {
				v, err := legendre.JakobiSymbol(2, td.b)
				if err != nil {
					t.Fatalf("Non-nill error: %s", err)
				}
				if v != td.res {
					t.Fatalf("\nWrong result: Expected %d, got %d", td.res, v)
				}
			},
		)
	}
}

func TestJakobiMinus1(t *testing.T) {
	tt := []struct {
		b, res int64
	}{
		{1, 1}, {7, -1},
		{3, -1}, {5, 1},
		{9, 1}, {15, -1},
		{11, -1}, {13, 1},
	}

	for _, td := range tt {
		t.Run(
			fmt.Sprintf("JS(-1, %d) == %d", td.b, td.res),
			func(t *testing.T) {
				v, err := legendre.JakobiSymbol(-1, td.b)
				if err != nil {
					t.Fatalf("Non-nill error: %s", err)
				}
				if v != td.res {
					t.Fatalf("\nWrong result: Expected %d, got %d", td.res, v)
				}
			},
		)
	}
}

func TestJakobiModB(t *testing.T) {
	tt := []struct {
		a, b, res int64
	}{
		{3, 1, 1}, {9, 7, 1},
		{5, 3, -1}, {7, 5, -1},
		{11, 9, 1}, {17, 15, 1},
		{13, 11, -1}, {15, 13, -1},
	}

	for _, td := range tt {
		t.Run(
			fmt.Sprintf("JS(%d, %d) == %d", td.a, td.b, td.res),
			func(t *testing.T) {
				v, err := legendre.JakobiSymbol(td.a, td.b)
				if err != nil {
					t.Fatalf("Non-nill error: %s", err)
				}
				if v != td.res {
					t.Fatalf("\nWrong result: Expected %d, got %d", td.res, v)
				}
			},
		)
	}
}

func TestJakobiPower2(t *testing.T) {
	tt := []struct {
		a, b, res int64
	}{
		{3, 1, 1},
		{11, 7, 1},
		{9, 5, 1},
		{17, 9, 1}, {23, 15, 1},
		{8, 11, -1}, {8, 13, -1},
	}

	for _, td := range tt {
		t.Run(
			fmt.Sprintf("JS(%d, %d) == %d", td.a, td.b, td.res),
			func(t *testing.T) {
				v, err := legendre.JakobiSymbol(td.a, td.b)
				if err != nil {
					t.Fatalf("Non-nill error: %s", err)
				}
				if v != td.res {
					t.Fatalf("\nWrong result: Expected %d, got %d", td.res, v)
				}
			},
		)
	}
}

func TestJakobiGaussInvertion(t *testing.T) {
	tt := []struct {
		a, b, res int64
	}{
		{3, 5, -1},
		{3, 7, -1},
		{5, 7, -1},
		{5, 11, 1},
		{7, 11, -1},
		{9, 11, 1},
		{11, 15, -1},
		{13, 17, 1},
	}

	for _, td := range tt {
		t.Run(
			fmt.Sprintf("JS(%d, %d) == %d", td.a, td.b, td.res),
			func(t *testing.T) {
				v, err := legendre.JakobiSymbol(td.a, td.b)
				if err != nil {
					t.Fatalf("Non-nill error: %s", err)
				}
				if v != td.res {
					t.Fatalf("\nWrong result: Expected %d, got %d", td.res, v)
				}
			},
		)
	}
}
