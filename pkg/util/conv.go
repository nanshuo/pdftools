package util

import (
	"unsafe"
	"reflect"
)

func B2s(b []byte)string{
	return *(*string)(unsafe.Pointer(&b))
}

func S2b(s string)[]byte{
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
