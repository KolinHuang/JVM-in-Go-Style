package main

import (
	"fmt"
	"jvmgo/chap05/classfile"
	"jvmgo/chap05/classpath"
	"jvmgo/chap05/rtda"
	"strings"
)

func main() {
	cmd := parseCmd()
	cmd.Xjre = "/Library/Java/JavaVirtualMachines/jdk1.8.0_201.jdk/Contents/Home/jre"
	cmd.clspath = "/Users/huangyucai/Documents/code/git_depositorys/github_KolinHuang/JVM-in-Go-Style/javafiles/"
	cmd.class = "GaussTest"
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
	cp := classpath.Parse(cmd.Xjre, cmd.clspath)
	className := strings.Replace(cmd.class, ".", "/", -1)
	cf := loadClass(className, cp)
	mainMethod := getMainMethod(cf)
	if mainMethod == nil {
		fmt.Printf("Main method not found in class %s\n", cmd.class)
	}else{
		interpret(mainMethod)
	}

}

func getMainMethod(cf *classfile.ClassFile) *classfile.MemberInfo{
	for _, m := range cf.Methods() {
		if m.Name() == "main" && m.Descriptor() == "([Ljava/lang/String;)V"{
			return m
		}
	}
	return nil
}

func loadClass(className string, cp *classpath.Classpath) *classfile.ClassFile{
	classData, _, err := cp.ReadClass(className)
	if err != nil {
		panic(err)
	}
	cf, err := classfile.Parse(classData)
	if err != nil{
		panic(err)
	}
	return cf
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
