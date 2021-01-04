package main

import(
	"flag"
	"fmt"
	"os"
)

type Cmd struct {
	helpFlag	bool//用于指定是否显示帮助信息
	versionFlag	bool//用于指定是否显示版本信息
	clspath	string//用于指定类路径
	class		string//用于指定主类
	args		[]string//用于指定其他参数
}

func parseCmd() *Cmd {
	//创建一个Cmd结构体对象
	cmd := &Cmd{}
	//如果Parse函数解析失败，就会调用printUsage函数把命令的用法打印到控制台
	flag.Usage = printUsage
	//将help和？这两个参数绑定到变量&cmd.helpFlag，初值为false，信息为print help message
	//这样如果在命令行中出现help和?参数，说明需要打印帮助信息
	flag.BoolVar(&cmd.helpFlag, "help", false, "print help message-打印帮助信息")
	flag.BoolVar(&cmd.helpFlag, "? ", false, "print help message-打印帮助信息")
	flag.BoolVar(&cmd.versionFlag, "version", false, "print version and exit-打印命令版本")
	flag.StringVar(&cmd.clspath, "classpath", "", "classpath-指定类文件路径")
	flag.StringVar(&cmd.clspath, "cp", "", "classpath-指定类文件路径")
	//当所有flag定义完毕，调用此方法来解析命令行参数到flags中
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		cmd.class = args[0]
		cmd.args = args[1: ]
	}
	return cmd
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Printf("%s [-options] class [args...]\n", os.Args[0])
	fmt.Println("Options:")
	flag.PrintDefaults()
}