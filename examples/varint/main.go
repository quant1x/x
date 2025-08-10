package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	// 编码
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutVarint(buf[:], -123456) // 编码有符号整数
	fmt.Printf("Encoded bytes: %v\n", buf[:n])

	// 解码
	value, _ := binary.Varint(buf[:n]) // 解码有符号整数
	fmt.Printf("Decoded value: %d\n", value)
}
