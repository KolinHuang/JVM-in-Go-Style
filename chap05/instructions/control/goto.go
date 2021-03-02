package control

import (
	"jvmgo/chap05/instructions/base"
	"jvmgo/chap05/rtda"
)

// Branch always
type GOTO struct{ base.BranchInstruction }

func (self *GOTO) Execute(frame *rtda.Frame) {
	base.Branch(frame, self.Offset)
}
