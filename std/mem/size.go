package mem

import (
	"reflect"
	"unsafe"
)

func v1TypeSize[T any]() uintptr {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	return typ.Size()
}

func TypeSize[T any]() uintptr {
	var e T
	return unsafe.Sizeof(e)
}
