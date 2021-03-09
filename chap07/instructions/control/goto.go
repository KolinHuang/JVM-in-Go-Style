package control

import (
	"jvmgo/chap07/instructions/base"
	"jvmgo/chap07/rtda"
)

// Branch always
type GOTO struct{ base.BranchInstruction }

func (self *GOTO) Execute(frame *rtda.Frame) {
	base.Branch(frame, self.Offset)
}
