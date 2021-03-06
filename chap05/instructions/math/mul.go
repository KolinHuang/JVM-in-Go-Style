package math

import (
"jvmgo/chap05/instructions/base"
"jvmgo/chap05/rtda"
)

type IMUL struct { base.NoOperandsInstruction }
type LMUL struct { base.NoOperandsInstruction }
type FMUL struct { base.NoOperandsInstruction }
type DMUL struct { base.NoOperandsInstruction }

func (self *IMUL) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v1 := stack.PopInt()
	v2 := stack.PopInt()
	result := v1 * v2
	stack.PushInt(result)
}

func (self *LMUL) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v1 := stack.PopLong()
	v2 := stack.PopLong()
	result := v1 * v2
	stack.PushLong(result)
}

func (self *FMUL) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v1 := stack.PopFloat()
	v2 := stack.PopFloat()
	result := v1 * v2
	stack.PushFloat(result)
}

func (self *DMUL) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v1 := stack.PopDouble()
	v2 := stack.PopDouble()
	result := v1 * v2
	stack.PushDouble(result)
}


