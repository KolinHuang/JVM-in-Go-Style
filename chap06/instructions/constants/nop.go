package constants

import (
	"jvmgo/chap06/instructions/base"
	"jvmgo/chap06/rtda"
)

type NOP struct {
	base.NoOperandsInstruction
}

func (self *NOP) Execute(frame *rtda.Frame){
	//do nothing
}
