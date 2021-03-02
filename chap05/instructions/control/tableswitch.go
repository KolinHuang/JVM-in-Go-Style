package control

import (
	"jvmgo/chap05/instructions/base"
	"jvmgo/chap05/rtda"
)

type TABLE_SWITCH struct {
	defaultOffset int32	//默认情况下执行跳转所需的字节码偏移量
	low int32	//记录case的取值下限
	high int32	//记录case的取值上限
	jumpOffsets []int32	//一个索引表，存放high - low + 1个int值，对应各种case下，执行跳转所需的字节码偏移量
}

func (self *TABLE_SWITCH) FetchOperands( reader *base.BytecodeReader){
	reader.SkipPadding()
	self.defaultOffset = reader.ReadInt32()
	self.low = reader.ReadInt32()
	self.high = reader.ReadInt32()
	jumpOffsetsCount := self.high - self.low + 1
	self.jumpOffsets = reader.ReadInt32s(jumpOffsetsCount)
}

func (self *TABLE_SWITCH) Execute(frame *rtda.Frame){
	index := frame.OperandStack().PopInt()
	var offset int
	if index >= self.low && index <= self.high {
		offset = int(self.jumpOffsets[index - self.low])	//查表
	}else {
		offset = int(self.defaultOffset)	//默认
	}
	base.Branch(frame, offset)
}
