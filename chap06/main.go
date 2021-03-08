package main

import (
	"fmt"
	"jvmgo/chap06/classpath"
	"jvmgo/chap06/rtda/heap"
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
	classLoader := heap.NewClassLoader(cp)
	className := strings.Replace(cmd.class, ".", "/", -1)
	mainClass := classLoader.LoadClass(className)
	mainMethod := mainClass.GetMainMethod()

	if mainMethod == nil {
		fmt.Printf("Main method not found in class %s\n", cmd.class)
	}else{
		interpret(mainMethod)
	}

}

