package array

import (
	"unsafe"
)

func PushFront(array []interface{}, elem interface{}) []interface{} {
	res := make([]interface{}, len(array)+1, unsafe.Sizeof(array[0]))
	res[0] = elem
	for i, elem := range array {
		res[i+1] = elem
	}
	return res
}
