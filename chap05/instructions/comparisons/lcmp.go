package comparisons

import (
	"jvmgo/chap05/instructions/base"
	"jvmgo/chap05/rtda"
)

type LCMP struct { base.NoOperandsInstruction }

func (self *LCMP) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v1 := stack.PopLong()
	v2 := stack.PopLong()
	i := int32(0)
	if v1 > v2 {
		i = 1
	}else if v1 < v2{
		i = -1
	}
	stack.PushInt(i)
}