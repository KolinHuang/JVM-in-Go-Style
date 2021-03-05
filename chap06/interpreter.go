package main

import (
	"fmt"
	"jvmgo/chap06/classfile"
	"jvmgo/chap06/instructions"
	"jvmgo/chap06/instructions/base"
	"jvmgo/chap06/rtda"
)

func interpret(methodInfo *classfile.MemberInfo){
	//获得到方法的Code属性
	codeAttr := methodInfo.CodeAttribute()
	//获得最大栈深度
	maxLocals := codeAttr.MaxLocals()
	//获得最大局部变量表长度
	maxStack := codeAttr.MaxStack()
	//获得字节码
	bytecode := codeAttr.Code()

	thread := rtda.NewThread()
	frame := thread.NewFrame( maxLocals, maxStack)
	thread.PushFrame(frame)

	defer catchErr(frame)
	loop(thread, bytecode)
}

func catchErr(frame *rtda.Frame){
	if r := recover(); r != nil{
		fmt.Printf("LocalVars:%v\n", frame.LocalVars())
		fmt.Printf("OperandStack:%v\n", frame.OperandStack())
		panic(r)
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