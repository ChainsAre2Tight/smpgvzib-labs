package mpi_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-1_montgomery/mpi"
)

func TestMul(t *testing.T) {
	tt := []struct {
		name      string
		a, b, res mpi.MPI
	}{
		{
			"9712x526=5108512",
			mpi.MPI{B: 10, V: []int64{2, 1, 7, 9}},
			mpi.MPI{B: 10, V: []int64{6, 2, 5}},
			mpi.MPI{B: 10, V: []int64{2, 1, 5, 8, 0, 1, 5}},
		},
	}

	for _, td := range tt {
		t.Run(
			td.name,
			func(t *testing.T) {
				res, err := mpi.Mul(td.a, td.b)
				if err != nil {
					t.Fatalf("non-nill error in mpi.Mul: %s", err)
				}
				if !reflect.DeepEqual(res.V, td.res.V) {
					t.Fatalf("Wrong result:\nExpected\t%d\nGot\t\t%d", td.res.V, res.V)
				}
			},
		)
	}
}

func TestShortDiv(t *testing.T) {
	tt := []struct {
		name      string
		dividend  mpi.MPI
		divisor   int64
		quotent   mpi.MPI
		remainder int64
	}{
		{
			"8789/7=1255, 4",
			mpi.MPI{B: 10, V: []int64{9, 8, 7, 8}},
			7,
			mpi.MPI{B: 10, V: []int64{5, 5, 2, 1}},
			4,
		},
	}

	for _, td := range tt {
		t.Run(
			td.name,
			func(t *testing.T) {
				q, r, err := mpi.ShortDiv(td.dividend, td.divisor)
				if err != nil {
					t.Fatalf("non-nill error in mpi.Mul: %s", err)
				}
				if !reflect.DeepEqual(q.V, td.quotent.V) || r != td.remainder {
					t.Fatalf("Wrong result:\nExpected\t%d\t|\t%d\nGot\t\t%d\t|\t%d", td.remainder, td.quotent.V, r, q.V)
				}
			},
		)
	}
}

func TestLongDiv(t *testing.T) {
	tt := []struct {
		name      string
		dividend  mpi.MPI
		divisor   mpi.MPI
		quotent   mpi.MPI
		remainder mpi.MPI
	}{
		{
			"115923 / 344 = 336, 339",
			mpi.MPI{B: 10, V: []int64{3, 2, 9, 5, 1, 1}},
			mpi.MPI{B: 10, V: []int64{4, 4, 3}},
			mpi.MPI{B: 10, V: []int64{6, 3, 3, 0}},
			mpi.MPI{B: 10, V: []int64{9, 3, 3}},
		},
		{
			"1956 / 55 = 35, 31",
			mpi.MPI{B: 10, V: []int64{6, 5, 9, 1}},
			mpi.MPI{B: 10, V: []int64{5, 5}},
			mpi.MPI{B: 10, V: []int64{5, 3, 0}},
			mpi.MPI{B: 10, V: []int64{1, 3}},
		},
	}

	for _, td := range tt {
		t.Run(
			td.name,
			func(t *testing.T) {
				q, r, err := mpi.LongDiv(td.dividend, td.divisor)
				if err != nil {
					t.Fatalf("non-nill error in mpi.LongDiv: %s", err)
				}
				if !reflect.DeepEqual(q.V, td.quotent.V) || !reflect.DeepEqual(r.V, td.remainder.V) {
					t.Fatalf("\nWrong result:\tR\t\t\tQ\nExpected\t\t%d\t|\t%d\nGot\t\t\t\t%d\t|\t%d", td.remainder.V, td.quotent.V, r.V, q.V)
				}
			},
		)
	}
}

func TestGCD(t *testing.T) {
	tt := []struct {
		name string
		a    mpi.MPI
		b    mpi.MPI
		x    mpi.MPI
		y    mpi.MPI
		d    mpi.MPI
	}{
		{
			"gcd(35, 20) = 5, -1, 2",
			mpi.MPI{B: 10, V: []int64{5, 3}},
			mpi.MPI{B: 10, V: []int64{0, 2}},
			mpi.MPI{B: 10, V: []int64{1}, N: true},
			mpi.MPI{B: 10, V: []int64{2}, N: false},
			mpi.MPI{B: 10, V: []int64{5}},
		},
		{
			"gcd(45, 127) = 1, 48, -17",
			mpi.MPI{B: 10, V: []int64{7, 2, 1}},
			mpi.MPI{B: 10, V: []int64{5, 4}},
			mpi.MPI{B: 10, V: []int64{7, 1}, N: true},
			mpi.MPI{B: 10, V: []int64{8, 4}, N: false},
			mpi.MPI{B: 10, V: []int64{1}},
		},
	}

	for _, td := range tt {
		t.Run(
			td.name,
			func(t *testing.T) {
				x, y, d, err := mpi.ExtendedGCD(td.a, td.b)
				if err != nil {
					t.Fatalf("non-nill error in mpi.ExtendedGCD: %s", err)
				}
				if !reflect.DeepEqual(x, td.x) || !reflect.DeepEqual(y, td.y) || !reflect.DeepEqual(td.d, d) {
					t.Fatalf("\nWrong result:\tR\t\t\tQ\nExpected\t\t%v\t|\t%v\t|\t%v\nGot\t\t\t\t%v\t|\t%v\t|\t%v", td.x, td.y, td.d, x, y, d)
				}
			},
		)
	}
}

func TestFromZnBaseB(t *testing.T) {
	tt := []struct {
		n uint64
		b uint32
		N mpi.MPI
	}{
		{
			n: 45,
			b: 8,
			N: mpi.MPI{B: 8, V: []int64{5, 5}, N: false},
		},
		{
			n: 45,
			b: 10,
			N: mpi.MPI{B: 10, V: []int64{5, 4}},
		},
	}

	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%d to base %d", td.n, td.b),
			func(t *testing.T) {
				N := mpi.FromZnBaseB(td.n, td.b)
				if !reflect.DeepEqual(N, td.N) {
					t.Fatalf("\nWrong result:\nExpected\t%v\nGot\t\t\t%v", td.N, N)
				}
			},
		)
		t.Run(
			fmt.Sprintf("%d from %v", td.n, td.N),
			func(t *testing.T) {
				n := td.N.ToZn()
				if td.n != n {
					t.Fatalf("\nWrong result:\nExpected\t%v\nGot\t\t\t%v", td.n, n)
				}
			},
		)
	}
}
