# JVM-in-Go-Style



## 1 命令行工具

### 1.1 准备工作

**操作系统**：macOS 10.15.7 

**java version**： "1.8.0_201"

**go version**： go1.15.3 darwin/amd64

#### 1.1.1/2 安装JDK和Go

省略了安装的步骤。

go 命令行希望所有的Go源代码都被放在一个工作空间中。所谓工作空间，实际上就是一个目录结构，这个目录结构包含三个子目录。

* .src目录中是Go语言的源代码
* .pkg目录中是编译好的包对象文件
* .bin目录是链接好的可执行文件

只有src目录是我们需要写的，go会自动创建其余两个目录，工作空间可以位于任何地方，本实现（Mac环境下）的工作空间为:

```shell
/Users/huangyucai/go
```

如果要自定义工作空间，可使用以下命令：

```shell
#查看目前的环境
go env

vim ~/.bash_profile
#修改GPPATH项
GOPATH=/xxxxxx
#使其生效
source ~/.bash_profile
#查看效果
go env
```

#### 1.1.3 创建目录结构

Go语言以包为单位组织源代码，包可以嵌套，形成层次关系。本实现的所有源代码放在jvmgo包中，所以首先创建目录如下

![image-20210104193557762](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210104193557762.png)







### 1.2 Java命令

Java程序有一个主类用来启动Java应用，此类中包含一个main()方法。Java虚拟机规范没有明确规定JVM如何获取和启动主类，所以应由虚拟机实现自行决定启动方式。Oracle的Java虚拟机实现使用Java命令来启动，主类名在命令行参数中指定。Java命令有以下4种形式：

```shell
java [-options] class [args]
java [-options] -jar jarfile [args]
javaw [-options] class [args]
javaw [-options] -jar jarfile [args]
```

通常，第一个非选项参数给出主类的完全限定名，但是如果用户指定了-jar参数，那第一个非选项参数表示JAR文件名，Java命令必须从这个JAR文件中寻找主类。javaw命令和java命令几乎完全一样，唯一的差别在于，javaw命令不显示命令行窗口，因此特别适合用于启动GUI应用程序。

**常用选项**

![image-20210104194820326](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210104194820326.png)





### 1.3 编写命令行工具

在chap01目录下创建文件`cmd.go`，并编写结构体代码：

```go
type Cmd struct {
	helpFlag	bool//用于指定是否显示帮助信息
	versionFlag	bool//用于指定是否显示版本信息
	cpOption	string//用于指定类路径
	class		string//用于指定主类
	args		[]string//用于指定其他参数
}
```

再编写一个用于解析命令行参数的函数，在Go中可以直接处理os.Args变量来完成以上需求，但是比较麻烦。Go还内置了flag包可以帮助我们优雅地处理命令行选项。[API地址](https://golang.google.cn/pkg/flag/)

常用的几个方法：

```go
//声明一个flag: "-n"，存储在类型为*int的指针nFlag中
import "flag"
var nFlag = flag.Int("n", 1234, "help message for flag n")
//还可以用Var方法将一个flag绑定到一个变量
var flagvar int
func init() {
	flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
}

//定义完flag后，调用以下方法将命令行参数解析到上面定义好的flags中
flag.Parse()
```

本实现中，处理命令行的函数

```go
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
```

printUsage函数会在flag.Parse解析失败后被调用，显示命令的用法，如果解析成功，命令行参数将被解析到结构体对应的变量中。



### 1.4 测试命令行工具

在chap01目录下编写main.go文件：

```go
package main
import "fmt"

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
```

测试命令`chap01 -help`，输出如下所示：

![image-20210104204359043](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210104204359043.png)

测试命令`chap01 -cp java/classpath MyApp arg1 arg2 arg3`

![image-20210104204610033](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210104204610033.png)





## 2 搜索class文件

启动JVM虚拟机后，需要进行类加载操作，要进行类加载就得知道需要加载的类的路径。



### 2.1 类路径

类加载器的类别主要有三种：

* 启动类加载器（Bootstrap ClassLoader）：加载Java的核心库，即Java_HOME/jre/lib/rt.jar、resources.jar或sun.boot.class.path路径下的内容
* 扩展类加载器（Extension ClassLoader）：从java.ext.dirs系统属性所指定的目录中加载类库，或从JDK的安装目录的jre/lib/ext子目录下加载类库
* 系统类加载器（Application ClassLoader）：负责加载环境变量classpath或系统属性java.class.path指定路径下的类库

系统类路径的默认值是当前目录，即“.”。可以设置classpath环境变量来修改用户类路径，但是不灵活，更好的办法是通过java命令传递-classpath参数来指定系统类加载器的加载路径。

-classpath选项即可以**指定目录**，也可以**指定JAR文件或者ZIP文件**，还可以同时**指定多个目录和文件**，用分隔符分开即可。在本实现下（MacOS），分隔符采用的是冒号。从Java6 开始，还可以**使用通配符（*）**指定某个目录下的所有JAR文件。



### 2.2 准备工作

创建目录chap02/classpath。将chap01的内容拷入chap02中。

Java虚拟机需要使用JDK的启动类路径来寻找和加载Java核心类，因此需要用某种方式指定jre目录的位置。因此采用命令行选项的方式指定。修改Cmd结构体，新增一个非标准选项-Xjre

```go
	Xjre		string//用于指定jre路径
```



### 2.3 实现类路径

如果我们把类路径想象成一个大的整体，那么它由启动类路径、扩展类路径和用户（系统）类路径三个小路径构成，三个小路径又分别由更小的路径构成。我们就可以套用组合模式（composite pattern）来设计和实现类路径。

> 插播组合模式
>
> #### 角色
>
> **Component（抽象构件）**：它可以是接口或抽象类，为叶子构件和容器构件对象声明接口，在该角色中可以包含所有子类共有行为的声明和实现。在抽象构件中定义了访问及管理它的子构件的方法，如增加子构件、删除子构件、获取子构件等。
>
> **Leaf（叶子构件）**：它在组合结构中表示叶子节点对象，叶子节点没有子节点，它实现了在抽象构件中定义的行为。对于那些访问及管理子构件的方法，可以通过异常等方式进行处理。
>
> **Composite（容器构件）**：它在组合结构中表示容器节点对象，容器节点包含子节点，其子节点可以是叶子节点，也可以是容器节点，它提供一个集合用于存储子节点，实现了在抽象构件中定义的行为，包括那些访问及管理子构件的方法，在其业务方法中可以递归调用其子节点的业务方法。
>
> 组合模式的**关键是定义了一个抽象构件类，它既可以代表叶子，又可以代表容器**，而客户端针对该抽象构件类进行编程，无须知道它到底表示的是叶子还是容器，可以对其进行统一处理。**同时容器对象与抽象构件类之间还建立一个聚合关联关系**，在容器对象中既可以包含叶子，也可以包含容器，以此实现递归组合，形成一个树形结构。
>
>
> 作者：小旋锋
> 链接：https://juejin.cn/post/6844903687228407821
> 来源：掘金
> 著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。





#### 2.3.1 Entry接口

定义entry接口来表示类路径项。

```go
package classpath

import(
"os"
"strings"
)

//定义路径分隔符常量，":"
const pathListSepatator = string(os.PathListSeparator)
type Entry interface {
	//负责寻找和加载class文件
	readClass(className string) ([]byte, Entry, error)
	//负责返回变量的字符串表示，相当于Java的toString()方法
	String() string
}

func newEntry(path string) Entry {
	if strings.Contains(path, pathListSepatator) {
		return newCompositeEntry(path)
	}

	if strings.HasSuffix(path, "*") {
		return newWildcardEntry(path)
	}

	if strings.HasSuffix(path, ".jar") ||
		strings.HasSuffix(path, ".JAR") ||
		strings.HasSuffix(path, ".zip") ||
		strings.HasSuffix(path, ".ZIP") {
		return newZipEntry(path)
	}

	return newDirEntry(path)
}
```

readClass()函数的参数是class文件的相对路径，路径之间用“/”分隔，文件名有.class后缀。返回值是读到的字节数据、最终定位到class文件的Entry以及错误信息。

newEntry()函数根据参数创建不同类型的Entry实例。正如2.1结尾说的，类路径可以有多种形式：

1. JAR、ZIP文件——>关键词：.JAR、.ZIP后缀
2. 含分隔符的多个文件和目录——>关键词：“:”分隔符
3. 通配符表示的多个文件——>关键词："*"通配符
4. 目录——>剩余的就是目录

根据以上分析，Entry接口可以有4个实现，分别是DirEntry（目录）, ZipEntry（文件）, CompositeEntry（多个文件和目录）和WildcardEntry（含通配符）。





#### 2.3.2 Entry接口的实现

> 插播：在Golang中，方法和函数的区别
>
> 在Go语言中，方法和函数只差了一个地方，方法在func和标识符之间多了一个参数，这个参数就称作接收者（receiver）。接收者有两种：一种是值接收者，一种是指针接收者。

和Java语言不同，Go结构体不需要显式地实现接口，只要方法匹配即可。Go没有专门的构造函数。本实现统一使用new开头的函数来创建结构体实例。

1. DirEntry实现：

   ```go
   package classpath
   
   import(
   	"io/ioutil"
   	"path/filepath"
   )
   
   type DirEntry struct {
   	absDir	string
   }
   
   //返回目录形式类路径项实例
   func newDirEntry(path string) *DirEntry {
   	//Abs函数返回路径的绝对路径表示
   	absDir, err := filepath.Abs(path)
   	if err != nil {
   		//当调用了panic函数，此函数的执行将会停止
   		panic(err)
   	}
   	//返回DirEntry实例
   	return &DirEntry{absDir}
   }
   
   func (self *DirEntry) readClass(className string) ([]byte, Entry, error){
   	//把目录和class文件名拼成一个完整的路径
   	fileName := filepath.Join(self.absDir, className)
   	//读取class文件的内容
   	data, err := ioutil.ReadFile(fileName)
   	return data, self, err
   }
   
   func (self *DirEntry) String() string {
   	return self.absDir
   }
   
   ```

2. ZipEntry实现：

   ```go
   package classpath
   
   import(
   	"archive/zip"	//提供读取和写入ZIP压缩包的操作
   	"errors"
   	"io/ioutil"	//IO工具包，提供一些IO操作
   	"path/filepath"
   )
   
   type ZipEntry struct {
   	absPath string	//存放JAR/ZIP文件的绝对路径
   }
   
   func newZipEntry(path string) *ZipEntry {
   	absPath, err := filepath.Abs(path)
   	if err != nil {
   		panic(err)
   	}
   	return &ZipEntry{absPath}
   }
   
   //从ZIP文件中提取class文件
   func (self *ZipEntry) readClass(className string) ([]byte, Entry, error) {
   	//OpenReader will open the Zip file specified by name and return a ReadCloser.
   	r, err := zip.OpenReader(self.absPath)
   	if err != nil {
   		return nil, nil, err
   	}
   	//defer类似于Java中的finally语句块，在函数返回之前，或者说在return语句后执行。
   	defer r.Close()
   	//遍历压缩包内的文件
   	for _, f := range r.File {
   		//如果找到了className对应的class文件，就打开并读取内容，然后返回
   		if f.Name == className {
   			//打开文件
   			rc, err := f.Open()
   			if err != nil {//打开失败，返回错误
   				return nil, nil, err
   			}
   			defer rc.Close()
   			//读取文件所有数据，返回字节数组
   			data, err := ioutil.ReadAll(rc)
   			if err != nil {//读取失败
   				return nil, nil, err
   			}
   			return data, self, nil
   		}
   	}
   	return nil, nil, errors.New("class not found: " + className)
   
   }
   
   func (self *ZipEntry) String() string {
   	return self.absPath
   }
   
   ```

3. CompositeEntry的实现：

   ```go
   package classpath
   
   import(
   	"errors"
   	"strings"
   )
   /**
   多个Entry构成的类路径
    */
   type CompositeEntry []Entry
   //将路径列表参数按分隔符分成小路径，然后将每个小路径转化为具体的Entry实例
   func newCompositeEntry(pathList string) CompositeEntry {
   	compositeEntry := []Entry{}
   	for _, path := range strings.Split(pathList, pathListSepatator) {
   		//调用Entry接口中的newEntry函数
   		entry := newEntry(path)
   		compositeEntry = append(compositeEntry, entry)
   	}
   	return compositeEntry
   }
   
   func (self CompositeEntry) readClass(className string) ([]byte, Entry, error) {
   	for _, entry := range self {
   		data, from, err := entry.readClass(className)
   		if err == nil {
   			return data, from, nil
   		}
   	}
   	return nil, nil, errors.New("class not found: " + className)
   }
   //调用每一个子路径的String方法，然后把得到的字符串用路径分隔符拼接起来即可
   func (self CompositeEntry) String() string {
   	strs := make([]string, len(self))
   
   	for i, entry := range self {
   		strs[i] = entry.String()
   	}
   	return strings.Join(strs, pathListSepatator)
   }
   ```

4. WildcardEntry实现：

   ```go
   package classpath
   
   import(
   	"os"
   	"path/filepath"
   	"strings"
   )
   //WildcardEntry实际上也是CompositeEntry类型的
   func newWildcardEntry(path string) CompositeEntry {
   	baseDir := path[: len(path) - 1]	//删除"*"号
   	compositeEntry := []Entry{}
   	//匿名函数
   	walkFn := func(path string, info os.FileInfo, err error) error {
   		if err != nil {
   			return err
   		}
   		//通配符类路径不能递归匹配子目录下的JAR文件
   		if info.IsDir() && path != baseDir {
   			return filepath.SkipDir
   		}
   		//找出JAR文件
   		if strings.HasSuffix(path, ".jar") ||
   			strings.HasSuffix(path,".JAR") {
   			jarEntry := newZipEntry(path)
   			compositeEntry = append(compositeEntry, jarEntry)
   		}
   		return nil
   	}
   	//遍历baseDir，将所有JAR文件创建为ZipEntry，放入数组compositeEntry中
   	//Walk walks the file tree rooted at root, calling walkFn for each file
   	//or directory in the tree, including root.
   	filepath.Walk(baseDir, walkFn)
   	return compositeEntry
   }
   ```

   

