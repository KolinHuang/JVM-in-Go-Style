package loads

import (
	"jvmgo/chap07/instructions/base"
	"jvmgo/chap07/rtda"
)

type ALOAD struct { base.Index8Instruction }
type ALOAD_0 struct { base.NoOperandsInstruction }
type ALOAD_1 struct { base.NoOperandsInstruction }
type ALOAD_2 struct { base.NoOperandsInstruction }
type ALOAD_3 struct { base.NoOperandsInstruction }


func _aload(frame *rtda.Frame, index uint){
	val := frame.LocalVars().GetRef(index)
	frame.OperandStack().PushRef(val)
}

func (self *ALOAD) Execute(frame *rtda.Frame){
	_dload(frame, self.Index)
}

func (self *ALOAD_0) Execute(frame *rtda.Frame){
	_dload(frame, 0)
}
func (self *ALOAD_1) Execute(frame *rtda.Frame){
	_dload(frame, 1)
}
func (self *ALOAD_2) Execute(frame *rtda.Frame){
	_dload(frame, 2)
}
func (self *ALOAD_3) Execute(frame *rtda.Frame){
	_dload(frame, 3)
}