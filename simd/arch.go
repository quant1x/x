package simd

// 波尔类型
type boolean interface {
	~bool
}

// 整型
type integer interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// 浮点类型
type float interface {
	~float32 | ~float64
}

// 数字类型
type number interface {
	integer | float
}

// 加速器接口
type accelerator[E number] interface {
	// 加
	add(x, y []E) []E
	// 减
	sub(x, y []E) []E
	// 乘
	mul(x, y []E) []E
	// 除
	div(x, y []E) []E
	// 取模
	mod(x, y []E) []E
}

type f32s []float32
