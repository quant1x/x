package simd

import (
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
	testDataFloat32  []float32
	testDataFloat32y []float32
	testDataFloat64  []float64
	testDataFloat64y []float64
)

func initTestData() {
	testDataBoolx = make([]bool, benchAlignInitNum)
	testDataBooly = make([]bool, benchAlignInitNum)
	testDataBoolr = make([]bool, benchAlignInitNum)
	testDataFloat32 = make([]float32, benchAlignInitNum)
	testDataFloat32y = make([]float32, benchAlignInitNum)
	testDataFloat64 = make([]float64, benchAlignInitNum)
	testDataFloat64y = make([]float64, benchAlignInitNum)
	for i := 0; i < benchAlignInitNum; i++ {
		testDataBoolx[i] = i%8 == 0
		testDataBooly[i] = i%16 == 0
		testDataBoolr[i] = testDataBoolx[i] && testDataBooly[i]
		testDataFloat32[i] = rand.Float32()
		testDataFloat32y[i] = rand.Float32()
		testDataFloat64[i] = rand.Float64()
		testDataFloat64y[i] = rand.Float64()
	}
}

func init() {
	initTestData()
}
