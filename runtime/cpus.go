package runtime

import "runtime"

// GoMaxProcs 设置CPU核数
func GoMaxProcs(n ...int) {
	numCPU := runtime.NumCPU() / 2
	if len(n) > 0 && n[0] > 0 {
		numCPU = n[0]
	}
	runtime.GOMAXPROCS(numCPU)
}
