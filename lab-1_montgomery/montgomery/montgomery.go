package montgomery

import (
	"fmt"
	"reflect"

	"github.com/ChainsAre2Tight/smpvzib-labs/lab-1_montgomery/mpi"
)

func compute_N_dashed(N, R mpi.MPI) (N1 mpi.MPI, err error) {
	_, x, _, err := mpi.ExtendedGCD(R, N)
	if err != nil {
		return mpi.MPI{}, fmt.Errorf("reduction.compute_n_dashed: ExtendedGCD: %s", err)
	}
	R, x = mpi.EqualizeLength(R, x)
	if !x.N {
		N1, _, err = mpi.Sub(R, x)
		if err != nil {
			return mpi.MPI{}, fmt.Errorf("reduction.compute_n_dashed: sub(R, x): %s", err)
		}
		if N1.N {
			panic("reduction.compute_n_dashed: sub(R, x) undeflow")
		}
	} else {
		N1 = mpi.MPI{
			B: x.B,
			V: make([]int64, len(x.V)),
			N: false,
		}
		copy(N1.V, x.V)
	}
	return N1, nil
}

type Montgomery struct {
	N       mpi.MPI
	R       mpi.MPI
	Rdashed mpi.MPI
	Ndashed mpi.MPI
	radix   uint32
}

func NewMontgomery(N mpi.MPI) (*Montgomery, error) {

	_, _, d, err := mpi.ExtendedGCD(N, mpi.FromZnBaseB(uint64(N.B), N.B))
	if err != nil {
		return nil, fmt.Errorf("NewMontgomery: ExtendedGCD: %s", err)
	}
	if !reflect.DeepEqual(d.Trim().V, []int64{1}) {
		return nil, fmt.Errorf("NewMontgomery: bad modulo: N is not coprime with radix")
	}

	R := mpi.MPI{
		B: N.B,
		V: make([]int64, len(N.V)+1),
		N: false,
	}
	R.V[len(N.V)] = 1
	N_dashed, err := compute_N_dashed(N, R)
	N_dashed = N_dashed.Trim()
	if err != nil {
		return nil, fmt.Errorf("NewMontgomery: %s", err)
	}

	R2, err := mpi.Mul(R, R)
	if err != nil {
		return nil, fmt.Errorf("NewMontgomery: R^2: %s", err)
	}

	_, Rdashed, err := mpi.LongDiv(R2, N)
	if err != nil {
		return nil, fmt.Errorf("NewMontgomery: R^2 mod N: %s", err)
	}

	return &Montgomery{
		N:       N,
		R:       R,
		Ndashed: N_dashed,
		Rdashed: Rdashed,
		radix:   N.B,
	}, nil
}

func (m *Montgomery) Redc(u mpi.MPI) (t mpi.MPI, err error) {
	if uint32(m.radix) != u.B {
		return mpi.MPI{}, fmt.Errorf("reduction.Redc: radix mismatch")
	}
	b := m.radix
	n := len(m.N.V)
	if len(u.V) > int(2*n) {
		return mpi.MPI{}, fmt.Errorf("reduction.Redc: expected u to be not larger than 2n in length")
	}

	t = mpi.MPI{
		B: u.B,
		V: make([]int64, 2*n),
	}

	copy(t.V, u.V)

	for i := range n {
		k1, err := mpi.Mul(m.Ndashed, mpi.FromInt64(t.V[i], 1, b))
		if err != nil {
			return mpi.MPI{}, fmt.Errorf("reduction.Redc: ti*N': %s", err)
		}
		k := k1.V[0]
		kibi := mpi.MPI{
			B: b,
			N: false,
			V: make([]int64, i+1),
		}
		kibi.V[i] = k
		kibiN, err := mpi.Mul(m.N, kibi)
		if err != nil {
			return mpi.MPI{}, fmt.Errorf("reduction.Redc: knbi: %s", err)
		}
		kibiN, t = mpi.EqualizeLength(kibiN, t)
		t, err = mpi.Add(t, kibiN)
		if err != nil {
			return mpi.MPI{}, fmt.Errorf("reduction.Redc: t+knbi: %s", err)
		}
	}

	t.V = t.V[n:]

	return t, nil
}

func (m *Montgomery) MontMul(u, v mpi.MPI) (t mpi.MPI, err error) {
	if uint32(m.radix) != u.B {
		return mpi.MPI{}, fmt.Errorf("montgomery.Mul: radix mismatch")
	}
	b := m.radix
	n := len(m.N.V)
	if len(u.V) > int(n) {
		return mpi.MPI{}, fmt.Errorf("montgomery.Mul: expected u to be not larger than n in length")
	}
	if len(v.V) > int(n) {
		return mpi.MPI{}, fmt.Errorf("montgomery.Mul: expected v to be not larger than n in length")
	}
	u, _ = mpi.EqualizeLength(u, m.N)
	v, _ = mpi.EqualizeLength(v, m.N)

	t = mpi.MPI{
		B: u.B,
		V: make([]int64, n),
		N: false,
	}

	for i := range n {
		t1 := t.V[0] + u.V[i]*v.V[0]
		t2, err := mpi.Mul(m.Ndashed, mpi.FromInt64(t1, 1, b))
		if err != nil {
			return mpi.MPI{}, fmt.Errorf("montgomery.Mul: (t0+uiv0)*p': %s", err)
		}
		uiv, err := mpi.Mul(v, mpi.FromInt64(u.V[i], 1, b))
		if err != nil {
			return mpi.MPI{}, fmt.Errorf("montgomery.Mul: u_i*v: %s", err)
		}
		mip, err := mpi.Mul(m.N, mpi.FromInt64(t2.V[0], 1, b))
		t4, err := mpi.Add(mip, uiv)
		if err != nil {
			return mpi.MPI{}, fmt.Errorf("montgomery.Mul: mip+uiv': %s", err)
		}
		t5, t4 := mpi.EqualizeLength(t, t4)
		t6, err := mpi.Add(t5, t4)
		if err != nil {
			return mpi.MPI{}, fmt.Errorf("montgomery.Mul: t+uiv+mip: %s", err)
		}
		t, err = mpi.Slice(t6, 1, n)
		if err != nil {
			return mpi.MPI{}, fmt.Errorf("montgomery.Mul: t = (t + uiv + mip) / b': %s", err)
		}
	}

	return t, nil
}

func (m *Montgomery) Phi(x mpi.MPI) (xR mpi.MPI, err error) {
	xRdashed, err := mpi.Mul(x, m.Rdashed)
	if err != nil {
		return mpi.MPI{}, fmt.Errorf("montgomery.Phi: xRdashed: %s", err)
	}
	xR, err = m.Redc(xRdashed)
	if err != nil {
		return mpi.MPI{}, fmt.Errorf("montgomery.Phi: Redc(xRdashed): %s", err)
	}
	return xR, nil
}

func (m *Montgomery) InversePhi(xR mpi.MPI) (x mpi.MPI, err error) {
	x, err = m.Redc(xR)
	if err != nil {
		return mpi.MPI{}, fmt.Errorf("montgomery.InversePhi: Redc(xR): %s", err)
	}
	return x.Trim(), nil
}

func (m *Montgomery) ZnMul(u, v mpi.MPI) (t mpi.MPI, err error) {
	uR, err := m.Phi(u)
	if err != nil {
		return mpi.MPI{}, fmt.Errorf("montgomery.ZnMul: Phi(u): %s", err)
	}
	vR, err := m.Phi(v)
	if err != nil {
		return mpi.MPI{}, fmt.Errorf("montgomery.ZnMul: Phi(v): %s", err)
	}
	tR, err := m.MontMul(uR, vR)
	if err != nil {
		return mpi.MPI{}, fmt.Errorf("montgomery.ZnMul: MontMul(uR, vR): %s", err)
	}
	t, err = m.InversePhi(tR)
	if err != nil {
		return mpi.MPI{}, fmt.Errorf("montgomery.ZnMul: InversePhi(tR): %s", err)
	}
	return t, nil
}

func (m *Montgomery) MontExp(u mpi.MPI, exp uint64) (t mpi.MPI, err error) {
	if uint32(m.radix) != u.B {
		return mpi.MPI{}, fmt.Errorf("montgomery.Mul: radix mismatch")
	}
	b := m.radix
	if exp == 0 {
		return mpi.FromInt64(1, 1, b), nil
	}
	t = u
	for range exp - 1 {
		t, err = m.MontMul(t, u)
		if err != nil {
			return mpi.MPI{}, fmt.Errorf("exp, mul: %s", err)
		}
		t = t.Trim()
	}

	return t, nil
}

func (m *Montgomery) ZnExp(u mpi.MPI, exp uint64) (t mpi.MPI, err error) {
	uR, err := m.Phi(u)
	if err != nil {
		return mpi.MPI{}, fmt.Errorf("montgomery.ZnExp: Phi(u): %s", err)
	}
	tR, err := m.MontExp(uR, exp)
	if err != nil {
		return mpi.MPI{}, fmt.Errorf("montgomery.ZnExp: MontExp(uR, exp): %s", err)
	}
	t, err = m.InversePhi(tR)
	if err != nil {
		return mpi.MPI{}, fmt.Errorf("montgomery.ZnExp: InversePhi(tR): %s", err)
	}
	return t, nil
}
