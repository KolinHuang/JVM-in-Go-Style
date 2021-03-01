package constants

import (
	"jvmgo/chap05/instructions/base"
	"jvmgo/chap05/rtda"
)

type NOP struct {
	base.NoOperandsInstruction
}

func (self *NOP) Excute(frame *rtda.Frame){
	//do nothing
}
