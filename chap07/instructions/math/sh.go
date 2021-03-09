package math

import (
	"jvmgo/chap07/instructions/base"
	"jvmgo/chap07/rtda"
)

type ISHL struct { base.NoOperandsInstruction }
type ISHR struct { base.NoOperandsInstruction }
type IUSHR struct { base.NoOperandsInstruction }
type LSHL struct { base.NoOperandsInstruction }
type LSHR struct { base.NoOperandsInstruction }
type LUSHR struct { base.NoOperandsInstruction }

func (self *ISHL) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v2 := stack.PopInt()	//v2为需要左移的位数
	v1 := stack.PopInt()
	s := uint32(v2) & 0x1f	//只需要取v2的前5比特就足够表示位移位数了
	result := v1 << s
	stack.PushInt(result)
}

func (self *ISHR) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v2 := stack.PopInt()	//v2为需要左移的位数
	v1 := stack.PopInt()
	s := uint32(v2) & 0x1f
	result := v1 >> s
	stack.PushInt(result)
}

func (self *IUSHR) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v2 := stack.PopInt()	//v2为需要左移的位数
	v1 := stack.PopInt()
	s := uint32(v2) & 0x1f
	result := int32(uint32(v1) >> s)	//先转无符号，右移后转有符号
	stack.PushInt(result)
}

func (self *LSHL) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v2 := stack.PopInt()	//v2为需要左移的位数
	v1 := stack.PopLong()
	s := uint32(v2) & 0x3f	//只需要取v2的前5比特就足够表示位移位数了
	result := v1 << s
	stack.PushLong(result)
}

func (self *LSHR) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v2 := stack.PopInt()	//v2为需要左移的位数
	v1 := stack.PopLong()
	s := uint32(v2) & 0x3f
	result := v1 >> s
	stack.PushLong(result)
}

func (self *LUSHR) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v2 := stack.PopInt()	//v2为需要左移的位数
	v1 := stack.PopLong()
	s := uint32(v2) & 0x3f
	result := int64(uint64(v1) >> s)	//先转无符号，右移后转有符号
	stack.PushLong(result)
}