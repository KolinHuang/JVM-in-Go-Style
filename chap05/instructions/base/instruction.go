package base

import (
	"jvmgo/chap05/rtda"
)

type Instruction interface {
	FetchOperands(reader *BytecodeReader)//从字节码中提取操作数
	Execute(frame *rtda.Frame)//执行指令逻辑
}
//零操作数的指令
type NoOperandsInstruction struct {}

func (self *NoOperandsInstruction) FetchOperands(reader *BytecodeReader) {
	//nothing to do
}

//跳转指令
type BranchInstruction struct {
	Offset int//跳转偏移量
}
func (self *BranchInstruction) FetchOperands(reader *BytecodeReader) {
	//nothing to do
	self.Offset = int(reader.ReadInt16())
}

//存储和加载类指令需要根据索引存取局部变量表，索引由单字节操作数给出，所以把这类指令抽象成Index8Instruction
type Index8Instruction struct {
	Index uint//局部变量表的索引
}
func (self *Index8Instruction) FetchOperands(reader *BytecodeReader){
	self.Index = uint(reader.ReadUint8())
}
//有一些指令需要访问常量池，常量池索引由两字节操作数给出，所以把这类指令抽象成Index16Instruction
type Index16Instruction struct {
	Index uint//局部变量表的索引
}
func (self *Index16Instruction) FetchOperands(reader *BytecodeReader){
	self.Index = uint(reader.ReadUint16())
}
