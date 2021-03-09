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

> 2020-1-7-21:12补充：由于在搜索class文件步骤中出现了错误，于是我把代码都放在了目录/Users/huangyucai/golang下，然后将环境变量配置为
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

经过排查，发现是编译器找不到`jvmgo/chap02/classpath`包。

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

### 3.1 class文件结构



构成class文件的基本数据单位是字节，可以把整个class文件当成一个字节流来处理。Java虚拟机规范定义了u1, u2, u4, u8四种数据类型来表示1、2、4和8字节无符号数，连续多个字节构成的数据按大端格式存储（高位在前，低位在后）。以上四种数据类型分别对应Go语言的uint8、uint16、uint32和uint64类型。

Class文件中还用了一种称为“表”的数据类型来表示数据。表是由多个无符号数或者其他表作为数据项构成的复合数据类型，所有表的命名都习惯性地以“_info”结尾。

无论是无符号数还是表，当需要描述同一类型但数量不定的多个数据时，经常会使用一个**前置的容量计数器**加若干**连续数据项**的形式描述。

Class的结构不像XML等描述语言，由于它没有任何分隔符号，所以所有数据项，无论是顺序还是数量，甚至于数据存储的字节序（Byte Ordering，Class 文件中字节序为Big-Endian）这样的细节，都是被严格限定的**，哪个字节代表什么含义，长度是多少， 先后顺序如何，全部都不允许改变。**

Class文件总体的结构如下表所示：

|      类型      |                名称                |          数量           |
| :------------: | :--------------------------------: | :---------------------: |
|       u4       |             magic/魔数             |            1            |
|       u2       |      minor_version/最小版本号      |            1            |
|       u2       |      major_version/主要版本号      |            1            |
|       u2       | constant_pool_count/常量池容量计数 |            1            |
|    cp_info     |      constant_pool/常量池集合      | constant_pool_count - 1 |
|       u2       |       access_flags/访问标志        |            1            |
|       u2       |         this_class/类索引          |            1            |
|       u2       |        super_class/父类索引        |            1            |
|       u2       |      interface_count/接口计数      |            1            |
|       u2       |        interfaces/接口索引         |    interfaces_count     |
|       u2       |       fileds_count/字段计数        |            1            |
|   field_info   |         fields/字段表集合          |       field_count       |
|       u2       |       methods_count/方法计数       |            1            |
|       u2       |         methods/方法表集合         |      methods_count      |
|       u2       |     attributes_count/属性计数      |            1            |
| attribute_info |       attributes/属性表集合        |    attributes_count     |

其中：

* **魔数**（0xCAFEBABE）用于确定这个文件是否为一个能被虚拟机接受的Class文件。

* **主要版本号**限定了能够解析本Class文件的JDK版本。目前**最小版本号**主要用于标识实现了新特性的Class文件。

* **常量池计数**指定了常量池中有几项常量，从1开始计数，引用第0位表示在特定情况下“不引用任何一个常量池项目”。**常量池**主要存放两大类常量：字面量和符号引用。字面量好理解，就是不会改变的文字常量，比如文本字符串、被声明为final的常量值等。符号引用是用于唯一描述某个方法、字段、句柄、动态调用点等的文字。常量池集合中的每一项常量都是一张表，共有17种常量表。

* **访问标志**用于标识类或接口的访问信息，如public、final、abstract，或标识这个类是一个接口、注解、枚举、模块等，或标识这个类并非由用户代码产生。

* **类索引、父类索引和接口索引集合**这三项数据用于确定这个类的继承关系。类索引确定这个类的全限定名、父类索引确定这个类的父类的全限定类名，接口索引集合描述这个类实现了哪些接口。

* **字段表集合**用于描述接口或类中声明的变量。Java语言中的“字段”包括静态变量和实例变量，不包括局部变量。字段表中有这么几项数据：

  * |   类型    |                名称                 |       数量       |
    | :-------: | :---------------------------------: | :--------------: |
    |    u2     |      access_flags/字段访问标识      |        1         |
    |    u2     |   name_index/简单名称的常量池引用   |        1         |
    |    u2     | descriptor_index/描述符的常量池引用 |        1         |
    |    u2     |     attributes_count/属性表计数     |        1         |
    | attribute |        attributes/属性表集合        | attributes_count |

    字段可以包括的修饰符有字段的**作用域（public、private、protected修饰符**）、是实例变量还是类变量（**static修饰符**）、**可变性（final）**、**并发可见性（volatile修饰符，是否强制从主内存读写）**、**可否被序列化（transient修饰符**）、字段**数据类型（基本类型、对象、数组）**、**字段名称**。

  * 简单名称指的是没有类型和参数修饰的方法或者字段名称。

  * 描述符：基本数据类型（byte、char、double、float、int、long、short、boolean）以及代表无返回值的void**类型都用一个大写字符来表示**，而对象类型则用**字符L加对象的全限定名**来表示。用描述符来描述方法时，按照先参数列表、后返回值的顺序描述。

  * 属性表用于存储一些额外的信息。

* **方法表集合**用于描述方法。结构和字段表几乎完全一致。

  * 访问标志：因为`volatile`关键字和`transient`关键字不能修饰方法，所以方法表的访问标志中没有了` ACC_VOLATILE`标志和`ACC_TRANSIENT`标志。与之相对，`synchronized`、`native`、`strictfp`和`abstract` 关键字可以修饰方法，方法表的访问标志中也相应地增加了`ACC_SYNCHRONIZED`、` ACC_NATIVE`、`ACC_STRICTFP`和`ACC_ABSTRACT`标志。

* **属性表集合**：Class文件、字段表、方法表都可以携带自己的属性表（attribute_info）集合，以描述某些场景专有的信息。方法表集合之后的属性表集合，指的是class文件所携带的辅助信息，比如该class文件的源文件的名称，以及任何带有RetentionPolicy.CLASS或者RetentionPolicy.RUNTIME的注解。这类信息通常被用于Java虚拟机的验证与运行，以及Java程序的调试

</br>

</br>



### 3.2 解析class文件

Go语言内置了丰富的数据类型，非常适合处理class文件。

#### 3.2.1 读取数据

直接操作字节流很不方便，所以先定义一个结构体来帮助读取数据。

```go
type ClassReader struct {
	data []byte
}
func (self *ClassReader) readUint8() uint8 {}//读取u1
func (self *ClassReader) readUint16() uint16 {}//读取u2
func (self *ClassReader) readUint32() uint32 {}//读取u4
func (self *ClassReader) readUint64() uint64 {}//读取u8
func (self *ClassReader) readUint16s() []uint16 {}//读取u2集合
func (self *ClassReader) readBytes(n uint32) []byte {}//读取指定数量的字节
```

利用`encoding/binary`包下的`BigEndian.Uintxx`方法读取字节，用Go的reslice语法跳过已读字节。



#### 3.2.2 整体结构

有了ClassReader，就可以开始解析class文件了。在chap03/classfile目录下创建class_file.go文件，添加结构体：

```go
type ClassFile struct {
	//magic	uint32
	minorVersion	uint16	//最小版本号
	majorVersion	uint16	//主要版本号
	constantPool	ConstantPool	//常量池表
	accessFlags		uint16	//访问标志
	thisClass		uint16	//类索引
	superClass		uint16	//父类索引
	interfaceCount	uint16	//接口计数
	interfaces		[]uint16	//接口信息
	fields			[]*MemberInfo	//字段表
	methods			[]*MemberInfo	//方法表
	attributes		[]AttributeInfo	//属性表
}
```

相比Java语言，Go的访问控制比较简单：只有公开和私有两种，因此用首字母大写标记某个类型、结构体、字段、变量、函数、方法等是公开的。

首先定义一个函数Parse()把byte数组解析成ClassFile结构体。由于Go语言没有异常处理机制，只有一个panic-recover机制，所以在解析时使用panic-recover来捕获异常。

> 插播go的recover函数：Recover 是一个Go语言的内建函数，可以让进入宕机流程中的 goroutine 恢复过来，recover 仅在延迟函数 defer 中有效，在正常的执行过程中，调用 recover 会返回 nil 并且没有其他任何效果，如果当前的 goroutine 陷入恐慌，调用 recover 可以捕获到 panic 的输入值，并且恢复正常的执行。
>
> 通常来说，不应该对进入 panic 宕机的程序做任何处理，但有时，需要我们可以从宕机中恢复，至少我们可以在程序崩溃前，做一些操作，举个例子，当 web 服务器遇到不可预料的严重问题时，在崩溃前应该将所有的连接关闭，如果不做任何处理，会使得客户端一直处于等待状态，如果 web 服务器还在开发阶段，服务器甚至可以将异常信息反馈到客户端，帮助调试。
>
> 在其他语言里，宕机往往以异常的形式存在，底层抛出异常，上层逻辑通过 try/catch 机制捕获异常，没有被捕获的严重异常会导致宕机，捕获的异常可以被忽略，让代码继续运行。
>
> Go语言没有异常系统，其使用 panic 触发宕机类似于其他语言的抛出异常，recover 的宕机恢复机制就对应其他语言中的 try/catch 机制。

```go
func Parse(classData []byte) (cf *ClassFile, err error){
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	cr := &ClassReader{classData}
	cf = &ClassFile{}
	cf.read(cr)//将数据读取到结构体中
	return
}
```

再定义一个read函数用于将字符数组通过ClassReader读取到结构体中，具体操作就是为Class文件的每一项数据定义一个读取函数，在读取函数中再调用ClassReader的读取方法。有点类似于JavaWeb开发中，将业务层和Dao层分离的意思。

```go
func (self *ClassFile) read(reader *ClassReader) {
	self.readAndCheckMagic(reader)
	self.readAndCheckVersion(reader)
	self.constantPool = readConConstantPool(reader)
	self.accessFlags = reader.readUint16()
	self.thisClass = reader.readUint16()
	self.superClass = reader.readUint16()
	self.interfaces = reader.readUint16s()
	//以下三者都需要用到常量池中的字面量或者符号引用
	self.fields = readMember(reader, self.constantPool)
	self.methods = readMember(reader, self.constantPool)
	self.attributes = readAttributes(reader, self.constantPool)
}
```



#### 3.2.3魔数

class文件以“0xCAFEBABE”开头，用于标识这是一个class文件。所以需要首先读取魔数，校验这是否是个class文件。Java虚拟机规定，如果加载的class文件不符合要求的格式，JVM实现需要抛出`java.lang.ClassFormatError`异常。

在读取魔数的时候， 直接使用ClassReader读取连续的4个字节即可。

由于我们还未实现异常机制，因此如果不符合要求，暂时调用panic终止程序运行。



#### 3.2.3版本号

我们参考Java 8，支持的版本号为45.0~52.0的class文件。根据JVM规范规定，如果遇到其他版本号，应当抛出`java.lang.UnsupportedClassVersionError`异常。目前我们暂时使用panic终止程序。





#### 3.2.5类访问标志

版本号之后是常量池，比较复杂，放到后面再讲。先讲类访问标志，是16比特的bitmask。



#### 3.2.6类和超类索引

两个u2类型的常量池索引，分别指向类名和超类名。class文件存储的类名类似与完全限定名，只是把点换成了斜线，Java语言规范把这种名字叫做二进制名(binary names)。



#### 3.2.7 接口索引表

接口索引表中存放的也是常量池索引，给出该类实现的所有接口的名字。



#### 3.2.8 字段和方法表

字段表和方法表的基本结构大致相同，差别仅在于属性表。所以我们用一个结构体统一表示方法和字段。在chap03/classfile目录下创建member_info.go文件：

```go
type MemberInfo struct {
	cp	ConstantPool	//常量池引用
	accessFlags	uint16	//访问标志
	nameIndex	uint16	//简单名称索引
	describetorIndex	uint16	//描述符索引
	attributes	[]Attributeinfo	//属性表集合
}
```



编写方法按序读取字段表或方法表：

```go
//读取字段表或方法表
func readMembers(reader *ClassReader, cp ConstantPool) []*MemberInfo {
	memberCount := reader.readUint16()//读取计数
	members := make([]*MemberInfo, memberCount)
	for i := range members {
		members[i] = readMember(reader, cp)
	}
	return members
}

func readMember(reader *ClassReader, cp ConstantPool) *MemberInfo {
	return &MemberInfo{
		cp:	cp,
		accessFlags: reader.readUint16(),
		nameIndex:	reader.readUint16(),
		describetorIndex: reader.readUint16(),
		attributes: readAttributes(reader, cp),
	}
}
```





### 3.3 解析常量池



#### 3.3.1 ConstantPool结构体

在chap03/classfile目录下创建constant_pool.go文件：

```go
package classfile

type ConstantPool []ConstantInfo
```

常量池其实也是一张表，但是需要注意的是，表头的计数会比实际的表长大1，因为索引为0的位置需要留出来表示不引用任何常量。

CONSTANT_Long_info和CONSTANT_Double_info各占两个位置。也就是说当常量池中存在这两种常量时，实际的常量数量会少于n-1个。





#### 3.3.2 ConstantInfo接口

JVM规范中一共定义了14种常量，每种常量对应一个tag。

![image-20210119200136453](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210119200136453.png)

所以创建/classfile/constant_info.go接口文件并定义常量：

```go
const (
	CONSTANT_Class              = 7		//类或接口的符号引用
	CONSTANT_Fieldref           = 9		//字段的符号引用
	CONSTANT_Methodref          = 10	//类中方法的符号引用
	CONSTANT_InterfaceMethodref = 11	//接口中方法的符号引用
	CONSTANT_String             = 8		//字符串类型字面量
	CONSTANT_Integer            = 3		//整型字面量
	CONSTANT_Float              = 4		//浮点型字面量
	CONSTANT_Long               = 5		//长整型字面量
	CONSTANT_Double             = 6		//双精度字面量
	CONSTANT_NameAndType        = 12	//字段或方法的部分符号引用
	CONSTANT_Utf8               = 1		//UTF-8编码的字符串
	CONSTANT_MethodHandle       = 15	//表示方法句柄
	CONSTANT_MethodType         = 16	//表示方法类型
	CONSTANT_InvokeDynamic      = 18	//表示一个动态方法调用点
)
```

定义接口：

```go
type ConstantInfo interface {
	//读取常量信息，由具体的常量结构体实现
	readInfo(reader *ClassReader)
}
```

读取tag，根据tag来创建对应的ConstantXxxInfo常量，然后将数据读入具体的ConstantXxxInfo产量。



接下来逐个定义具体的常量结构体，并实现ConstantInfo接口

#### 3.3.3 CONSTANT_Integer_info

CONSTANT_Integer_info使用4字节存储整形常量。由于CONSTANT_Integer_info\CONSTANT_Float\CONSTANT_Long\CONSTANT_Double这四种数字型字面量十分相似，因此将这些结构体定义在同一个文件中。在chap03/classfile下创建目录cp_numeric.go

```go
type ConstantIntegerInfo struct {
	val int32
}
type ConstantFloatInfo struct {
	val float32
}
type ConstantLongInfo struct {
	val int64
}
type ConstantDoubleInfo struct {
	val float64
}
```

再分别实现readInfo方法，从字节流中读取数字型字面量。







#### 3.3.4 CONSTANT_Utf8_info

CONSTANT_Utf8_info常量中存放的是MUTF-8编码（Modified UTF-8）的字符串。

MUTF-8编码方式和UTF-8大致相同，但并不兼容。主要有以下两点区别：

* null字符（U+0000）会被编码成2字节：0xC0, 0x80
* 补充字符（U+FFFF的Unicode字符）按UTF-16拆分为代理对（Surrogate Pair）分别编码的。

以下内容摘自[JVM规范4.4.7](https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.4.7)：

*String content is encoded in modified UTF-8. Modified UTF-8 strings are encoded so that code point sequences that contain only non-null ASCII characters can be represented using only 1 byte per code point, but all code points in the Unicode codespace can be represented. Modified UTF-8 strings are not null-terminated.*

...

*There are two differences between this format and the "standard" UTF-8 format. First, the null character `(char)0` is encoded using the 2-byte format rather than the 1-byte format, so that modified UTF-8 strings never have embedded nulls. Second, only the 1-byte, 2-byte, and 3-byte formats of standard UTF-8 are used. The Java Virtual Machine does not recognize the four-byte format of standard UTF-8; it uses its own two-times-three-byte format instead.*

在chap03/classfile下创建cp_utf8.go文件：

```go
type ConstantUtf8Info struct {
	str string
}
```

由于go语言使用的是UTF-8编码，而字节码文件使用的是MUTF-8编码，因此需要进行解码。





#### 3.3.5 CONSTANT_String_info

CONSTANT_String_info常量表示java.lang.String字面量。CONSTANT_String_info本身不存放字符串数据，只存了常量池索引，这个索引指向一个CONSTANT_Utf-8_info常量。

在chpa03/classfile目录下创建cp_string.go 文件，在其中定义ConstantStringInfo结构体：

```go
type ConstantStringInfo struct {
	cp	ConstantPool
	stringIndex uint16
}
```



#### 3.3.6 CONSTANT_Class_info

CONSTANT_Class_info常量表示类或者接口的符号引用。和CONSTANT_String_info类似，name_index是常量池索引，指向CONSTANT_Utf-8_info常量。

在chap03/classfile目录下创建cp_class.go文件

```go
type ConstantClassInfo struct {
	cp	ConstantPool
	nameIndex	uint16
}
```



#### 3.3.7 CONSTANT_NameAndType_info

给出字段或方法的名称和描述符

CONSTANT_Class_Info和CONSTANT_NameAndType_info加在一起可以唯一确定一个字段或者方法。

描述符：基本数据类型（byte、char、double、float、int、long、short、boolean）以及代表无返回值的void**类型都用一个大写字符来表示**，而对象类型则用**字符L加对象的全限定名**来表示

![image-20201218093840471](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20201218093840471.png)

对于数组类型，每一维度将使用一个前置的“[”字符来描述，如一个定义为**“java.lang.String[][]”类型** 的二维数组将被记录成“**[[Ljava/lang/String；**”，一个整型数组“int[]”将被记录成“[I”。

用描述符来描述方法时，按照先参数列表、后返回值的顺序描述，参数列表按照参数的严格顺序放在一组小括号“()”之内。如方法`void inc()`的描述符为“`()V`”，方法`java.lang.String toString()`的描述符 为“`()Ljava/lang/String；`”，方法`int indexOf(char[]source，int sourceOffset，int sourceCount，char[]target， int targetOffset，int targetCount，int fromIndex)`的描述符为“`([CII[CIII)I`”。



#### 3.3.8 CONSTANT_Fieldref_info

表示字段符号引用

#### 3.3.9 CONSTANT_Methodref_info

标识普通方法符号引用

#### 3.3.10 CONSTANT_InterfaceMethodref_info

标识接口方法符号引用

以上三种类型常量结构相同，故定义一个统一的结构体ConstantMemberrefInfo来表示这3中常量。在chap03/classfile目录下创建cp_member_ref.go文件:

```go
type ConstantMemberrefInfo struct {
	cp ConstantPool
	classIndex	uint16
	nameAndTypeIndex	uint16
}
```

然后定义三个结构体“继承”ConstantMemberrefInfo。Go语言并没有“继承”这个概念，但是可以通过结构体嵌套来模拟：

```go
type ConstantFieldrefInfo struct {
	ConstantMemberrefInfo
}

type ConstantMethodrefInfo struct {
	ConstantMemberrefInfo
}

type ConstantInterfaceMethodrefInfo struct {
	ConstantMemberrefInfo
}
```



剩余三个常量：CONSTANT_MethodType_info、CONSTANT_MethodHandle_info和CONSTANT_InvokeDynamic_info在Java7才添加到class文件中，目的是支持新增的invokedynamic指令。



**小结**
可以把常量池中的常量分为两类：字面量（literal）和符号引用（symbolic reference）。字面量包括数字常量和字符串常量，符号引用包括类和接口名、字段和方法信息等。除了字面量，其他常量都是通过索引直接或间接指向CONSTANT_Utf-8_info常量。





### 3.4 解析属性表

Class文件、字段表、方法表都可以携带自己的属性表（attribute_info）集合，以描述某些场景专有的信息。

与Class文件中其他的数据项目要求严格的顺序、长度和内容不同，属性表集合的限制稍微宽松一 些，**不再要求各个属性表具有严格顺序**，并且《Java虚拟机规范》允许只要不与已有属性名重复，任何人实现的编译器都可以向属性表中写入自己定义的属性信息，Java虚拟机运行时会忽略掉它不认识的属性。

对于每一个属性，它的名称都要从常量池中引用一个CONSTANT_Utf8_info类型的常量来表示， 而属性值的结构则是完全自定义的，只需要通过一个u4的长度属性去说明属性值所占用的位数即可。

一个符合规则的属性表应该满足

![image-20201219102930077](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20201219102930077.png)



#### 3.4.1 AttributeInfo接口

和常量池类似，各种属性表达的信息也各不相同，因此无法用统一的结构来定义。

在chap03/classfile目录下创建attribute_info文件，在其中定义AttributeInfo接口：

```go
type AttributeInfo interface {
	readInfo(reader *ClassReader)
  func readAttributes(reader *ClassReader, cp ConstantPool) []AttributeInfo
  func newAttributeInfo(attrName string, attrLen uint32, cp ConstantPool) AttributeInfo
}
```

和ConstantInfo接口一样，AttributeInfo接口也只定义了一个readInfo方法，需要具体的属性实现。

再定义readAttribute和readAttributes方法，从class文件中读取属性集合，根据属性名创建属性实例。

JVM规范定义了23种属性，先解析以下8种：

```go
func newAttributeInfo(attrName string, attrLen uint32, cp ConstantPool) AttributeInfo {
	switch attrName {
	case "Code":
		return &CodeAttribute{cp : cp}
	case "ConstantValue":
		return &ConstantValueAttribute{}
	case "Deprecated":
		return &DeprecatedAttribute{}
	case "Exceptions":
		return &ExceptionsAttribute{}
	case "LineNumberTable":
		return &LineNumberTableAttribute{}
	case "LocalVariableTable":
		return &LocalVariableTableAttribute{}
	case "SourceFile":
		return &SourceFileAttribute{cp : cp}
	case "Systhetic":
		return &SyntheticAttribute{}
	default:
		return &UnparsedAttribute{attrName, attrLen, nil}
	}
}
```



按照用途，23种预定义的属性可以分为三组。第一组属性是实现Java虚拟机所必需的，共有5种；第二组属性是Java类库所必需的，共有12种；第三组属性主要提供给工具使用，共有6种。第三组是可选属性。

![image-20210224184558142](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210224184558142.png)



#### 3.4.2 Deprecated和Synthetic属性

Deprecated和Synthetic属性是最简单的两种属性，仅起标记作用，不包含任何数据。

正因为不包含任何数据，所以attribute_length的值必须是0。

Deprecated属性用于指出类、接口、字段或方法已经不建议使用，使用@Deprecated注解可以添加此属性。

Synthetic属性用来标记源文件中不存在、由编译器生成的类成员，引入Synthetic属性主要是为了支持嵌套类和嵌套接口。

```go
type DepreactedAttribute struct {
	MarkerAttribute
}

type SyntheticAttribute struct {
	MarkerAttribute
}
type MarkerAttribute struct {
}
```



#### 3.4.3 SourceFile属性

SourceFile是可选定长属性，只会出现在ClassFile结构中，用于指出源文件名。

```go
type SourceFileAttribute struct {
	cp ConstantPool
	sourceFileIndex	uint16
}
```

sourceFileIndex是一个指向CONSTANT_utf-8_info的索引，指出源文件名。



#### 3.4.4 ConstantValue属性

ConstantValue是定长属性，只会出现在field_info结构中，用于表示常量表达式的值。

```go
type ConstantValueAttribute struct {
	cosntantValueIndex uint16
}
```



#### 3.4.5 Code属性

Code是变长属性，用于存放字节码等方法相关信息。Code属性比较复杂。

```go
type CodeAttribute struct {
	cp	ConstantPool
	maxStack	uint16	//操作数栈的最大深度
	maxLocal	uint16	//局部变量表的大小
	code	[]byte	//字节码
	exceptionTable	[]*ExceptionTableEntry	//异常表
	attributes	[]AttributeInfo	//属性表
}

type ExceptionTableEntry struct {
	startPC	uint16
	endPC	uint16
	handlerPC	uint16
	catchType	uint16
}
```





#### 3.4.6 Exceptions属性

Exceptions是变长属性，记录方法抛出的异常表。

```go
type ExceptionsAttribute struct {
	exceptionIndexTable []uint16
}
```



#### 3.4.7 LineNumberTable和LocalVariableTable属性

LineNumberTable属性表存放方法的行号信息，LocalVariableTable属性表中存放方法的局部变量信息。

这两种属性和SourceFile属性都属于调试信息，都不是运行时所必须的，默认会生成。

```go
type LineNumberTableAttribute struct {
	lineNumberTable []*LineNumberTableEntry
}

type LineNumberTableEntry struct {
	startPC	uint16
	lineNumber	uint16
}
```

```go
type LocalVariableTableAttribute struct {
	localVariableTable	[]*LocalVariableTableEntry
}

type LocalVariableTableEntry struct {
	startPC uint16
	length uint16
	nameIndex uint16
	descriptorIndex uint16
	index uint16
}
```



### 3.5 测试本章代码

修改main.go文件：

```go
func startJVM(cmd *Cmd){
	clsp := classpath.Parse(cmd.Xjre, cmd.clspath)
	className := strings.Replace(cmd.class,".","/",-1)
	cf := loadClass(className, clsp)
	fmt.Println(cmd.class)
	printClassInfo(cf)
}
//读取Class文件，将属性load到ClassFile对象中
func loadClass(className string, cp *classpath.Classpath) *classfile.ClassFile {
	classData, _, err := cp.ReadClass(className)//读数据
	if err != nil {
		panic(err)
	}
	cf, err := classfile.Parse(classData)//解析Class文件
	if err != nil {
		panic(err)
	}
	return cf
}
//打印ClassFile对象
func printClassInfo(cf *classfile.ClassFile) {
	fmt.Printf("version: %v.%v\n", cf.MajorVersion(), cf.MinorVersion())
	fmt.Printf("constants count: %v\n", len(cf.ConstantPool()))
	fmt.Printf("access flags: 0x%x\n", cf.AccessFlags())
	fmt.Printf("this class: %v\n", cf.ClassName())
	fmt.Printf("super class: %v\n", cf.SuperClassName())
	fmt.Printf("interfaces: %v\n", cf.InterfaceNames())
	fmt.Printf("fields count: %v\n", len(cf.Fields()))
	for _, f := range cf.Fields() {
		fmt.Printf("  %s\n", f.Name())
	}
	fmt.Printf("methods count: %v\n", len(cf.Methods()))
	for _, m := range cf.Methods() {
		fmt.Printf("  %s\n", m.Name())
	}
}

```

执行命令：

```bash
go install ../src/jvmgo/chap03
chap03 -Xjre /Library/Java/JavaVirtualMachines/jdk1.8.0_201.jdk/Contents/Home/jre java.lang.String
```

结果如下:

![image-20210224201159809](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210224201159809.png)

Class文件解析成功。





## 4 运行时数据区

多线程共享的内存区域主要存放两类数据：类数据和类实例。对象数据存放在Heap中，类数据存放在方法区当中。

线程私有的内存区域用于辅助执行Java字节码。

Go本身有垃圾回收功能，所以可以直接使用Go的堆和垃圾收集器。



### 4.1 数据类型

Go语言提供了非常丰富的数据类型，包括各种整数和两种精度的浮点数。Java和Go的浮点数都采用IEEE 754规范。对于基本类型，可以直接在Go和Java之间建立映射关系；对于引用类型自然选择使用指针实现。

首先在chap04/rtda/Object.go中定义一个结构体来表示Object对象：

```go
type Object struct {
	//todo
}
```



### 4.2 实现运行时数据区

本节会先实现线程私有的运行时数据区。下面先从线程入手。



#### 4.2.1 线程

在chap04/rtda下创建thread.go文件，在其中定义Thread结构体：

```go
type Thread struct {
	pc	int
	stack *Stack//虚拟机栈
}
```

线程能够操作PC、操作虚拟机栈的栈帧，因此需要定义相关方法：

```go
func(self *Thread) PC() int {
	return self.pc
}
func(self *Thread) SetPC(pc int) {
	self.pc = pc
}
func(self *Thread) PushFrame(frame *Frame){
	self.stack.push(frame)
}
func(self *Thread) PopFrame(frame *Frame){
	self.stack.pop(frame)
}
func(self *Thread) CurrentFrame() *Frame{
	return self.stack.top()
}
```





#### 4.2.2 Java虚拟机栈

Java虚拟机规范对Java虚拟机栈的约束非常宽松。

每个JVM线程都有一个JVM栈，在线程创建的时候，随之创建。JVM栈与常规语言（如C）中的栈非常类似，它保存局部变量以及中间运算结果，参与方法的调用和返回。JVM栈的内存分配无需连续。

- If the computation in a thread requires a larger Java Virtual Machine stack than is permitted, the Java Virtual Machine throws a `StackOverflowError`.
- If Java Virtual Machine stacks can be dynamically expanded, and expansion is attempted but insufficient memory can be made available to effect the expansion, or if insufficient memory can be made available to create the initial Java Virtual Machine stack for a new thread, the Java Virtual Machine throws an `OutOfMemoryError`.

我们用链表数据结构来实现Java虚拟机栈，这样栈就可以按需使用内存空间，而且弹出的栈帧也可以及时地被Go的垃圾收集器回收。在chap04/rtda目录下创建jvm_stack.go文件，在其中定义Stack结构体。

```go
type Stack struct {
	maxSize	uint
	size	uint
	_top	*Frame
}
```

同时在frame.go文件中编写栈帧的结构体：

```go
type Frame struct {
	lower	*Frame	//指向此栈帧下方的第一个栈帧
	localVars	LocalVars
	operandStack	*OperandStack
  thread       *Thread
	nextPC       int // the next instruction after the call
}
```

我们在栈帧中定义了5个属性：

* lower作为链表的next
* localVars和operandStack分别表示局部变量表和操作数栈
* thread表示当前所在线程
* nextPC表示下一条指令的地址



#### 4.2.3 局部变量表

局部变量表是按索引访问的，可以设置为一个数组。根据JVM规范， 这个数组的每个槽至少可以容纳一个int或者reference值，两个连续的元素可以容纳一个long或double值。

在Go中，最容易想到用`[]int`来表示这个数组。Go的int类型因平台而异，在64位系统上是int64，在32位系统上是int32，总之足够容纳Java的int类型，另外它和内置的uintptr类型宽度一样，所以也足够放下一个内存地址。但Go的垃圾回收机制并不能有效处理uintptr指针。也就是说，如果一个结构体实例，除了uintptr类型指针保存它的地址之外，其他地方没有引用这个实例，它就会被当作垃圾回收。

另一个方案是用[]interface{}类型，但是可读性比较差。

第三种方案是定义一个结构体，让它可以同时容纳一个int值和一个引用值。在chap/rtda目录下创建slot.go文件

```go
type Slot struct {
	num int32
	ref *Object
}
```

实现局部变量表

```go
type LocalVars []Slot

func newLocalVars(maxLocals uint) LocalVars {
	if maxLocals > 0 {
		return make([]Slot, maxLocals)
	}
	return nil
}
```

同时定义一些存取变量的方法

```go
//一些存取变量的方法
func (self LocalVars) SetInt(index uint, val int32)

func (self LocalVars) GetInt(index uint) int32
//float变量先转成int类型，然后按int变量来处理
func (self LocalVars) SetFloat(index uint, val float32) 
func (self LocalVars) GetFloat(index uint) float32 
// long consumes two slots
func (self LocalVars) SetLong(index uint, val int64) 
func (self LocalVars) GetLong(index uint) int64 

// double consumes two slots
func (self LocalVars) SetDouble(index uint, val float64)
func (self LocalVars) GetDouble(index uint) float64 

func (self LocalVars) SetRef(index uint, ref *Object) 
func (self LocalVars) GetRef(index uint) *Object
```





#### 4.2.4 操作数栈

在/chap04/rtda/operand_stack.go文件中定义OperandStack结构体：

```go
type OperandStack struct {
	size int
	slots []Slot
}
func newOperandStack(maxStack uint) *OperandStack{
	if maxStack > 0 {
		return &OperandStack{
			slots: make([]Slot, maxStack),
		}
	}
	return nil
}
```

操作数栈的大小是编译期就已经确定的，所以可以用[]Slot实现。size字段用于记录栈顶位置。

也定义一些操作数据的方法





### 4.3 测试本章代码

修改main.go:

```go
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

```

运行命令：

```shell
go install chap04
chap04 test
```

结果：

![image-20210225195538426](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210225195538426.png)



## 5 指令集和解释器

本章将在前两章的基础上编写一个简单的解释器，并实现大约150条指令。



### 5.1 字节码和指令集

每个类或接口都会被Java编译器编译成一个class文件，类或接口的方法信息就放在class文件的method_info结构中。如果方法不是抽象的，也不是本地方法，方法的Java代码就会被编译器编译成字节码，存放在method_info结构的Code属性中。

字节码中存放编码后的Java虚拟机指令。每条指令都以一个单字节的操作码（opcode）开头，这就是字节码名称的由来。

Java虚拟机使用的是变长指令，操作码后面可以跟零字节或多字节的操作数。为了让编码后的字节更加紧凑，很多操作码本身就隐含了操作数，比如把常数0推入操作数栈的指令是iconst_0。



由于操作数栈和局部变量表只存放数据的值，并不记录数据类型，所以指令必须知道自己在操作什么类型的数据。这一点也直接反映在了操作码的助记符上。例如iadd指令就是对int值进行加法操作；dstore指令把操作数栈顶的double值弹出，存储到局部变量表中。areturn从方法中返回引用值。

总结其规律：**如果某类指令可以操作不同类型的变量，则助记符的第一个字母表示变量类型。**

![image-20210226181936799](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210226181936799.png)



JVM规范把已经定义的205条指令按用途分成了11类，分别是：

1. 常量（constants）指令
2. 加载（loads）指令
3. 存储（stores）指令
4. 操作数栈（stack）指令
5. 数学（math）指令
6. 转换（conversions）指令
7. 比较（comparisons）指令
8. 控制（control）指令
9. 引用（references）指令
10. 扩展（extended）指令
11. 保留（reserved）指令

保留指令一共有三条：breakpoint（0xCA）、impdep1（0xFE）、impdep2（0xFF）。这三条指令不允许出现在class文件中。



### 5.2 指令和指令解码



JVM规范给出了JVM解释器的大致逻辑：

```java
do {
    atomically calculate pc and fetch opcode at pc;
    if (operands) fetch operands;
    execute the action for the opcode;
} while (there is more to do);
```

每次循环重复做三件事：

1. 自动计算PC，根据PC值取操作码
2. 如果存在操作数，就取操作数
3. 执行指令



如果用for循环+switch case的方法来实现解释器，那么代码的可读性将非常之差，而且很不优雅。所以我们把指令抽象成接口，解码和执行逻辑写在具体的指令实现中，如：

```go
for{
  pc := calculatePC()
  opcode := bytecode[pc]
  inst := createInst(opcode)
  inst.fetchOperands(bytecode)
  inst.execute()
}
```



#### 5.2.1 Instructions接口

```go
type Instruction interface {
	FetchOperands(reader *BytecodeReader)//从字节码中提取操作数
	Execute(frame *rtda.Frame)//执行指令逻辑
}
```

按照操作数类型定义一些结构体，并实现FetchOperands方法。这些结构体相当于Java中的抽象类，实现了Instruction接口，并规范了取指令方法。具体的指令继承这些结构体，然后专注实现Execute()方法即可。

```go
//零操作数的指令
type NoOperandsInstruction struct {}

func (self *NoOperandsInstruction) FetchOperands(reader *BytecodeReader) {
	//nothing to do
}

//跳转指令
type BranchInstruction struct {
	Offset int//跳转偏移量
}
func (self *BranchInstruction) FetchOperands(reader *BytecodeReader) {
	//nothing to do
	self.Offset = int(reader.ReadInt16())
}

//存储和加载类指令需要根据索引存取局部变量表，索引由单字节操作数给出，所以把这类指令抽象成Index8Instruction
type Index8Instruction struct {
	Index uint//局部变量表的索引
}
func (self *Index8Instruction) FetchOperands(reader *BytecodeReader){
	self.Index = uint(reader.ReadUint8())
}
//有一些指令需要访问常量池，常量池索引由两字节操作数给出，所以把这类指令抽象成Index16Instruction
type Index16Instruction struct {
	Index uint//局部变量表的索引
}
func (self *Index16Instruction) FetchOperands(reader *BytecodeReader){
	self.Index = uint(reader.ReadUint16())
}
```



#### 5.2.2 BytecodeReader

在base目录下创建bytecode_reader.go文件，在其中定义BytecodeReader结构体：

```go
//字节码读取器
type BytecodeReader struct {
	code []byte//存放字节码
	pc int//记录读取到了哪个字节
}
```

为了避免每次解码指令都新创建一个BytecodeReader实例，给它定义一个Reset()方法：

```go
func (self *BytecodeReader) Reset(code []byte, pc int){
	self.code = code
	self.pc = pc
}
```

再定义一些读取字节码的方法：

```go
//读8比特，也就是一个字节
func (self *BytecodeReader) ReadUint8() uint8
func (self *BytecodeReader) Readint8() int8
//连续读取两字节
func (self *BytecodeReader) ReadUint16() uint16
func (self *BytecodeReader) ReadInt16() int16
//连续读取四字节
func (self *BytecodeReader) ReadUint32() int32
func (self *BytecodeReader) ReadInt32s() []int32
func (self *BytecodeReader) SkipPadding()
```



### 5.3 常量指令

常量指令把常量推入操作数栈顶。常量可以来自三个地方：隐含在操作码里、操作数和运行时常量池。

常量指令共有21条，本节实现其中的18条，另外三条是ldc指令，用于从运行时常量池中加载常量，将在第6章介绍。

#### 5.3.1 nop指令

即使是在JVM规范上，也是除了donothing之外，没有别的介绍了。

```go
type NOP struct {
	base.NoOperandsInstruction
}

func (self *NOP) Excute(frame *rtda.Frame){
	//do nothing
}
```



#### 5.3.2 const系列指令

这一系列指令把隐含在操作码中的常量值推入操作数栈顶。

创建constants/const.go文件：

```go
type ACONST_NULL struct { base.NoOperandsInstruction }	//将nil入栈
type DCONST_0 struct { base.NoOperandsInstruction }	//将double 0入栈
type DCONST_1 struct { base.NoOperandsInstruction }	
type FCONST_0 struct { base.NoOperandsInstruction }	//将float 0入栈
type FCONST_1 struct { base.NoOperandsInstruction }
type ICONST_M1 struct { base.NoOperandsInstruction }	//将int型-1入栈
type ICONST_0 struct { base.NoOperandsInstruction }	//将int 0入栈
type ICONST_1 struct { base.NoOperandsInstruction }
type ICONST_2 struct { base.NoOperandsInstruction }
type ICONST_3 struct { base.NoOperandsInstruction }
type ICONST_4 struct { base.NoOperandsInstruction }
type ICONST_5 struct { base.NoOperandsInstruction }
type LCONST_0 struct { base.NoOperandsInstruction }	//将long 0入栈
type LCONST_1 struct { base.NoOperandsInstruction }
```

分别实现其Execute方法：

```go
func (self *ACONST_NULL) Execute(frame *rtda.Frame){
	frame.OperandStack().PushRef(nil)
}
func (self *DCONST_0) Execute(frame *rtda.Frame){
	frame.OperandStack().PushDouble(float64(0))
}
func (self *DCONST_1) Execute(frame *rtda.Frame){
	frame.OperandStack().PushDouble(float64(1))
}
...
```



#### 5.3.3 bipush和sipush指令

bipush指令从操作数中获取一个byte型整数，并将其扩展成int型，然后推入栈顶。

sipush指令从操作数中获取一个short型整数，并将其扩展成int型，然后推入栈顶。

创建constants/ipush.go文件：

```go
//push byte
type BIPUSH struct {
	val int8
}
//push short
type SIPUSH struct {
	val int16
}
```

分别实现bipush和sipush的FetchOperands以及Execute方法。FetchOperands从字节码中读取规定位数的操作数，Execute将操作数转为int，并推入栈顶。





### 5.4 加载指令

加载指令从局部变量表获取变量，然后推入操作数栈顶。加载指令共33条，按照所操作变量的类型可分为6类：

1. aload系列指令操作引用类型变量
2. dload系列指令操作double类型变量
3. fload系列指令操作float变量
4. iload系列指令操作int变量
5. lload系列指令操作long变量
6. xaload操作数组

创建loads/iload.go文件：

```go
type ILOAD struct { base.Index8Instruction }
type ILOAD_0 struct { base.NoOperandsInstruction }
type ILOAD_1 struct { base.NoOperandsInstruction }
type ILOAD_2 struct { base.NoOperandsInstruction }
type ILOAD_3 struct { base.NoOperandsInstruction }
```

对于ILOAD指令：

The *index* is an unsigned byte that must be an index into the local variable array of the current frame ([§2.6](https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-2.html#jvms-2.6)). The local variable at *index* must contain an `int`. **The *value* of the local variable at *index* is pushed onto the operand stack.**

对于ILOAD_n指令：

Load `int` from local variable

The <*n*> must be an index into the local variable array of the current frame ([§2.6](https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-2.html#jvms-2.6)). The local variable at <*n*> must contain an `int`. The *value* of the local variable at <*n*> is pushed onto the operand stack.操作数的索引来自操作码。

其余4条指令类似，xaload放到后面实现。



### 5.5 存储指令

和加载指令正好相反，存储指令把变量从操作数栈顶弹出，然后存入局部变量表。

和加载指令一样，存储指令也可以分为6类。

1. astore系列指令操作引用类型变量
2. dstore系列指令操作double类型变量
3. fstore系列指令操作float变量
4. istore系列指令操作int变量
5. lstore系列指令操作long变量
6. xastore操作数组(bastore, castore)

创建store/lstore.go文件：

```go
type LSTORE struct { base.Index8Instruction }
type LSTORE_0 struct { base.NoOperandsInstruction }
type LSTORE_1 struct { base.NoOperandsInstruction }
type LSTORE_2 struct { base.NoOperandsInstruction }
type LSTORE_3 struct { base.NoOperandsInstruction }
```

The *index* is an unsigned byte. Both *index* and *index*+1 must be indices into the local variable array of the current frame ([§2.6](https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-2.html#jvms-2.6)). The *value* on the top of the operand stack must be of type `long`. It is popped from the operand stack, and the local variables at *index* and *index*+1 are set to *value*.

The *lstore* opcode can be used in conjunction with the *wide* instruction ([§*wide*](https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-6.html#jvms-6.5.wide)) to access a local variable using a two-byte unsigned index.

其余4条指令类似，xastore放到后面实现。



### 5.6 栈指令

栈指令直接对操作数栈进行操作，共9条指令：

1. pop和pop2指令将栈顶变量弹出；
2. dup系列指令复制栈顶变量；
3. swap指令交换栈顶两个变量。

和其他类型的指令不同，栈指令并不关心变量类型。在操作数栈的数据结构中再封装两个方法：PushSlot和PopSlot，用于操作栈的元素。

```go
func (self *OperandStack) PushSlot(slot Slot){
	self.slots[self.size] = slot
	self.size++
}

func (self *OperandStack) PopSlot() Slot{
	self.size--
	return self.slots[self.size]
}
```



#### 5.6.1 pop指令

关于pop指令：

Pop the top value from the operand stack.

The *pop* instruction must not be used unless *value* is a value of a category 1 computational type ([§2.11.1](https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-2.html#jvms-2.11.1)).

除非value是1类计算类型的值，否则不得使用pop指令

pop2用于弹出double和long变量。

```go
type POP struct { base.NoOperandsInstruction }
type POP2 struct { base.NoOperandsInstruction }
```



#### 5.6.2 dup指令

dup：Duplicate the top operand stack value.

```markdown
bottom -> top
[...][c][b][a]
             \_
               |
               V
[...][c][b][a][a]
```

dup_x1：Duplicate the top operand stack value and insert two values down

```markdown
bottom -> top
[...][c][b][a]
          __/
         |
         V
[...][c][a][b][a]
```

dup_x2：Duplicate the top operand stack value and insert two or three values down

```markdown
bottom -> top
[...][c][b][a]
       _____/
      |
      V
[...][a][c][b][a]
```

dup2：Duplicate the top one or two operand stack values

```markdown
bottom -> top
[...][c][b][a]____
          \____   |
               |  |
               V  V
[...][c][b][a][b][a]
```

dup2_x1：Duplicate the top one or two operand stack values and insert two or three values down

```markdown
bottom -> top
[...][c][b][a]
       _/ __/
      |  |
      V  V
[...][b][a][c][b][a]
```

dup2_x2：Duplicate the top one or two operand stack values and insert two, three, or four values down

```markdown
bottom -> top
[...][d][c][b][a]
       ____/ __/
      |   __/
      V  V
[...][b][a][d][c][b][a]
```



```go
type DUP struct { base.NoOperandsInstruction }
type DUP_X1 struct { base.NoOperandsInstruction }
type DUP_X2 struct { base.NoOperandsInstruction }
type DUP2 struct { base.NoOperandsInstruction }
type DUP2_X1 struct { base.NoOperandsInstruction }
type DUP2_X2 struct { base.NoOperandsInstruction }
```



#### 5.6.3 swap指令

Swap the top two operand stack values

```go
type SWAP struct {
	base.NoOperandsInstruction
}
```



### 5.7 数学指令

数学指令大致对应Java语言中的加、减、乘、除等数学运算符。数学指令包括算数指令、位移指令和布尔运算指令，共37条。



#### 5.7.1 算数指令

算数指令又可以进一步分为加法指令（add）、减法指令（sub）、乘法指令（mul）、除法指令（div）、求余指令（rem）和取反指令（neg）6种。

求余指令：

```go
type DREM struct { base.NoOperandsInstruction }//double求余
type FREM struct { base.NoOperandsInstruction }//float求余
type IREM struct { base.NoOperandsInstruction }//int求余
type LREM struct { base.NoOperandsInstruction }//long求余
```

其余5条指令都比较简单。



#### 5.7.2 位移指令

位移指令可以分为左移和右移，右移指令又可以分为算数右移（右符号右移）和逻辑右移（无符号右移）两种。

```go
type ISHL struct { base.NoOperandsInstruction }
type ISHR struct { base.NoOperandsInstruction }
type IUSHL struct { base.NoOperandsInstruction }
type LSHL struct { base.NoOperandsInstruction }
type LSHR struct { base.NoOperandsInstruction }
type LUSHL struct { base.NoOperandsInstruction }
```

```go
func (self *ISHL) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v2 := stack.PopInt()	//v2为需要左移的位数
	v1 := stack.PopInt()
	s := uint32(v2) & 0x1f	//只需要取v2的前5比特就足够表示位移位数了
	result := v1 << s
	stack.PushInt(result)
}

```

1. int变量只有32位，所以只取v2的前5比特就足够表示位移位数了；

2. Go语言位移操作符右侧必须是无符号整数，所以需要对v2进行类型转换。



#### 5.7.3 布尔运算指令

布尔运算指令只能操作int和long变量，分为按位与（and）、按位或（or）、按位异或（xor）三种。以与为例：

```go
type IAND struct { base.NoOperandsInstruction }
type LAND struct { base.NoOperandsInstruction }

func (self *IAND) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v1 := stack.PopInt()
	v2 := stack.PopInt()
	result := v1 & v2
	stack.PushInt(result)
}

func (self *LAND) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v1 := stack.PopLong()
	v2 := stack.PopLong()
	result := v1 & v2
	stack.PushLong(result)
}
```

代码较为简单，不多解释了。



#### 5.7.4 iinc指令

iinc指令给局部变量表中的int变量增加常量值，局部变量表索引和常量值都由**指令的操作数提供。**

```go
type IINC struct {
	Index uint
	Const int32
}

func (self *IINC) FetchOperands(reader *base.BytecodeReader){
	self.Index = uint(reader.ReadUint8())
	self.Const = int32(reader.ReadInt8())
}

func (self *IINC) Execute(frame *rtda.Frame){
	localVars := frame.LocalVars()
	val := localVars.GetInt(self.Index)
	val += self.Const
	localVars.SetInt(self.Index, val)
}
```





### 5.8 类型转换指令

类型转换指令大致对应Java语言中的基本类型强制转换操作。类型转换指令共有15条。

按照被转换变量的类型，类型转换指令可以分为3种：

1. i2x系列指令把int变量强转为其他类型。
2. l2x系列指令把long变量强转为其他类型。
3. f2x系列指令把float变量强转为其他类型。
4. d2x系列指令把double变量强转为其他类型。

```go
type D2F struct { base.NoOperandsInstruction }
type D2I struct { base.NoOperandsInstruction }
type D2L struct { base.NoOperandsInstruction }

func (self *D2F) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	d := stack.PopDouble()
	f := float32(d)
	stack.PushFloat(f)
}

func (self *D2I) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	d := stack.PopDouble()
	i := int32(d)
	stack.PushInt(i)
}

func (self *D2L) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	d := stack.PopDouble()
	l := int64(d)
	stack.PushLong(l)
}
```

Go基本类型转换非常方便，因此有利于实现上述指令。



### 5.9 比较指令

比较指令可以分为两类：一类将比较结果推入操作数栈顶，一类根据比较结果跳转。

比较指令是编译器实现if-else, for, while等语句的基础，共有19条指令。



#### 5.9.1 lcmp指令

lcmp指令用于比较long变量。

```go
func (self *LCMP) Execute(frame *rtda.Frame){
	stack := frame.OperandStack()
	v1 := stack.PopLong()
	v2 := stack.PopLong()
	i := int32(0)
	if v1 > v2 {
		i = 1
	}else if v1 < v2{
		i = -1
	}
	stack.PushInt(i)
}
```



#### 5.9.2 fcmp<op>和dcmp<op>指令

fcmpg和fcmpl指令用于比较float变量。这两条指令和lcmp类似，但是出了比较的变量类型不同以外，还有一个重要的区别。由于浮点数计算有可能产生NaN值，所以除了大于、小于、等于之外，还有第四种结果：无法比较。fcmpg和fcmpl指令的区别就在于对第四种结果的定义。当两个float变量中至少有一个NaN时，用fcmpg指令比较的结果是1，用fcmpl指令比较的结果是-1。

```go
type FCMPG struct {base.NoOperandsInstruction}
type FCMPL struct {base.NoOperandsInstruction}

func _fcmp(frame *rtda.Frame, gFlag bool){
	stack := frame.OperandStack()
	v2 := stack.PopFloat()
	v1 := stack.PopFloat()
	if v1 > v2 {
		stack.PushInt(1)
	}else if v1 == v2 {
		stack.PushInt(0)
	}else if v1 < v2 {
		stack.PushInt(-1)
	}else if gFlag{
		stack.PushInt(1)
	}else{
		stack.PushInt(-1)
	}
}

func (self *FCMPG) Execute(frame *rtda.Frame){
	_fcmp(frame, true)
}
func (self *FCMPG) Execute(frame *rtda.Frame){
	_fcmp(frame, false)
}
```



dcmp<op>与fcmp<op>代码几乎相同。



#### 5.9.3 if<cond>指令

if<cond>指令把操作数栈顶的int变量弹出，然后跟0比较，满足条件就跳转。假设弹出的变量是x

* ifeq : x == 0 (equal)
* ifne : x != 0 (not equal)
* iflt : x < 0 (less than)
* ifle : x <= 0 (less equal than)
* ifgt : x > 0 (greater than)
* ifge : x >= 0 (greater equal than)



#### 5.9.5 if_acmp<cond>指令

if_acmpeq和if_acmpne指令把栈顶的两个引用弹出，根据引用是否相同进行跳转。

if_icmp<cond>指令把栈顶的两个int变量弹出，然后进行比较，满足条件跳转。





### 5.10 控制指令

控制指令共有11条。jsr和ret指令在Java6之后不再使用了；return系列有6条指令，放到后面实现。剩下goto, tableswitch和lookupswitch三条指令。

#### 5.10.1 goto指令

goto指令进行无条件跳转。



#### 5.10.2 tableswitch指令

Java语言中的switch-case语句有两种实现方式：如果case值可以编码成一个索引表，则实现成tableswitch指令；否则实现成lookupswitch指令。

可用tableswitch实现：

![image-20210302203905279](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210302203905279.png)



需要用lookupswitch实现：

![image-20210302203935136](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210302203935136.png)



```go
type TABLE_SWITCH struct {
	defaultOffset int32	//默认情况下执行跳转所需的字节码偏移量
	low int32	//记录case的取值下限
	high int32	//记录case的取值上限
	jumpOffsets []int32	//一个索引表，存放high - low + 1个int值，对应各种case下，执行跳转所需的字节码偏移量
}
```



#### 5.10.3 lookupswitch指令

```go
type LOOKUP_SWITCH struct {
   defaultOffset int32
   npairs int32
   matchOffsets []int32
}
```

matchOffsets有点像Map，它的key是case值，value是跳转偏移量。所以查找时，需要遍历matchOffsets。





### 5.11 扩展指令

扩展指令共有6条，jsr_w已经不再使用。multianewarray指令用于创建多维数组，放到后面实现。



#### 5.11.1 wide指令

加载类指令、存储类指令、ret指令和iinc指令需要按索引访问局部变量表，索引以uint8的形式存在字节码中。对于大部分方法来说，局部变量表大小都不会超过256，所以用一字节来表示索引足矣。但是如果有方法的局部变量表超过了256，JVM规范定义了wide指令来扩展前述指令。

```go
type WIDE struct {
	modifiedInstruction base.Instruction
}
```

wide指令改变其他指令的行为，modifiedInstruction字段存放被改变的指令。wide指令需要自己解码出modifiedInstruction。

```go
func (self *WIDE) FetchOperands(reader *base.BytecodeReader) {
	opcode := reader.ReadUint8()
	switch opcode {//根据操作码来解码指令
	case 0x15:
		inst := &loads.ILOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x16:
		inst := &loads.LLOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x17:
		inst := &loads.FLOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x18:
		inst := &loads.DLOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x19:
		inst := &loads.ALOAD{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x36:
		inst := &stores.ISTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x37:
		inst := &stores.LSTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x38:
		inst := &stores.FSTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x39:
		inst := &stores.DSTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x3a:
		inst := &stores.ASTORE{}
		inst.Index = uint(reader.ReadUint16())
		self.modifiedInstruction = inst
	case 0x84:
		inst := &math.IINC{}
		inst.Index = uint(reader.ReadUint16())
		inst.Const = int32(reader.ReadInt16())
		self.modifiedInstruction = inst
	case 0xa9: // ret
		panic("Unsupported opcode: 0xa9!")
	}
}
```



#### 5.11.2 ifnull和ifnonnull指令

根据引用是否为空进行跳转。





#### 5.11.3 goto_w指令

goto_w指令和goto指令的唯一区别就是索引从2字节变成了4字节。







### 5.12 解释器

实现一个简单的解释器。

```go
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
```



### 5.13 测试

修改main方法：

```go
package main

import (
	"fmt"
	"jvmgo/chap05/classfile"
	"jvmgo/chap05/classpath"
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
		if m.Name() == "main" && m.Descriptor() == "(Ljava/lang/String)V"{
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


```

执行流程：

* 首先调用loadClass()方法读取并解析class文件
* 再调用getMainMethod函数查找类的main方法
* 最后调用interpret函数解析执行main方法

创建测试文件GaussTest.java，用Javac编译器编译成class文件

```java
public class GaussTest{

	public static void main(String[] args){
		int sum = 0;
		for(int i = 1; i <= 100; ++i){
			sum += i;
		}
		System.out.println(sum);
	}
}
```

执行命令：

```shell
go install chap05
```

再转到bin目录下执行以下命令：

```shell
chap05 -Xjre /Library/Java/JavaVirtualMachines/jdk1.8.0_201.jdk/Contents/Home/jre -cp /Users/huangyucai/Documents/code/git_depositorys/github_KolinHuang/JVM-in-Go-Style/javafiles/ GaussTest
```

运行结果：

![image-20210303192941959](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210303192941959.png)

在局部变量表中得到了sum = 5050的结果，测试成功。

至此分析一波方法执行的流程：

* 首先将命令行参数解析到cmd对象中，然后调用startJVM开始运行JVM；

* 调用classPath.Parse()方法解析类路径：

  * 使用-Xjre选项指定的类路径来解析启动类路径和扩展类路径；
  * 使用-clspath选项指定的类路径来解析用户类路径。

* 调用loadClass方法加载类：

  * 调用cp.ReadClass方法读取主类文件

    * 优先从启动类路径下查找类文件：self.bootClasspath.readClass(className)
    * 其次在扩展类路径下查找类文件：self.extClasspath.readClass(className)
    * 最后在用户类路径下查找类文件：self.userClasspath.readClass(className)
    * 最终在用户类路径下查找到类文件，然后调用ioutil包下的ReadFile读取类文件，并返回字符数组。

  * 调用classfile.Parse方法对该类文件的字符数组进行解析，严格按JVM规范将所有数据填入ClassFile结构体中：

    * 创建了一个类文件读取器ClassReader来辅助读取类文件，在ClassReader中定义了许多读取类文件的方法，可以按比特数量读取字符流。

    * 调用了ClassFile的read方法来解析字符数组：

      * ```go
        self.readAndCheckMagic(reader)//读魔数
        self.readAndCheckVersion(reader)//检查版本号
        self.constantPool = readConstantPool(reader)//读常量池
        self.accessFlags = reader.readUint16()
        self.thisClass = reader.readUint16()
        self.superClass = reader.readUint16()
        self.interfaces = reader.readUint16s()
        //以下三者都需要用到常量池中的字面量或者符号引用
        self.fields = readMembers(reader, self.constantPool)
        self.methods = readMembers(reader, self.constantPool)
        self.attributes = readAttributes(reader, self.constantPool)
        ```

* 调用getMainMethod获取Main方法

* 调用interpret方法来解释运行Main方法

  * 获取方法的Code属性、最大栈深度、最大局部变量表长度以及字节码。

  * 创建一个线程，再为线程创建Java虚拟机栈帧，根据最大栈深度、最大局部变量表长度创建操作数栈和局部变量表。

  * 将栈帧压入Java虚拟机栈中。

  * 运行loop方法：

    * 将栈帧弹出，创建一个字节码读取器BytecodeReader，循环读取指令并执行：

      * 取PC值，并设置到PC计数器中。
      * 根据PC值取操作码opcode，第一个操作码为3，对应了指令iconst_0，将int型常量0压入操作数栈中。
      * 根据opcode创建指令实例。
      * 根据指令实例取操作数。
      * 执行指令。
      * 重复上述步骤。

    * 

      

![image-20210305181242965](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210305181242965.png)



## 6 类和对象



### 6.1 方法区



#### 6.1.1 类信息

```go
type Class struct {
	accessFlags       uint16	//访问标志
	name              string // thisClassName
	superClassName    string	//父类名
	interfaceNames    []string	//接口名
	constantPool      *classfile.ConstantPool	//常量池引用
	fields            []*Field	//字段引用
	methods           []*Method	//方法引用
	loader            *ClassLoader	//类加载器引用
	superClass        *Class	//父类引用
	interfaces        []*Class	//接口引用
	instanceSlotCount uint		//实例槽总数
	staticSlotCount   uint		//静态槽总数
	staticVars        Slots		//静态变量
}
```

accessFlags是类的访问标志，字段和方法也有访问标志，但具体标志位的含义有所不同，因此根据JVM规范，将各个比特为的含义统一定义在access.flags.go文件中。

![image-20210305174714029](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210305174714029.png)

```go
const (
	ACC_PUBLIC       = 0x0001 // class field method
	ACC_PRIVATE      = 0x0002 //       field method
	ACC_PROTECTED    = 0x0004 //       field method
	ACC_STATIC       = 0x0008 //       field method
	ACC_FINAL        = 0x0010 // class field method
	ACC_SUPER        = 0x0020 // class
	ACC_SYNCHRONIZED = 0x0020 //             method
	ACC_VOLATILE     = 0x0040 //       field
	ACC_BRIDGE       = 0x0040 //             method
	ACC_TRANSIENT    = 0x0080 //       field
	ACC_VARARGS      = 0x0080 //             method
	ACC_NATIVE       = 0x0100 //             method
	ACC_INTERFACE    = 0x0200 // class
	ACC_ABSTRACT     = 0x0400 // class       method
	ACC_STRICT       = 0x0800 //             method
	ACC_SYNTHETIC    = 0x1000 // class field method
	ACC_ANNOTATION   = 0x2000 // class
	ACC_ENUM         = 0x4000 // class field
)
```

name, superClassName和interfaceNames字段分别存放类名、超类名以及接口名，这些类名都是完全限定名，具有java/lang/Object前缀。constantPool字段存放运行时常量池指针,fields和methods字段分别存放字段表和方法表。



#### 6.1.2 字段信息

字段与方法都属于类的成员，它们有一些共同的信息（访问标志、名字、描述符）。我们将这些信息抽象出来，放到一个结构体中：

```go
type ClassMember struct {
	accessFlags uint16	//访问标志
	name        string	//简单名称
	descriptor  string	//描述符
	class       *Class	//所属类引用
}
```

接下来定义字段结构：

```go
type Field struct {
	ClassMember
	constValueIndex uint
	slotId          uint
}
```





#### 6.1.3 方法信息

方法比字段稍微复杂一点，因为方法中有字节码：

```go
type Method struct {
	ClassMember
	maxStack  uint16
	maxLocals uint16
	code      []byte
}
```



到此为止，出了ConstantPool还没有介绍以外，已经定义了4个结构体，这些结构体之间的关系如下：

![image-20210305191409647](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210305191409647.png)





### 6.2 运行时常量池

运行时常量池主要存放两类信息：字面量（literal）和符号引用（symbolic reference）。字面量包括整数、浮点数和字符串字面量；符号引用包括类符号引用、字段符号引用、类方法符号引用和接口方法符号引用。

```go
type Constant interface{}

type ConstantPool struct {
	class  *Class
	consts []Constant//存放常量
}
```

```go
//核心逻辑就是把[]classfile.ConstantInfo转换成[]heap.Constant。
func newConstantPool(class *Class, cfCp classfile.ConstantPool) *ConstantPool {
	cpCount := len(cfCp)
	consts := make([]Constant, cpCount)
	rtCp := &ConstantPool{class, consts}

	for i := 1; i < cpCount; i++ {
		cpInfo := cfCp[i]
		switch cpInfo.(type) {
		case 
      ....
		default:
			// todo
		}
	}
	return rtCp
}
```



#### 6.2.1 类符号引用

类符号引用、字段符号引用、类方法符号引用和接口方法符号引用有一些共性，所以将其抽取出建立结构体：

```go
type SymRef struct {
	cp        *ConstantPool
	className string
	class     *Class
}
```

cp字段存放符号引用所在的运行时常量池指针，这样就可以通过符号引用访问到运行时常量池，进一步访问到类数据。className存放类的全限定名。class字段缓存解析后的类结构体指针。

对于类符号引用，只要有类名就可以解析符号引用。

对于字段，首先要解析类符号引用得到类数据，然后用字段名和描述符查找字段数据。

方法符号引用的解析过程和字段符号引用类似。



定义ClassRef结构体：

```go
type ClassRef struct {
	SymRef
}
```

根据class文件中存储的类常量创建ClassRef实例：

```go
func newClassRef(cp *ConstantPool, classInfo *classfile.ConstantClassInfo) *ClassRef {
	ref := &ClassRef{}
	ref.cp = cp
	ref.className = classInfo.Name()
	return ref
}
```



#### 6.2.2 字段符号引用

定义MemberRef结构体来存放字段和方法符号引用共有的信息：

```go
type MemberRef struct {
	SymRef
	name       string
	descriptor string
}
```

在JVM规范中，并没有规定一个类的字段名不能相同，其只规定了两个字段的描述符+字段名必须不同。这也就是为什么字段符号引用还需要描述符来表示的原因。

创建字段符号引用：

```go
type FieldRef struct {
	MemberRef
	field *Field
}
```

field字段缓存解析后的字段指针。



#### 6.2.3 方法符号引用

```go
type MethodRef struct {
	MemberRef
	method *Method
}
```



#### 6.2.4 接口方法符号引用

```go
type InterfaceMethodRef struct {
	MemberRef
	method *Method
}
```



至此所有的符号引用都已经定义完毕，他们的继承结构图如下：

![image-20210305195051976](/Users/huangyucai/Library/Application Support/typora-user-images/image-20210305195051976.png)



### 6.3 类加载器

Java虚拟机的类加载系统十分复杂，本节将初步实现一个简化版的类加载器:

```go
/*
class names:
    - primitive types: boolean, byte, int ...
    - primitive arrays: [Z, [B, [I ...
    - non-array classes: java/lang/Object ...
    - array classes: [Ljava/lang/Object; ...
*/
type ClassLoader struct {
	cp       *classpath.Classpath
	classMap map[string]*Class // loaded classes
}
```

ClassLoader依赖Classpath来搜索和读取class文件，cp字段保存Classpath指针。

classMap字段记录已经加载的类数据，key是类的完全限定名。我们可以把classMap当作方法区的具体实现，其中存储着类数据。

LoadClass方法把类数据加载到方法区：

```go
func (self *ClassLoader) LoadClass(name string) *Class {
	if class, ok := self.classMap[name]; ok {
		// already loaded
		return class
	}
	return self.loadNonArrayClass(name)
}
```

先查找classMap，看类是否已经被加载，如果是，直接返回类数据，否则调用loadNonArrayClass方法加载类。数组类和普通类又很大不同，它的数据并不是来自class文件，而是由Java虚拟机在运行期间生成的。目前我们展示不考虑数组类。loadNonArrayClass方法如下：

```go
func (self *ClassLoader) loadNonArrayClass(name string) *Class {
	data, entry := self.readClass(name)
	class := self.defineClass(data)
	link(class)
	fmt.Printf("[Loaded %s from %s]\n", name, entry)
	return class
}
```

类的加载大致可以分为三个步骤：首先找到class文件并把数据读取到内存；然后解析class文件，生成虚拟机可以使用的类数据，并放入方法区；最后进行链接。



#### 6.3.1 readClass方法

```go
func (self *ClassLoader) readClass(name string) ([]byte, classpath.Entry) {
	data, entry, err := self.cp.ReadClass(name)
	if err != nil {
		panic("java.lang.ClassNotFoundException: " + name)
	}
	return data, entry
}
```

实际上调用Classpath的ReadClass方法。



#### 6.3.2 defineClass方法

```go
func (self *ClassLoader) defineClass(data []byte) *Class {
	class := parseClass(data)
	class.loader = self
	resolveSuperClass(class)
	resolveInterfaces(class)
	self.classMap[class.name] = class
	return class
}
```

首先调用parseClass函数把class文件数据转换成Class结构体。

```go
func parseClass(data []byte) *Class {
	cf, err := classfile.Parse(data)//将数据读取为classfile
	if err != nil {
		//panic("java.lang.ClassFormatError")
		panic(err)
	}
	return newClass(cf)//将classfile解析成Class结构
}
```

Class结构体的superClass和interfaces字段存放超类名以及直接接口表，这些类名都是符号引用。根据JVM规范5.3.5，调用resolveSuperClass()和resolveInterfaces()函数解析这些类符号引用。

```go
// jvms 5.4.3.1
func resolveSuperClass(class *Class) {
	if class.name != "java/lang/Object" {
		class.superClass = class.loader.LoadClass(class.superClassName)
	}
}
func resolveInterfaces(class *Class) {
	interfaceCount := len(class.interfaceNames)
	if interfaceCount > 0 {
		class.interfaces = make([]*Class, interfaceCount)
		for i, interfaceName := range class.interfaceNames {
			class.interfaces[i] = class.loader.LoadClass(interfaceName)
		}
	}
}
```

递归调用LoadClass方法来加载它的超类和直接接口。

#### 6.3.3 link()

解析分为以下三个阶段：

* **验证（Verify）**：目的在于确保class文件的字节流中包含的信息符合当前虚拟机的要求，保证被加载类的正确性，不会危害虚拟机的自身安全。主要包括四种验证：**文件格式验证**、**元数据验证**、**字节码验证**、**符号引用验证**。

  * **文件格式验证：**

    * 是否以魔数`0xCAFEBABE`开头。
    * 主次版本号是否在当前Java虚拟机接受范围之内。
    * 常量池的常量中是否有不被支持的常量类型（检查常量tag标志）。
    * 指向常量的各种索引值中是否有**指向不存在的常量**或**不符合类型的常量**。
    * CONSTANT_Utf8_info型的常量中是否有不符合UTF-8编码的数据。
    * Class文件中各个部分及文件本身是否有被删除的或附加的其他信息。

    此验证阶段是**基于二进制字节流进行的**，只有通过了这个阶段的验证之后，这段字节流才被允许进入Java虚拟机内存的方法区中进行存储，所以后面的三个验证阶段全部是**基于方法区的存储结构上进行的**，不会再直接读取、操作字节流了。

  * **元数据验证：**对字节码描述的信息进行语义分析，以保证其描述的信息符合《Java语言规范》的要求：

    * 这个类是否有父类（除了java.lang.Object之外，所有的类都应当有父类）
    * 这个类的父类是否继承了不允许被继承的类（被final修饰的类）。
    * 如果这个类不是抽象类，是否实现了其父类或接口之间要求实现的所有方法。
    * 类中的字段、方法是否与父类产生矛盾（例如覆盖了父类的final字段，或者出现了不符合规则的方法重载，例如方法参数都一致，但返回值类型却不同等）

  * **字节码验证：**通过数据流分析和控制流分析，确定程序语义是合法的、符合逻辑的。对类的方法体（Class文件的Code属性）进行校验分析，保证被校验类的方法在运行时不会做出危害虚拟机安全的行为。

    * 保证任意时刻操作数栈的数据类型与指令代码序列都能配合工作，例如不会出现类似于“在操作数栈放置了一个int类型的数据，使用时却按long类型来加载入本地变量表”这样的情况。
    * 保证任何跳转指令都不会跳转到方法体以外的字节码指令上。
    * 保证方法体中的类型转换总是有效的。例如可以把一个子类对象赋值给父类数据类型。

    即使字节码验证阶段进行了再大量、再严密的检查，也依然不能保证方法体一定是安全的。涉及到离散数学中的停机问题，即不能通过程序准确地检查出程序是否能在有限的时间之内结束运行。

    由于数据流分析和控制流分析的高度复杂性，Java虚拟机的设计团队为了避免过多的执行时间消 耗在字节码验证阶段中，在JDK 6之后的Javac编译器和Java虚拟机里进行了一项联合优化，把尽可能 多的校验辅助措施挪到Javac编译器里进行。

  * **符号引用验证：**此校验行为发生在虚拟机将符号引用转化为直接引用的时候。符号引用验证可以看作是对类自身以外（常量池中的各种符号引用）的各类信息进行匹配性校验，通俗来说就是，该类是否缺少或者被禁止访问它依赖的某些外部类、方法、字段等资源。

    * 符号引用中通过字符串描述的全限定类名是否能找到对应的类。
    * 在指定类中是否存在符合方法的字段描述符及简单名称说描述的方法和字段。
    * 符号引用中的类、字段、方法的可访问性（private、protected、public、<package>）是否可被当前类访问。

* **准备（Prepare）**：为**类变量**分配内存并且设置该类变量的**默认初始值，即零值**。这里**不包含用final修饰的static变量**，因为final在编译的时候就会分配了，准备阶段会显式初始化。这里**不会为实例变量分配初始化**，类变量会分配在方法区中，而**实例变量是会随着对象一起分配到Java堆**当中，所以在这个阶段对象还没创建，故实例变量不会被初始化。

* **解析（Resolve）**：将常量池的**符号引用转换为直接引用**的过程。事实上，解析操作往往会在JVM执行完初始化之后再执行。符号引用就是一组符号来描述所引用的目标。符号引用的字面量形式明确定义在《Java虚拟机规范》的Class文件格式中。直接引用就是直接指向目标的指针、相对偏移量或间接定位到目标的句柄。下图就是符号引用：

  <img src="/Users/huangyucai/Library/Application Support/typora-user-images/image-20201212155028986.png" alt="image-20201212155028986" style="zoom:50%;" />

  解析动作主要针对类或接口、字段、类方法、接口方法、方法类型等。对应常量池中的`CONSTANT_Class_info`、`CONSTANT_Fieldref_info`、`CONSTANT_Methodref_info`等。

  《Java虚拟机规范》之中并未规定解析阶段发生的具体时间，只要求了在执行ane-warray、 checkcast、getfield、getstatic、instanceof、invokedynamic、invokeinterface、invoke-special、 invokestatic、invokevirtual、ldc、ldc_w、ldc2_w、multianewarray、new、putfield和putstatic这17个用于 操作符号引用的字节码指令之前，先对它们所使用的符号引用进行解析。所以虚拟机实现可以根据需要来自行判断，到底是在类被加载器加载时就对常量池中的符号引用进行解析，还是等到一个符号引用将要被使用前才去解析它。

  类似地，对方法或者字段的访问，也会在解析阶段中对它们的可访问性（public、protected、 private、<package>）进行检查

  * **类或接口的解析**
  * **字段解析：**如果如果有一个同名字段同时出现在某个类的接口和父类中，Javac编译器将提示“The fild xxx is ambiguous”，并且拒绝编译这段代码。首先搜索当前类、再搜索当前类实现的接口及其父接口，最后搜索父类
  * **方法解析**：首先搜索当前类，再搜索父类，最后搜索接口和父接口。
  * **接口方法解析**



在这本书中，作者没有选择实现验证过程。本来我想实现一下的，但是看到JVM规范中验证算法的篇幅：

![image-20210305203445251](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210305203445251.png)

还是算了吧...

不过验证的主要内容已经在上文说清楚了。



### 6.4 对象、实例变量和类变量

和局部变量类似，实例变量和类变量也可以放在Slot中。所以给实例变量和类变量同样定义一个变量表：

```go
type Slot struct {
	num int32
	ref *Object
}

type Slots []Slot
```



接下来给Object结构体添加两个字段：

```go
type Object struct {
	class  *Class
	fields Slots
}
```

接下来的问题是，如何知道静态变量和实例变量需要多少空间，以及哪个字段对应Slots中的哪个位置呢？

第一个问题：数数。递归地数父类静态变量和实例变量的数字。

第二个问题：在数字段的时候，给字段按顺序编号即可。有几个注意点：

* 首先，静态字段和实例字段要分开编号。
* 对于实例字段，一定要从继承关系的最顶端，也就是java.lang.Object开始编号。
* 最后，编号时，要考虑long和double是占两个slot的。

定义prepare函数，完成以上需求：

```go
// jvms 5.4.2
func prepare(class *Class) {
	calcInstanceFieldSlotIds(class)
	calcStaticFieldSlotIds(class)
	allocAndInitStaticVars(class)
}
```

calcInstanceFieldSlotIds函数计算实例字段的个数，同时给它们编号，代码如下：

```go
func calcInstanceFieldSlotIds(class *Class) {
	slotId := uint(0)
	if class.superClass != nil {
		slotId = class.superClass.instanceSlotCount
	}
	for _, field := range class.fields {
		if !field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.instanceSlotCount = slotId
}
```

calcStaticFieldSlotIds函数计算静态字段的个数，同时给它们编号，代码如下：

```go
func calcStaticFieldSlotIds(class *Class) {
	slotId := uint(0)
	for _, field := range class.fields {
		if field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.staticSlotCount = slotId
}
```

allocAndInitStaticVars函数给**类变量(注意：这里没有 实例变量初始化)**分配空间，然后给它们赋予初始值，代码如下：

```go
func allocAndInitStaticVars(class *Class) {
	class.staticVars = newSlots(class.staticSlotCount)
	for _, field := range class.fields {
		if field.IsStatic() && field.IsFinal() {//
			initStaticFinalVar(class, field)
		}
	}
}
```

如果静态变量属于基本类型或者String类型，且用final修饰，那么它的值在编译期就已知，存储在class文件的常量池中，initStaticFinalVar函数从常量池中加载常量值，然后给静态变量赋值：

```go
func initStaticFinalVar(class *Class, field *Field) {
	vars := class.staticVars
	cp := class.constantPool
	cpIndex := field.ConstValueIndex()//获取static final的常量值索引
	slotId := field.SlotId()

	if cpIndex > 0 {
		switch field.Descriptor() {
		case "Z", "B", "C", "S", "I":
			val := cp.GetConstant(cpIndex).(int32)
			vars.SetInt(slotId, val)
		case "J":
			val := cp.GetConstant(cpIndex).(int64)
			vars.SetLong(slotId, val)
		case "F":
			val := cp.GetConstant(cpIndex).(float32)
			vars.SetFloat(slotId, val)
		case "D":
			val := cp.GetConstant(cpIndex).(float64)
			vars.SetDouble(slotId, val)
		case "Ljava/lang/String;":
			panic("todo")
		}
	}
}
```

Go语言会保证新创建的Slot结构体有默认的值，所以我们无需任何操作就能保证静态变量有默认的初始值0和ni l。





### 6.5 类和字段符号引用解析



#### 6.5.1 类符号引用解析

在cp_symref.go文件中定义ResolvedClass方法：

```go
func (self *SymRef) ResolvedClass() *Class {
	if self.class == nil {
		self.resolveClassRef()
	}
	return self.class
}
```

如果类符号引用已经解析，ResolvedClass方法直接返回类指针，否则调用resolveClassRef方法进行解析：

```go
// jvms8 5.4.3.1
func (self *SymRef) resolveClassRef() {
	d := self.cp.class
	c := d.loader.LoadClass(self.className)
	if !c.isAccessibleTo(d) {
		panic("java.lang.IllegalAccessError")
	}

	self.class = c
}
```

解析类的步骤：

To resolve an unresolved symbolic reference from D to a class or interface C denoted by `N`, the following steps are performed:

1. The defining class loader of D is used to create a class or interface denoted by `N`. This class or interface is C. The details of the process are given in [§5.3](https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-5.html#jvms-5.3).

   Any exception that can be thrown as a result of failure of class or interface creation can thus be thrown as a result of failure of class and interface resolution.

2. If C is an array class and its element type is a `reference` type, then a symbolic reference to the class or interface representing the element type is resolved by invoking the algorithm in [§5.4.3.1](https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-5.html#jvms-5.4.3.1) recursively.

3. Finally, access permissions to C are checked.

   - If C is not accessible ([§5.4.4](https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-5.html#jvms-5.4.4)) to D, class or interface resolution throws an `IllegalAccessError`.

     This condition can occur, for example, if C is a class that was originally declared to be `public` but was changed to be non-`public` after D was compiled.

If steps 1 and 2 succeed but step 3 fails, C is still valid and usable. Nevertheless, resolution fails, and D is prohibited from accessing C.

![image-20210305211104696](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210305211104696.png)





#### 6.5.2 字段符号引用解析

Therefore, any exception that can be thrown as a result of failure of resolution of a class or interface reference can be thrown as a result of failure of field resolution. If the reference to C can be successfully resolved, an exception relating to the failure of resolution of the field reference itself can be thrown.

When resolving a field reference, field resolution first attempts to look up the referenced field in C and its superclasses:

1. If C declares a field with the name and descriptor specified by the field reference, field lookup succeeds. The declared field is the result of the field lookup.
2. Otherwise, field lookup is applied recursively to the direct superinterfaces of the specified class or interface C.
3. Otherwise, if C has a superclass S, field lookup is applied recursively to S.
4. Otherwise, field lookup fails.

Then:

- If field lookup fails, field resolution throws a `NoSuchFieldError`.

- Otherwise, if field lookup succeeds but the referenced field is not accessible ([§5.4.4](https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-5.html#jvms-5.4.4)) to D, field resolution throws an `IllegalAccessError`.

- Otherwise, let `<`E, `L1``>` be the class or interface in which the referenced field is actually declared and let `L2` be the defining loader of D.

  Given that the type of the referenced field is Tf, let T be Tf if Tf is not an array type, and let T be the element type ([§2.4](https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-2.html#jvms-2.4)) of Tf otherwise.



![image-20210305211512717](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210305211512717.png)



### 6.6 类和对象相关指令

* new创建类实例。
* putstatic和getstatic用于存取静态变量。
* putfield和getfield用于存取实例变量。
* instanceof和checkcast用于判断对象是否属于某种类型。
* ldc系列指令把运行时常量池中的常量推到操作数栈顶。



#### 6.6.1 new指令

new是专门用来创建类实例的，数组由专门的指令创建。

```go
type NEW struct{ base.Index16Instruction }
```

new指令的操作数是一个uint16索引，来自字节码。通过这个索引，可以从当前类的运行时常量池中找到一个类符号引用。解析这个符号引用，拿到类数据，然后创建对象，并把对象引用推入栈顶，new指令的工作就完成了。

按照JVM规定， 接口 和抽象类是不能被实例化的，需要抛出InstantiationError异常。另外，如果解析后的类还没有初始化，则需要先初始化类。初始化放在后面实现了方法调用后再实现。



#### 6.6.2 putstatic和getstatic指令

用于存取类的静态变量。

putstatic指令给类的某个静态变量赋值，它需要两个操作数：

* 第一个操作数是uint16索引，来自字节码。通过这个索引可以从当前类的运行时常量池中找到一个字段符号引用，解析这个符号引用就可以知道给类的哪个静态变量赋值。
* 第二个操作数是要赋给静态变量的值，从操作数栈中弹出。

getstatic指令取出某个静态变量值，然后推入栈顶。getstatic指令只需要一个操作数。

> 目前，指令的操作数来源有三项：
>
> 1. 来自指令本身
> 2. 来自字节码
> 3. 来自操作数栈



#### 6.6.3 putfiled和putfield指令

putfield指令给实例变量赋值，需要三个操作数：

* 前两个操作数是常量池索引和变量值，用法和putstatic一样。
* 第三个操作数是对象引用，从操作数中弹出。



getfield指令获取实例变量值，然后推入操作数栈，需要两个操作数：

* 第一个是uint16常量池索引。
* 第二个操作数是对象引用。



#### 6.6.4 instanceof和checkcast指令

instance指令判断对象是否是某个类的实例（活着对象的类是否实现了某个接口），并把结果推入操作数栈。需要两个操作数：

* 第一个是uint16索引，通过这个索引可以从当前类的运行时常量池找到一个类符号引用。
* 第二个操作数是对象引用，从操作数栈中弹出。



checkcast指令和instanceof指令很像，区别在于：instanceof指令会改变操作数栈（弹出对象引用，推入判断结果）；checkcast则不能改变操作数栈（如果判断失败，直接抛出ClassCastException异常）。



#### 6.6.5 ldc指令

ldc系列指令从运行时常量池加载常量值，并把它推入操作数栈。ldc系列指令属于常量类指令，共3条：

* ldc和ldc_w指令用于加载int、float和字符串常量，java.lang.Class实例或者MethodType和MethodHandle实例。ldc和ldc_w的区别仅在于操作数的宽度。
* ldc2_w指令用于加载long和double常量。



### 6.7 测试

修改main方法：

```go
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
```

先创建ClassLoader实例，然后用它来加载主类，最后执行主类的main方法。

Class结构体的GetMainMethod：

```go
func (self *Class) GetMainMethod() *Method {
   return self.getStaticMethod("main", "([Ljava/lang/String;)V")
}

func (self *Class) getStaticMethod(name, descriptor string) *Method {
   for _, method := range self.methods {
      if method.IsStatic() &&
         method.name == name &&
         method.descriptor == descriptor {

         return method
      }
   }
   return nil
}
```

修改解释器：

```go
func interpret(method *heap.Method){
	thread := rtda.NewThread()//创建线程
	frame := thread.NewFrame(method)//创建栈帧
	thread.PushFrame(frame)//栈帧入Java虚拟机栈
	defer catchErr(frame)
	loop(thread, method.Code())//执行方法
}
```



对象需要初始化，所以每个类都至少有一个构造函数。即使用户自己不定义，编译器也会自动生成一个默认构造器。在创建类实例的时候，编译器会在new指令的后面加入invokespecial指令来调用构造器初始化对象。

测试结果：

![image-20210308192445710](/Users/huangyucai/Library/Application Support/typora-user-images/image-20210308192445710.png)



### 小结

本章实现了方法区、运行时常量池、类和对象结构体、一个简单的类加载器，以及ldc和部分引用指令。





## 7 方法调用和返回



### 7.1 方法调用概述

从调用的角度来看，方法可以分为两类：静态方法（类方法）和实例方法。静态方法通过类来调用，实例方法通过对象引用来调用。静态方法是静态绑定的，也就是说，最终调用的是哪个方法在编译期就已经确定。实例方法则支持动态绑定，最终调用哪个方法可能要推迟到运行期才能知道。

从实现的角度来看，方法可以分为三类：没有实现（抽象方法）、用Java语言（或者JVM上的其他语言，如Groovy和Scala）实现以及用本地语言（C和C++）实现。静态方法和抽象方法是互斥的。在Java8 之前，接口只能包含抽象方法。为了实现Lambda表达式，Java 8放宽了这一限制，在接口中也可以定义静态方法和默认方法（本实现不考虑实现这两种方法）。

在Java7之前，JVM规范一共一共了4条方法调用指令：

1. invokestatic指令用来调用静态方法。
2. invokespecial指令用来调用无需动态绑定的实例方法，包括构造函数、私有方法和通过super关键字调用的超类方法。
3. 如果针对接口的引用调用方法，就使用invokeinterface指令，否则使用invokevirtual指令。

为了更好地支持动态类型语言，Java7增加了一条方法调用指令invokedynamic指令。

**JVM是如何调用方法的？**

首先，方法调用指令需要n+1个操作数，其中第1个操作数是uint16索引，在字节码中紧跟在指令操作码后，通过这个索引，可以从当前类的运行时常量池中找到一个方法符号引用，解析这个符号引用就可以得到一个方法。这个方法并不一定就是最终调用的方法，可能还需要一个查找过程（invokevirtual指令）才能找到最终要调用的方法。剩下的n个操作数是要被传递给调用方法的参数，从操作数栈中弹出。

如果要执行的是Java方法（非native），下一步是给这个方法创建一个新的帧，并把它推到Java虚拟机栈顶。传递参数后，新的方法就可以开始执行了。

![image-20210308195750042](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210308195750042.png)

方法的最后一条指令是某个返回指令，这个指令负责把方法的返回值推入前一帧的操作数栈顶，然后把当前帧从JVM虚拟机栈中弹出。



### 7.2 解析方法符号引用

#### 7.2.1 非接口方法符号引用

如果类D想通过方法符号引用访问类C的某个方法，先要解析符号引用得到类C。如果C是接口，则抛出IncompatibleClassChangeError异常，否则根据方法和描述符查找方法。如果找不到对应的方法，则抛出NoSuchMethodError异常，否则检查类D是否有权限访问该方法，如果没有，则抛出IllegalAccessError异常。



#### 7.2.2 接口方法符号引用

如果能在接口中找到方法，就返回找到的方法，否则调用loopupMethodInInterfaces函数在超接口中寻找。



### 7.3 方法调用和参数传递

首先，要确定方法的参数在局部变量表中占用多少位置：

* double和long占用2个位置
* 对于实例方法Java编译器会在参数列表的前面添加一个参数，这个隐藏的参数就是this引用。

假设司机的参数占据n个位置，依次把这n个变量从调用者的操作数栈中弹出，放入被调用方法的局部变量表中，参数传递完成。

在解析方法描述符时，判断传入的参数类型是否为double和long，如果遇到了这两种类型，变量表槽数就多+1。

```go
func (self *Method) calcArgSlotCount() {
	parsedDescriptor := parseMethodDescriptor(self.descriptor)
	for _, paramType := range parsedDescriptor.parameterTypes {
		self.argSlotCount++
		if paramType == "J" || paramType == "D" {
			self.argSlotCount++
		}
	}
	if !self.IsStatic() {
		self.argSlotCount++ // `this` reference
	}
}
```



### 7.4 返回指令

方法执行完毕之后，需要把结果返回给调用方。返回指令属于控制指令，一共有6条。6条指令都不需要操作数。

* return：无返回值。把当前帧从Java虚拟机栈中弹出即可。
* areturn：返回引用。将当前帧弹出，然后把此帧的操作数弹出，然后放入当前帧（与前不同）的操作数栈顶。
* ireturn：返回int。
* lreturn：返回long。
* freturn：返回float。
* dreturn：返回double。



### 7.5 方法调用指令



#### 7.5.1 invokestatic指令

```go
// Invoke a class (static) method
type INVOKE_STATIC struct{ base.Index16Instruction }
```

假定解析符号引用后得到方法M。M必须是静态方法，否则抛出IncompatibleClassChangeError异常。如果声明M方法的类还没有初始化，则先要初始化该类。

对于invokestatic指令，M就是最终要执行的方法，调用InvokeMethod直接执行就可。



#### 7.5.2 invokespecial指令

先拿到当前类、当前常量池、方法符号引用，然后解析符号引用，拿到解析后的类和方法。

假定从方法符号引用中解析出来的类是C，方法是M。如果M是构造函数，则声明M的类必须是C，否则抛出NoSuchMethodError异常。如果M是静态方法，则抛出IncompatibleClassChangeError异常。



#### 7.5.3 invokecirtual指令和Invokeinterface指令

An *invokevirtual* instruction is type safe iff all of the following are true:

- Its first operand, `CP`, refers to a constant pool entry denoting a method named `MethodName` with descriptor `Descriptor` that is a member of a class `MethodClassName`.
- `MethodName` is not `<init>`.
- `MethodName` is not `<clinit>`.
- One can validly replace types matching the class `MethodClassName` and the argument types given in `Descriptor` on the incoming operand stack with the return type given in `Descriptor`, yielding the outgoing type state.
- If the method is `protected`, the usage conforms to the special rules governing access to `protected` members ([§4.10.1.8](https://docs.oracle.com/javase/specs/jvms/se8/html/jvms-4.html#jvms-4.10.1.8)).

An *invokeinterface* instruction is type safe iff all of the following are true:

- Its first operand, `CP`, refers to a constant pool entry denoting an interface method named `MethodName` with descriptor `Descriptor` that is a member of an interface `MethodIntfName`.
- `MethodName` is not `<init>`.
- `MethodName` is not `<clinit>`.
- Its second operand, `Count`, is a valid count operand (see below).
- One can validly replace types matching the type `MethodIntfName` and the argument types given in `Descriptor` on the incoming operand stack with the return type given in `Descriptor`, yielding the outgoing type state.



#### 小结

>  为什么要单独定义invokeinterface指令？

可以统一用invokevirtual指令，但是会降低效率。

当JVM通过invokevirtual调用方法时，this引用指向某个类（或其子类）的实例。因为类的继承关系是固定的，所以虚拟机采用了一种叫做vtable(virtual method table)的技术加速方法查找。但是当通过invokeinterface指令调用接口方法时，this引用可以指向任何实现了该接口的类的实例，所以无法使用vtable技术。



### 7.6 改进解释器

改进解释器，令其支持方法调用。

```go
func loop(thread *rtda.Thread, logInst bool) {
	reader := &base.BytecodeReader{}
	for {
		frame := thread.CurrentFrame()
		pc := frame.NextPC()
		thread.SetPC(pc)

		// decode
		reader.Reset(frame.Method().Code(), pc)
		opcode := reader.ReadUint8()
		inst := instructions.NewInstruction(opcode)
		inst.FetchOperands(reader)
		frame.SetNextPC(reader.PC())

		if logInst {
			logInstruction(frame, inst)
		}

		// execute
		inst.Execute(frame)
		if thread.IsStackEmpty() {
			break
		}
	}
}
```

在每次循环开始前，先拿到当前帧，然后根据PC从当前方法中解码出一条指令。指令执行完毕后，判断Java虚拟机栈中是否还有帧，如果没有则退出循环，否则继续。

一个方法执行到return指令时，会将此方法的栈帧从Java虚拟机栈中弹出。



### 7.7 测试方法调用

为命令行工具增加两个选项：

* -verbose: class选项，可以控制是否把类加载信息输出到控制台。
* -verbose: inst选项，用来控制是否把指令执行信息输出到控制台。

```go
type Cmd struct {
	helpFlag	bool//用于指定是否显示帮助信息
	versionFlag	bool//用于指定是否显示版本信息
	verboseClassFlag	bool
	verboseInstFlag	bool
	clspath		string//用于指定类路径
	Xjre		string//用于指定jre路径
	class		string//用于指定主类
	maxStackSize	uint//用于指定虚拟机栈的最大深度
	args		[]string//用于指定其他参数
}
```

写一个斐波那契数列计算

```java
public class Fiboo{

    private static long dfs(long n){
    	if(n <= 1)	{return n;}
    	return dfs(n-1) + dfs(n-2);
    }

    public static void main(String[] args){
        long m = dfs(30);
        System.out.println(m);
    }
}
```

![image-20210309200853071](https://hyc-pic.oss-cn-hangzhou.aliyuncs.com/image-20210309200853071.png)

测试成功。

### 7.8 类初始化

类初始化就是执行类的初始化方法<clinit>。类的初始化在以下情况出发：

* 执行new指令创建类实例，但类还没被初始化。
* 执行putstatic, getstatic指令存取类的静态变量，但声明该字段的类还没有被初始化。
* 执行invokestatic调用类的静态方法，但声明该方法的类还没被初始化。
* 当初始化一个类时，如果类的超类还没有被初始化，要先初始化类的超类。
* 执行某些反射操作时。

为了判断类是否已经初始化，需要给Class结构体添加一个字段：

```go
initStarted bool
```

类的初始化其实分为几个阶段，但是由于我们的类加载器还不够完善，所以先使用一个布尔状态表示就足够了。

修改new指令：

```go
func (self *NEW) Execute(frame *rtda.Frame) {
	cp := frame.Method().Class().ConstantPool()
	classRef := cp.GetConstant(self.Index).(*heap.ClassRef)
	class := classRef.ResolvedClass()
	// todo: init class

	if !class.InitStarted() {//如果还没初始化
		frame.RevertNextPC()//让PC重新指向当前指令
		base.InitClass(frame.Thread(), class)//初始化类
		return
	}

	if class.IsInterface() || class.IsAbstract() {
		panic("java.lang.InstantiationError")
	}

	ref := class.NewObject()
	frame.OperandStack().PushRef(ref)
}
```

其他指令类似。



