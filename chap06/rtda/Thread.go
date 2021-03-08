package rtda

import "jvmgo/chap06/rtda/heap"

type Thread struct {
	pc	int
	stack *Stack
}
//创建线程实例

func NewThread() *Thread{
	return NewThread0(1024)
}

func NewThread0(maxStackSize uint) *Thread{
	return &Thread{
		stack: newStack(maxStackSize),
	}
}

func(self *Thread) PC() int {
	return self.pc
}

func(self *Thread) SetPC(pc int) {
	self.pc = pc
}

func(self *Thread) PushFrame(frame *Frame){
	self.stack.push(frame)
}

func(self *Thread) PopFrame() *Frame{
	return self.stack.pop()
}

func(self *Thread) CurrentFrame() *Frame{
	return self.stack.top()
}

func (self *Thread) NewFrame(method *heap.Method) *Frame{
	return newFrame(self, method)
}