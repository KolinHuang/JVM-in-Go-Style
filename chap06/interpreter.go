package main

import (
	"fmt"
	"jvmgo/chap06/instructions"
	"jvmgo/chap06/instructions/base"
	"jvmgo/chap06/rtda"
	"jvmgo/chap06/rtda/heap"
)

func interpret(method *heap.Method){
	thread := rtda.NewThread()
	frame := thread.NewFrame(method)
	thread.PushFrame(frame)
	defer catchErr(frame)
	loop(thread, method.Code())
}

func catchErr(frame *rtda.Frame){
	if r := recover(); r != nil{
		//fmt.Printf("LocalVars:%v\n", frame.LocalVars())
		//fmt.Printf("OperandStack:%v\n", frame.OperandStack())
		//panic(r)
	}
}

func loop(thread *rtda.Thread, bytecode []byte){
	frame := thread.PopFrame()
	reader := &base.BytecodeReader{}

	for{
		pc := frame.NextPC()
		thread.SetPC(pc)

		//decode
		reader.Reset(bytecode, pc)
		opcode := reader.ReadUint8()//取操作码
		inst := instructions.NewInstruction(opcode)//创建指令
		inst.FetchOperands(reader)//取操作数
		frame.SetNextPC(reader.PC())//设置PC

		//execute
		fmt.Printf("pc:%2d inst:%T %v\n", pc, inst, inst)
		inst.Execute(frame)//执行
	}
}