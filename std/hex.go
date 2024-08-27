package std

import "encoding/hex"

func ToHexString(data []byte) string {
	return hex.EncodeToString(data)
}

func FromHexString(s string) ([]byte, error) {
	return hex.DecodeString(s)
}
