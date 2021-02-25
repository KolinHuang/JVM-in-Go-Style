package main

import (
	"fmt"
	"jvmgo/chap04/rtda"
)

func main() {
	cmd := parseCmd()
	if cmd.versionFlag {
		fmt.Println("version: 0.0.1")
	} else if cmd.helpFlag || cmd.class == ""{
		//用户指定了helpFlag参数或者未指定主类，就打印命令用法
		printUsage()
	} else {
		//一切正常就启动Java虚拟机
		startJVM(cmd)
	}
}

func startJVM(cmd *Cmd){
	frame := rtda.NewFrame(100, cmd.maxStackSize)//创建栈帧
	testLocalVars(frame.LocalVars())//操作局部变量表
	testOperandStack(frame.OperandStack())//操作操作数栈
}

func testLocalVars(vars rtda.LocalVars) {
	vars.SetInt(0, 100)
	vars.SetInt(1, -100)
	vars.SetLong(2, 2997924580)
	vars.SetLong(4, -2997924580)
	vars.SetFloat(6, 3.1415926)
	vars.SetDouble(7, 2.71828182845)
	vars.SetRef(9, nil)
	println(vars.GetInt(0))
	println(vars.GetInt(1))
	println(vars.GetLong(2))
	println(vars.GetLong(4))
	println(vars.GetFloat(6))
	println(vars.GetDouble(7))
	println(vars.GetRef(9))
}

func testOperandStack(ops *rtda.OperandStack) {
	ops.PushInt(100)
	ops.PushInt(-100)
	ops.PushLong(2997924580)
	ops.PushLong(-2997924580)
	ops.PushFloat(3.1415926)
	ops.PushDouble(2.71828182845)
	ops.PushRef(nil)
	println(ops.PopRef())
	println(ops.PopDouble())
	println(ops.PopFloat())
	println(ops.PopLong())
	println(ops.PopLong())
	println(ops.PopInt())
	println(ops.PopInt())
}
