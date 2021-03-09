package constants

import (
	"jvmgo/chap07/instructions/base"
	"jvmgo/chap07/rtda"
)

type NOP struct {
	base.NoOperandsInstruction
}

func (self *NOP) Execute(frame *rtda.Frame){
	//do nothing
}
