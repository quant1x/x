package mem

import "reflect"

func TypeSize[T any]() uintptr {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	return typ.Size()
}
