package main

import (
	"fmt"
	"jvmgo/chap02/classpath"
	"strings"
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
	clsp := classpath.Parse(cmd.Xjre, cmd.clspath)
	fmt.Printf("classpath: %v class: %v args: %v\n", clsp, cmd.class, cmd.args)
	className := strings.Replace(cmd.class,".","/",-1)
	classData, _, err := clsp.ReadClass(className)
	if err != nil {
		fmt.Printf("Could not find or load main class %s\n", cmd.class)
		return
	}
	fmt.Printf("class data: %v\n", classData)
}