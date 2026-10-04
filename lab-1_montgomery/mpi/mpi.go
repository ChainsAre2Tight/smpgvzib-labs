package mpi

import (
	"fmt"
	"math"
)

type MPI struct {
	B uint32
	V []int64
	N bool
}

func FromInt64(val int64, length int, radix uint32) MPI {
	res := MPI{
		B: uint32(radix),
		V: make([]int64, length),
		N: false,
	}

	res.V[0] = val
	return res
}

func FromZnBaseB(a uint64, radix uint32) (N MPI) {
	N = MPI{
		B: radix,
		N: false,
		V: make([]int64, 0, 4),
	}

	n := int64(a)
	for n > 0 {
		N.V = append(N.V, n%int64(radix))
		n = n / int64(radix)
	}
	return N
}

func (m *MPI) ToZn() (a uint64) {
	for i := range m.V {
		exp := uint64(math.Pow(float64(uint64(m.B)), float64(uint64(i))))
		a += uint64(m.V[i]) * exp
	}
	return a
}

func Add(u, v MPI) (w MPI, err error) {
	if u.B != v.B {
		return MPI{}, fmt.Errorf("MPI.add: radix mismatch")
	}
	b := u.B

	if len(u.V) != len(v.V) {
		return MPI{}, fmt.Errorf("MPI.add length mismatch")
	}
	n := len(u.V)

	w = MPI{
		B: b,
		V: make([]int64, n),
		N: false,
	}

	var k int64 = 0
	for i := range n {
		temp := u.V[i] + v.V[i] + k
		w.V[i] = temp % int64(b)
		k = temp / int64(b)
	}

	w.V[n-1] += k
	return w, nil
}

func Sub(u, v MPI) (w MPI, undeflow bool, err error) {
	if u.B != v.B {
		return MPI{}, false, fmt.Errorf("MPI.sub: radix mismatch")
	}
	b := u.B

	if len(u.V) != len(v.V) {
		return MPI{}, false, fmt.Errorf("MPI.sub: length mismatch (U: %d, V: %d)", len(u.V), len(v.V))
	}
	n := len(u.V)

	w = MPI{
		B: b,
		V: make([]int64, n),
	}

	var k int64 = 0
	for i := range n {
		temp := int64(u.V[i]-v.V[i]) + k
		if temp >= 0 {
			w.V[i] = temp % int64(b)
			k = temp / int64(b)
		} else {
			r := (int64(b) + temp) % int64(b)
			w.V[i] = r
			k = temp / int64(b)
			if r != 0 {
				k--
			}
		}
	}

	if k < 0 {

		big := MPI{B: b, V: make([]int64, n+1)}
		big.V[n] = 1
		w.V = append(w.V, 0)

		w, _, err = Sub(big, w)
		if err != nil {
			return MPI{}, false, fmt.Errorf("MPI.Sub: error during invertion: %s", err)
		}
		w.N = true

		w.V = w.V[:len(w.V)-1]
	}

	return w, k < 0, nil
}

func Mul(u, v MPI) (w MPI, err error) {
	if u.B != v.B {
		return MPI{}, fmt.Errorf("MPI.mul: radix mismatch")
	}
	b := u.B
	m := len(u.V)
	n := len(v.V)
	w = MPI{
		B: b,
		V: make([]int64, n+m),
		N: u.N != v.N,
	}

	for i := range n {
		var k int64 = 0
		if v.V[i] == 0 {
			w.V[m+i] = 0
			continue
		}

		for j := range m {
			t := v.V[i]*u.V[j] + w.V[i+j] + k
			w.V[i+j] = t % int64(b)
			k = t / int64(b)
		}
		w.V[m+i] = k
	}

	return w, nil
}

func Slice(a MPI, start, end int) (r MPI, err error) {
	if start > end || end > len(a.V) || start < 0 {
		return MPI{}, fmt.Errorf("mpi.Slice: bad start/end")
	}

	r = MPI{
		B: a.B,
		V: make([]int64, end-start+1),
		N: false,
	}

	copy(r.V, a.V[start:end+1])

	return r, nil
}

// modifies a[start:start+len(b.V)]
func Insert(a, b MPI, start int) (r MPI, err error) {
	if start < 0 || len(b.V)+start > len(a.V) {
		return MPI{}, fmt.Errorf("mpi.slice: bad start/length")
	}
	if b.B != a.B {
		return MPI{}, fmt.Errorf("mpi.insert: radix mismatch")
	}

	r = MPI{
		B: a.B,
		V: make([]int64, len(a.V)),
		N: false,
	}

	copy(r.V, a.V)
	copy(r.V[start:start+len(b.V)], b.V)

	return r, nil
}

func ShortDiv(u MPI, v int64) (q MPI, r int64, err error) {
	r = 0
	n := len(u.V)

	if v == 0 {
		return MPI{}, 0, fmt.Errorf("MPI.ShortDiv: DivisionByZero")
	}

	q = MPI{
		B: u.B,
		V: make([]int64, n),
		N: false,
	}
	b := u.B

	for i := n - 1; i >= 0; i-- {
		t := r*int64(b) + u.V[i]
		q.V[i] = t / v
		r = t % v
	}

	return q, r, nil
}

func (a MPI) Trim() (b MPI) {

	length := len(a.V)
	for i := len(a.V) - 1; i >= 1; i-- {
		if a.V[i] != 0 {
			break
		}
		length--
	}

	b = MPI{
		B: a.B,
		V: make([]int64, length),
		N: a.N,
	}

	copy(b.V, a.V[:length])
	return b
}

func EqualizeLength(a, b MPI) (a1, b1 MPI) {
	mx := max(len(a.V), len(b.V))
	a1 = MPI{
		B: a.B,
		N: a.N,
		V: make([]int64, mx),
	}
	b1 = MPI{
		B: b.B,
		N: b.N,
		V: make([]int64, mx),
	}
	copy(a1.V, a.V)
	copy(b1.V, b.V)

	return a1, b1
}

func LongDiv(u MPI, v MPI) (q MPI, r MPI, err error) {

	fail := func(err error) (MPI, MPI, error) {
		return MPI{}, MPI{}, fmt.Errorf("MPI.LongDiv: %s", err)
	}

	if u.B != v.B {
		return fail(fmt.Errorf("radix mismatch"))
	}

	n := len(v.V)
	if n == 1 {
		q, r1, err := ShortDiv(u, v.V[0])
		if err != nil {
			return fail(fmt.Errorf("ShortDiv: %s", err))
		}
		r = MPI{
			B: u.B,
			V: []int64{r1},
			N: false,
		}
		return q, r, nil
	}
	if n < 2 || v.V[n-1] == 0 {
		return fail(fmt.Errorf("divisionbyzero"))
	}
	m := len(u.V) - n
	b := u.B

	// u.V[n+m] = 0
	u.V = append(u.V, 0)
	d := 1

	for v.V[n-1] < int64(b/2) {
		v, err = Mul(v, FromInt64(2, 1, b))
		if err != nil {
			return fail(fmt.Errorf("V *= 2: %s", err))
		}
		if len(v.V) > n {
			v.V = v.V[:n]
		}
		u, err = Mul(u, FromInt64(2, 1, b))
		// if len(u.V) > m+n {
		// 	u.V = u.V[:n+m]
		// }
		if err != nil {
			return fail(fmt.Errorf("U *= 2: %s", err))
		}
		d *= 2
	}

	q = MPI{
		B: b,
		V: make([]int64, m+1),
		N: false,
	}
	r = MPI{
		B: b,
		V: make([]int64, n),
		N: false,
	}

	for i := m; i >= 0; i-- {
		line1 := u.V[i+n]*int64(b) + u.V[i+n-1] // panic
		q1 := min(int64(b-1), line1/v.V[n-1])
		for {
			debug2 := (u.V[i+n]*int64(b*b) + u.V[i+n-1]*int64(b) + u.V[i+n-2])
			debug1 := q1 * (v.V[n-1]*int64(b) + v.V[n-2])
			if debug1 <= debug2 {
				break
			}
			q1 -= 1
		}
		t1, err := Mul(v, FromInt64(q1, 1, b))
		if err != nil {
			return fail(fmt.Errorf("q1*v: %s", err))
		}
		t2, err := Slice(u, i, i+n)
		if err != nil {
			return fail(fmt.Errorf("slice u: %s", err))
		}
		// for len(t1.V) < len(t2.V) {
		// 	t1.V = append(t1.V, 0)
		// }
		t3, undeflow, err := Sub(t2, t1)
		if err != nil {
			return fail(fmt.Errorf("u-q1*b: %s", err))
		}
		if undeflow {
			q1--
			t4 := v
			for len(t4.V) < len(t3.V) {
				t4.V = append(t4.V, 0)
			}
			t3, _, err = Sub(t4, t3)
			if err != nil {
				return fail(fmt.Errorf("+0v: %s", err))
			}
		}
		u, err = Insert(u, t3, i)
		if err != nil {
			return fail(fmt.Errorf("u[i:i+n] = t3: %s", err))
		}
		q.V[i] = q1
	}

	r, _, err = ShortDiv(u, int64(d))
	if err != nil {
		return fail(fmt.Errorf("u/d: %s", err))
	}

	q.V = q.V[:m+1]
	r.V = r.V[:n]

	return q, r, nil
}

func ExtendedGCD(a, b MPI) (x, y, d MPI, err error) {
	fail := func(err error) (MPI, MPI, MPI, error) {
		return MPI{}, MPI{}, MPI{}, fmt.Errorf("mpi.ExtendedGCD: %s", err)
	}

	if a.B != b.B {
		return fail(fmt.Errorf("radix mismatch"))
	}

	// if b == 0
	b_flag := false
	for _, el := range b.V {
		if el != 0 {
			b_flag = true
		}
	}

	if !b_flag {
		return FromInt64(1, 1, a.B), FromInt64(0, 1, a.B), a, nil
	}

	adivb, amodb, err := LongDiv(a, b)
	if err != nil {
		return fail(fmt.Errorf("a/b, a mod b: %s", err))
	}
	adivb = adivb.Trim()
	amodb = amodb.Trim()
	x1, y1, d1, err := ExtendedGCD(b, amodb)
	if err != nil {
		return fail(err)
	}

	t1, err := Mul(adivb, y1)
	if err != nil {
		return fail(fmt.Errorf("a/b * y: %s", err))
	}

	var t2 MPI
	t1, x1 = EqualizeLength(t1, x1)
	if !t1.N && !x1.N {
		t2, _, err = Sub(x1, t1)
		if err != nil {
			return fail(fmt.Errorf("+x - +y: %s", err))
		}
	} else if !x1.N && t1.N {
		t2, err = Add(x1, t1)
		if err != nil {
			return fail(fmt.Errorf("+x + -y: %s", err))
		}
	} else {
		t2, err = Add(x1, t1)
		if err != nil {
			return fail(fmt.Errorf("-(-x + -y): %s", err))
		}
		t2.N = true
	}
	y1 = y1.Trim()
	t2 = t2.Trim()
	d1 = d1.Trim()

	// fmt.Println(a, b, amodb, adivb, x1, y1, y1, t2)
	return y1, t2, d1, nil
}
