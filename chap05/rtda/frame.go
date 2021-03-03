package rtda

type Frame struct {
	lower	*Frame	//指向此栈帧下方的第一个栈帧
	localVars	LocalVars
	operandStack	*OperandStack
	thread       *Thread
	nextPC       int // the next instruction after the call
}
//创建栈帧：执行方法的局部变量表大小和操作数栈深度是由编译器预先计算好的，存储在class文件method_info结构的Code属性中
func newFrame(thread *Thread, maxLocals, maxStack uint16) *Frame {
	return &Frame{
		thread:thread,
		localVars:    newLocalVars(maxLocals),
		operandStack: newOperandStack(maxStack),
	}
}

func (self *Frame) LocalVars() LocalVars {
	return self.localVars
}
func (self *Frame) OperandStack() *OperandStack {
	return self.operandStack
}
func (self *Frame) Thread() *Thread {
	return self.thread
}
func (self *Frame) NextPC() int {
	return self.nextPC
}
func (self *Frame) SetNextPC(nextPC int) {
	self.nextPC = nextPC
}