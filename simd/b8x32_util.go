package simd

// bool转int8
func boolToInt8(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

// int8ToBool 将非零的int值转换为true，零值转换为false
func int8ToBool(i int8) bool {
	return i != 0
}

func bool_xor(x, y bool) bool {
	a := boolToInt8(x)
	b := boolToInt8(y)
	c := a ^ b
	return int8ToBool(c)
}
