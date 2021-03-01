package constants

import (
	"jvmgo/chap05/instructions/base"
	"jvmgo/chap05/rtda"
)

//push byte
type BIPUSH struct {
	val int8
}
//push short
type SIPUSH struct {
	val int16
}

func(self *BIPUSH) FetchOperands(reader *base.BytecodeReader) {
	self.val = reader.ReadInt8()
}

func(self *BIPUSH) Execute(frame rtda.Frame){
	frame.OperandStack().PushInt(int32(self.val))
}

func(self *SIPUSH) FetchOperands(reader *base.BytecodeReader) {
	self.val = reader.ReadInt16()
}

func(self *SIPUSH) Execute(frame rtda.Frame){
	frame.OperandStack().PushInt(int32(self.val))
}