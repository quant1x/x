package simd

import "slices"

func bs_and(a, b []bool) []bool {
	return v1_bs_and(a, b)
}

func v1_bs_and(a, b []bool) []bool {
	an := len(a)
	bn := len(b)
	if an != bn {
		panic("Add: bad len")
	}
	length := an
	result := make([]bool, length)
	n := b8x32_and(a, b, result)

	for i := length - n; i < length; i++ {
		result[i] = a[i] && b[i]
	}
	return result
}

func v2_bs_and(a, b []bool) []bool {
	an := len(a)
	bn := len(b)
	if an != bn {
		panic("Add: bad len")
	}
	length := an
	result := make([]bool, length)
	n := b32x8_and(a, b, result)

	for i := length - n; i < length; i++ {
		result[i] = a[i] && b[i]
	}
	return result
}

// for vek
func v3_bs_and(a, b []bool) []bool {
	an := len(a)
	bn := len(b)
	if an != bn {
		panic("Add: bad len")
	}
	result := slices.Clone(a)
	vek_And_AVX2(result, b)

	return result
}

func bs_or(a, b []bool) []bool {
	return v1_bs_or(a, b)
}

func v1_bs_or(a, b []bool) []bool {
	an := len(a)
	bn := len(b)
	if an != bn {
		panic("Or: bad len")
	}
	length := an
	result := make([]bool, length)
	n := b8x32_or(a, b, result)

	for i := length - n; i < length; i++ {
		result[i] = a[i] || b[i]
	}
	return result
}

func bs_xor(a, b []bool) []bool {
	return v1_bs_xor(a, b)
}

func v1_bs_xor(a, b []bool) []bool {
	an := len(a)
	bn := len(b)
	if an != bn {
		panic("Or: bad len")
	}
	length := an
	result := make([]bool, length)
	n := b8x32_xor(a, b, result)

	for i := length - n; i < length; i++ {
		result[i] = bool_xor(a[i], b[i])
	}
	return result
}
