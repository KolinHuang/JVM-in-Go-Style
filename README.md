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

> 补充：由于在搜索class文件步骤中出现了错误，于是我把代码都放在了目录/Users/huangyucai/golang下，然后将环境变量配置为
>
> `export GOROOT=/Users/huangyucai/go
> export GOBIN=$GOROOT/bin
> export GOPATH=/User/huangyucai/golang
> export PATH=$PATH:$GOBIN`
>
> 其中GOROOT代表go的根路径
>
> GOPATH代表工作空间
>
> 具体为什么，我还不清楚，可能与go的包管理机制有关，后面有时间再研究吧。

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

本实现中，处理命令行的函数为printUsage

printUsage函数会在flag.Parse解析失败后被调用，显示命令的用法，如果解析成功，命令行参数将被解析到结构体对应的变量中。



### 1.4 测试

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

newEntry()函数根据-clspath参数创建不同类型的Entry实例。正如2.1结尾说的，类路径可以有多种形式：

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

1. **目录形式的类路径**，只需要有一个绝对路径能够表示其位置即可，因此DirEntry的结构体只需要一个属性来存放绝对路径。在读取类文件时，只需要将指定的类和传入的类路径参数拼接，就能得到类文件的完整路径。
2. 同理，**压缩文件形式的类路径**，也只需要一个绝对路径来表示其位置。只不过在读取类文件的时候，需要用到`"archive/zip"`工具包来遍历Zip文件中的class文件，找到符合要求的class文件。
3. **多个文件和目录形式的类路径**，由于是由多个类路径组成，因此需要一个数组来存放这些类路径，数组元素的类型设置为[]Entry，可以聚合4种形式的类路径。在读取类路径时，按照entry的类型调用各自的方法，是目录形式的就按1所说的读取类文件，是压缩文件形式的就按2所说的读取类文件。
4. **通配符形式的类路径**和第3个无本质区别，但是有一点，通配符形式的类路径无法读取子目录下的类文件，因此在遍历过程中遇到目录之后需要跳过。



#### 2.3.3 Classpath

创建classpath结构体，存放三种类路径。

```go
type Classpath struct {
	bootClasspath	Entry	//启动类路径
	extClasspath	Entry	//扩展类路径
	userClasspath	Entry	//用户类路径
}
```

前面提过，我们利用命令行方式来指定以上三种类路径的加载路径，那就需要将命令行参数`-Xjre`和`-clspath`对应的值解析成类路径Entry的形式。其中，启动类路径是`xx/../jre/lib/*`，扩展类路径是`xx/../jre/lib/ext/*`

举个例子：

在执行java命令时传递了参数`-Xjre /User/hyc/JAVA/jre`，那么我们就需要根据路径`/User/hyc/JAVA/jre`解析出启动类路径和扩展类路径。

1. 首先要判断此路径是否存在，用到了os包下的`Stat()`函数来查看目录状态和`isExist()`函数来判断目录是否存在。
2. 如果路径存在，那我们就可以直接返回此路径作为jre路径进行进一步解析；如果路径不存在，我们就要使用默认的jre路径了，本实现给了2个途径获取，一个是当前目录下查找jre，一个是到系统环境变量`JAVA_HOME`表示的目录下查找jre。如果都没找到，就中断方法，抛出error。
3. 假设我们找到了jre路径，就将此路径与`/lib/*`拼接，创建启动类路径Entry；再将此路径与`/lib/ext/*`拼接，创建扩展类路径Entry。
4. 用户类路径Entry的创建过程比较简单，如果用户没有指定`-clspath`参数，我们就用当前目录作为用户类路径。



**在读取类文件时，按照双亲委派模型，我们优先从启动类路径加载，然后再从扩展类路径加载，最后在用户类路径加载类文件。**



至此，整个类路径的查找和解析过程已经实现了，接下来测试一下！



### 2.4 测试

修改main.go如下：

```go
...
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
```

运行命令`go install ./jvmgo/chap02`

![image-20210107202738898](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210107202738898.png)

报错：

```shell
jvmgo/chap02/main.go:5:2: cannot find package "." in:
	/Users/huangyucai/go/src/vendor/jvmgo/chap02/classpath
```

经过排查，发现是编译器找不到`/jvm/chap02/classpath`包。

翻看博客，猜测可能是环境变量配错了，将环境变量修改为

```shell
export GOROOT=/Users/huangyucai/go
export GOBIN=$GOROOT/bin
export GOPATH=/User/huangyucai/golang
export PATH=$PATH:$GOBIN
```

然后将jvmgo文件夹移动至在`/User/huangyucai/golang`目录下，再次运行命令`go install ./jvmgo/chap02`，运行成功。

接下来测试一下能否查找到本机jre目录下的java.lang.Object类文件。本机的jre路径是：

```shell
/Library/Java/JavaVirtualMachines/jdk1.8.0_201.jdk/Contents/Home/jre
```

执行命令：`chap02 -Xjre /Library/Java/JavaVirtualMachines/jdk1.8.0_201.jdk/Contents/Home/jre java.lang.Object`

结果如下：

![image-20210107205154169](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210107205154169.png)

成功读取到了java.lang.Object类文件的字节码。



### 小结

首先本章套用了组合模式来设计统一的类路径表示，定制了接口Entry以及其4种实现。

然后利用`-Xjre`参数来指定启动类加载路径和扩展类加载路径，利用`-clspath`参数来指定用户类加载路径。根据三种类加载路径的要求对命令行参数传递进来的路径字符串进行解析，转换成绝对路径。

最后通过命令行传递的主类名按双亲委派模型查找类文件，并读取类文件。



## 3 解析class文件

