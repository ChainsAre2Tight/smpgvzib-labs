package interfaces

type PrimalityTest interface {
	IsPrime(N uint64, rounds int) (res bool, err error)
}

type PrimalityTestFunc func(N uint64, rounds int) (res bool, err error)

func (f PrimalityTestFunc) IsPrime(N uint64, rounds int) (res bool, err error) {
	return f(N, rounds)
}
