package rtda

import "jvmgo/chap07/rtda/heap"

type Slot struct {
	num int32
	ref *heap.Object
}