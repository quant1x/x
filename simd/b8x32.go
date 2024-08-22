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
	result := make([]bool, length, length)
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
	result := make([]bool, length, length)
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
