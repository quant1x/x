package assert

type assertType interface {
	~bool | ~float32 | ~float64
}

func Equal[E assertType](x, y []E) bool {
	if len(x) != len(y) {
		return false
	}
	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}
