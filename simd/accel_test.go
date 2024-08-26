package simd

import (
	"math"
	"math/rand"
	"sync"
)

const (
	benchAlignLength  = 5000
	benchAlignInitNum = 1000
)

var (
	testalignOnce    sync.Once
	testDataBoolx    []bool
	testDataBooly    []bool
	testDataBoolr    []bool
	testDataFloat32x []float32
	testDataFloat32y []float32
	testDataFloat32r []float32
	testDataFloat64x []float64
	testDataFloat64y []float64
	testDataFloat64r []float64
)

func randFloat32() float32 {
	return rand.Float32()*(math.MaxFloat32-math.SmallestNonzeroFloat32) + math.SmallestNonzeroFloat32
}

func randFloat64() float64 {
	return rand.Float64()*(math.MaxFloat64-math.SmallestNonzeroFloat64) + math.SmallestNonzeroFloat64
}

func initTestData() {
	testDataBoolx = make([]bool, benchAlignInitNum)
	testDataBooly = make([]bool, benchAlignInitNum)
	testDataBoolr = make([]bool, benchAlignInitNum)
	testDataFloat32x = make([]float32, benchAlignInitNum)
	testDataFloat32y = make([]float32, benchAlignInitNum)
	testDataFloat32r = make([]float32, benchAlignInitNum)
	testDataFloat64x = make([]float64, benchAlignInitNum)
	testDataFloat64y = make([]float64, benchAlignInitNum)
	testDataFloat64r = make([]float64, benchAlignInitNum)
	for i := 0; i < benchAlignInitNum; i++ {
		testDataBoolx[i] = i%8 == 0
		testDataBooly[i] = i%16 == 0
		testDataBoolr[i] = testDataBoolx[i] && testDataBooly[i]
		testDataFloat32x[i] = randFloat32()
		testDataFloat32y[i] = randFloat32()
		testDataFloat32r[i] = testDataFloat32x[i] + testDataFloat32y[i]
		testDataFloat64x[i] = randFloat64()
		testDataFloat64y[i] = randFloat64()
		testDataFloat64r[i] = testDataFloat64x[i] + testDataFloat64y[i]
	}
}

func init() {
	initTestData()
}
