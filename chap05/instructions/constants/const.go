package constants

import (
	"jvmgo/chap05/instructions/base"
	"jvmgo/chap05/rtda"
)

type ACONST_NULL struct { base.NoOperandsInstruction }	//将nil入栈
type DCONST_0 struct { base.NoOperandsInstruction }	//将double 0入栈
type DCONST_1 struct { base.NoOperandsInstruction }
type FCONST_0 struct { base.NoOperandsInstruction }	//将float 0入栈
type FCONST_1 struct { base.NoOperandsInstruction }
type FCONST_2 struct { base.NoOperandsInstruction }
type ICONST_M1 struct { base.NoOperandsInstruction }	//将int型-1入栈
type ICONST_0 struct { base.NoOperandsInstruction }	//将int 0入栈
type ICONST_1 struct { base.NoOperandsInstruction }
type ICONST_2 struct { base.NoOperandsInstruction }
type ICONST_3 struct { base.NoOperandsInstruction }
type ICONST_4 struct { base.NoOperandsInstruction }
type ICONST_5 struct { base.NoOperandsInstruction }
type LCONST_0 struct { base.NoOperandsInstruction }	//将long 0入栈
type LCONST_1 struct { base.NoOperandsInstruction }

func (self *ACONST_NULL) Execute(frame *rtda.Frame){
	frame.OperandStack().PushRef(nil)
}
func (self *DCONST_0) Execute(frame *rtda.Frame){
	frame.OperandStack().PushDouble(float64(0))
}
func (self *DCONST_1) Execute(frame *rtda.Frame){
	frame.OperandStack().PushDouble(float64(1))
}
func (self *FCONST_0) Execute(frame *rtda.Frame){
	frame.OperandStack().PushFloat(float32(0))
}
func (self *FCONST_1) Execute(frame *rtda.Frame){
	frame.OperandStack().PushFloat(float32(1))
}
func (self *FCONST_2) Execute(frame *rtda.Frame) {
	frame.OperandStack().PushFloat(2.0)
}
//将int型 -1推入操作数栈顶
func(self *ICONST_M1) Execute(frame *rtda.Frame){
	frame.OperandStack().PushInt(-1)
}
func(self *ICONST_0) Execute(frame *rtda.Frame){
	frame.OperandStack().PushInt(0)
}
func(self *ICONST_1) Execute(frame *rtda.Frame){
	frame.OperandStack().PushInt(1)
}
func(self *ICONST_2) Execute(frame *rtda.Frame){
	frame.OperandStack().PushInt(2)
}
func(self *ICONST_3) Execute(frame *rtda.Frame){
	frame.OperandStack().PushInt(3)
}
func(self *ICONST_4) Execute(frame *rtda.Frame){
	frame.OperandStack().PushInt(4)
}
func(self *ICONST_5) Execute(frame *rtda.Frame){
	frame.OperandStack().PushInt(5)
}

func(self *LCONST_0) Execute(frame *rtda.Frame){
	frame.OperandStack().PushLong(0)
}
func(self *LCONST_1) Execute(frame *rtda.Frame){
	frame.OperandStack().PushLong(1)
}

