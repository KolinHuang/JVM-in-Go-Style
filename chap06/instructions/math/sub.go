package math

import (
	"jvmgo/chap06/instructions/base"
	"jvmgo/chap06/rtda"
)

type ISUB struct { base.NoOperandsInstruction }
type LSUB struct { base.NoOperandsInstruction }
type FSUB struct { base.NoOperandsInstruction }
type DSUB struct { base.NoOperandsInstruction }

func (self *ISUB) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v1 := stack.PopInt()
	v2 := stack.PopInt()
	result := v1 - v2
	stack.PushInt(result)
}

func (self *LSUB) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v1 := stack.PopLong()
	v2 := stack.PopLong()
	result := v1 - v2
	stack.PushLong(result)
}

func (self *FSUB) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v1 := stack.PopFloat()
	v2 := stack.PopFloat()
	result := v1 - v2
	stack.PushFloat(result)
}

func (self *DSUB) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v1 := stack.PopDouble()
	v2 := stack.PopDouble()
	result := v1 - v2
	stack.PushDouble(result)
}


