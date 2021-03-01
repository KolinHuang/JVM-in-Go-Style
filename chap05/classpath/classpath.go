package classpath

import(
	"os"
	"path/filepath"
)

type Classpath struct {
	bootClasspath	Entry	//启动类路径
	extClasspath	Entry	//扩展类路径
	userClasspath	Entry	//用户类路径
}
//使用-Xjre选项指定的类路径来解析启动类路径和扩展类路径
//使用-clspath选项指定的类路径来解析用户类路径
func Parse(Xjre, clspath string) *Classpath {
	clsp := &Classpath{}
	clsp.parseBootAndExtClasspath(Xjre)
	clsp.parseUserClasspath(clspath)
	return clsp
}
//双亲委派模型，优先从启动类路径查找类文件，然后再到扩展类路径下查找，最后在用户类路径下查找
func (self *Classpath) ReadClass(className string) ([]byte, Entry, error) {
	className = className + ".class"
	if data, entry, err := self.bootClasspath.readClass(className); err == nil {
		return data, entry, err
	}
	if data, entry, err := self.extClasspath.readClass(className); err == nil {
		return data, entry, err
	}
	return self.userClasspath.readClass(className)
}

func (self *Classpath) String() string {
	return self.userClasspath.String()
}

func (self *Classpath) parseBootAndExtClasspath(Xjre string) {
	jreDir := getJreDir(Xjre)
	//jre/lib/*
	jreLibPath := filepath.Join(jreDir, "lib", "*")
	//为启动类路径创建一个通配符形式的Entry
	self.bootClasspath = newWildcardEntry(jreLibPath)
	//jre/lib/ext/*
	jreExtPath := filepath.Join(jreDir, "lib", "ext", "*")
	self.extClasspath = newWildcardEntry(jreExtPath)
}

func (self *Classpath) parseUserClasspath(clspath string){
	if clspath == "" {
		clspath = "."
	}
	self.userClasspath = newEntry(clspath)
}

func getJreDir(Xjre string) string {
	//如果指定的jre路径存在，就返回
	if Xjre != "" && exsits(Xjre){
		return Xjre
	}
	//如果指定的jre路径不存在，就到当前目录下寻找jre目录
	if exsits("./jre") {
		return "./jre"
	}
	//如果当前目录下不存在jre目录，就到系统环境变量JAVA_HOME指示的路径下找jre目录
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		return filepath.Join(jh, "jre")
	}
	panic("Can not find jre folder! ")
}
//用于判断目录是否存在
func exsits(path string) bool {
	//os.Stat用来返回path代表的路径是否有问题，如果有问题就会返回一个错误
	//这个错误将作为os.IsNotExist的输入，用来指示是否是因为这个path不存在而引发的错误
	//这么绕，也太晦涩了，无非就是想表达这个路径存不存在
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err){
			return false
		}
		//if os.IsExist(err) {
		//	return true
		//}else{
		//	return false
		//}
	}
	return true
}

