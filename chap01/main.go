package main

import (
	"fmt"
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
	fmt.Println("JVM start running...")
	fmt.Printf("classpath: %s class: %s args: %v\n",
		cmd.clspath, cmd.class, cmd.args)
}