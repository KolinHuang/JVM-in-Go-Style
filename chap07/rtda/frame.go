package rtda

import "jvmgo/chap07/rtda/heap"

type Frame struct {
	lower	*Frame	//指向此栈帧下方的第一个栈帧
	localVars	LocalVars
	operandStack	*OperandStack
	thread       *Thread
	method       *heap.Method
	nextPC       int // the next instruction after the call
}
//创建栈帧：执行方法的局部变量表大小和操作数栈深度是由编译器预先计算好的，存储在class文件method_info结构的Code属性中
//从method中读取栈深度和局部变量表长度
func newFrame(thread *Thread, method *heap.Method) *Frame {
	return &Frame{
		thread:thread,
		method:       method,
		localVars:    newLocalVars(method.MaxLocals()),
		operandStack: newOperandStack(method.MaxStack()),
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
func (self *Frame) Method() *heap.Method {
	return self.method
}
func (self *Frame) NextPC() int {
	return self.nextPC
}
func (self *Frame) SetNextPC(nextPC int) {
	self.nextPC = nextPC
}

func (self *Frame) RevertNextPC() {
	self.nextPC = self.thread.pc
}