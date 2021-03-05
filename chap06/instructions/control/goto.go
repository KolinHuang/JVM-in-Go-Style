package control

import (
	"jvmgo/chap06/instructions/base"
	"jvmgo/chap06/rtda"
)

// Branch always
type GOTO struct{ base.BranchInstruction }

func (self *GOTO) Execute(frame *rtda.Frame) {
	base.Branch(frame, self.Offset)
}
