package montgomery_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-1_montgomery/montgomery"
	"github.com/ChainsAre2Tight/smpvzib-labs/lab-1_montgomery/mpi"
)

func TestNewMontgomeryRdashed(t *testing.T) {
	tt := []struct {
		name    string
		N       mpi.MPI
		Rdashed mpi.MPI
	}{
		{
			name:    "2011, b=2^3 = 8, R=4096, R' = 1454",
			N:       mpi.MPI{B: 8, V: []int64{3, 3, 7, 3}},
			Rdashed: mpi.MPI{B: 8, V: []int64{6, 5, 6, 2}},
		},
	}

	for _, td := range tt {
		t.Run(
			td.name,
			func(t *testing.T) {
				m, err := montgomery.NewMontgomery(td.N)
				if err != nil {
					t.Fatalf("non-nill error in montgomery.NewMontgomery: %s", err)
				}
				if !reflect.DeepEqual(m.Rdashed, td.Rdashed) {
					t.Fatalf("Wrong result:\nExpected\t%v\nGot\t\t%v", td.Rdashed, m.Rdashed)
				}
			},
		)
	}
}

func TestNewMontgomeryCoprimeCheck(t *testing.T) {
	tt := []struct {
		N        uint64
		radix    uint32
		expected bool
	}{
		{2011, 10, true},
		{2010, 10, false},
	}

	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%d base %d - %v", td.N, td.radix, td.expected),
			func(t *testing.T) {
				_, err := montgomery.NewMontgomery(mpi.FromZnBaseB(td.N, td.radix))
				if td.expected != (err == nil) || !td.expected && err != nil && err.Error() != "NewMontgomery: bad modulo: N is not coprime with radix" {
					t.Fatalf("\nWrong state:\nExpected error:%v\nGot:%s", td.expected, err)
				}
			},
		)
	}
}

func TestRedc(t *testing.T) {
	tt := []struct {
		name string
		N    mpi.MPI
		b    uint32
		u    mpi.MPI
		res  mpi.MPI
	}{
		{
			name: "2011, b=2^3 = 8, u = 8170821",
			N:    mpi.MPI{B: 8, V: []int64{3, 3, 7, 3}},
			b:    8,
			u:    mpi.MPI{B: 8, V: []int64{5, 0, 5, 6, 2, 1, 7, 3}},
			res:  mpi.MPI{B: 8, V: []int64{4, 1, 4, 5}},
		},
	}

	for _, td := range tt {
		t.Run(
			td.name,
			func(t *testing.T) {
				m, err := montgomery.NewMontgomery(td.N)
				if err != nil {
					t.Fatalf("non-nill error in montgomery.NewMontgomery: %s", err)
				}
				res, err := m.Redc(td.u)
				if err != nil {
					t.Fatalf("non-nill error in montgomery.Redc: %s", err)
				}
				if !reflect.DeepEqual(res.V, td.res.V) {
					t.Fatalf("Wrong result:\nExpected\t%d\nGot\t\t%d", td.res.V, res.V)
				}
			},
		)
	}
}

func TestPhi(t *testing.T) {
	N := mpi.MPI{B: 8, V: []int64{3, 3, 7, 3}}

	tt := []struct {
		name string
		u    mpi.MPI
		uR   mpi.MPI
	}{
		{
			name: "[45] = (2447)_8",
			u:    mpi.MPI{B: 8, V: []int64{5, 5}},
			uR:   mpi.MPI{B: 8, V: []int64{7, 4, 4, 2}},
		}, {
			name: "[97] = (2447)_8",
			u:    mpi.MPI{B: 8, V: []int64{1, 4, 1}},
			uR:   mpi.MPI{B: 8, V: []int64{1, 7, 1, 2}},
		},
	}

	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s Phi", td.name),
			func(t *testing.T) {
				m, err := montgomery.NewMontgomery(N)
				if err != nil {
					t.Fatalf("non-nill error in montgomery.NewMontgomery: %s", err)
				}
				uR, err := m.Phi(td.u)
				if err != nil {
					t.Fatalf("non-nill error in montgomery.Phi: %s", err)
				}
				if !reflect.DeepEqual(uR, td.uR) {
					t.Fatalf("Wrong result:\nExpected\t%v\nGot\t\t%v", td.uR, uR)
				}
			},
		)
		t.Run(
			fmt.Sprintf("%s Inverse Phi", td.name),
			func(t *testing.T) {
				m, err := montgomery.NewMontgomery(N)
				if err != nil {
					t.Fatalf("non-nill error in montgomery.NewMontgomery: %s", err)
				}
				u, err := m.InversePhi(td.uR)
				if err != nil {
					t.Fatalf("non-nill error in montgomery.InversePhi: %s", err)
				}
				if !reflect.DeepEqual(u, td.u) {
					t.Fatalf("Wrong result:\nExpected\t%v\nGot\t\t%v", td.u, u)
				}
			},
		)
	}
}

func TestMontMul(t *testing.T) {
	tt := []struct {
		name string
		N    mpi.MPI
		u    mpi.MPI
		v    mpi.MPI
		res  mpi.MPI
	}{
		{
			name: "2011, b=2^3 = 8, u = 8170821",
			N:    mpi.MPI{B: 8, V: []int64{3, 3, 7, 3}},
			u:    mpi.MPI{B: 8, V: []int64{7, 4, 4, 2}},
			v:    mpi.MPI{B: 8, V: []int64{1, 7, 1, 2}},
			res:  mpi.MPI{B: 8, V: []int64{2, 4, 3, 2}},
		},
	}

	for _, td := range tt {
		t.Run(
			td.name,
			func(t *testing.T) {
				m, err := montgomery.NewMontgomery(td.N)
				if err != nil {
					t.Fatalf("non-nill error in montgomery.NewMontgomery: %s", err)
				}
				res, err := m.MontMul(td.u, td.v)
				if err != nil {
					t.Fatalf("non-nill error in montgomery.MontMul: %s", err)
				}
				if !reflect.DeepEqual(res, td.res) {
					t.Fatalf("Wrong result:\nExpected\t%v\nGot\t\t%v", td.res, res)
				}
			},
		)
	}
}

func TestZnMul(t *testing.T) {
	tt := []struct {
		name string
		N    mpi.MPI
		u    mpi.MPI
		v    mpi.MPI
		res  mpi.MPI
	}{
		{
			name: "2011, b=2^3 = 8, u = 8170821",
			N:    mpi.MPI{B: 8, V: []int64{3, 3, 7, 3}},
			u:    mpi.MPI{B: 8, V: []int64{5, 5}},
			v:    mpi.MPI{B: 8, V: []int64{1, 4, 1}},
			res:  mpi.MPI{B: 8, V: []int64{7, 2, 5}},
		},
		{
			name: "2011, b=2^3 = 10, u = 8170821",
			N:    mpi.MPI{B: 10, V: []int64{1, 1, 0, 2}},
			u:    mpi.MPI{B: 10, V: []int64{5, 4}},
			v:    mpi.MPI{B: 10, V: []int64{7, 9}},
			res:  mpi.MPI{B: 10, V: []int64{3, 4, 3}},
		},
	}

	for _, td := range tt {
		t.Run(
			td.name,
			func(t *testing.T) {
				m, err := montgomery.NewMontgomery(td.N)
				if err != nil {
					t.Fatalf("non-nill error in montgomery.NewMontgomery: %s", err)
				}
				res, err := m.ZnMul(td.u, td.v)
				if err != nil {
					t.Fatalf("non-nill error in montgomery.ZnMul: %s", err)
				}
				if !reflect.DeepEqual(res, td.res) {
					t.Fatalf("Wrong result:\nExpected\t%v\nGot\t\t%v", td.res, res)
				}
			},
		)
	}
}

func TestZnExp(t *testing.T) {
	tt := []struct {
		name string
		n    uint64
		a    uint64
		exp  uint64
		res  uint64
	}{
		{
			name: "1",
			n:    77,
			a:    5,
			exp:  3,
			res:  48,
		},
	}

	for _, td := range tt {
		t.Run(
			td.name,
			func(t *testing.T) {
				m, err := montgomery.NewMontgomery(mpi.FromZnBaseB(td.n, 10))
				if err != nil {
					t.Fatalf("non-nill error in montgomery.NewMontgomery: %s", err)
				}
				res, err := m.ZnExp(mpi.FromZnBaseB(td.a, 10), td.exp)
				if err != nil {
					t.Fatalf("non-nill error in montgomery.ZnExp: %s", err)
				}
				if res.ToZn() != td.res {
					t.Fatalf("Wrong result:\nExpected\t%v\nGot\t\t%v", td.res, res.ToZn())
				}
			},
		)
	}
}

func FuzzMultiplication(f *testing.F) {
	f.Add(uint64(45), uint64(97), uint64(2011))
	f.Add(uint64(5), uint64(10), uint64(31))
	f.Add(uint64(40), uint64(2), uint64(51))
	f.Fuzz(func(t *testing.T, a, b, N uint64) {
		expected := (a * b) % N
		n := mpi.FromZnBaseB(N, 10)
		mont, err := montgomery.NewMontgomery(n)
		if err != nil {
			t.Fatalf("NewMont: %s", err)
		}
		a1 := mpi.FromZnBaseB(a, 10)
		b1 := mpi.FromZnBaseB(b, 10)
		got, err := mont.ZnMul(a1, b1)
		if err != nil {
			t.Fatalf("ZnMul: %s", err)
		}
		res := got.ToZn()
		if res != expected {
			t.Fatalf("\nExpected\t%d\nGot\t\t%d", expected, res)
		}
	})
}
