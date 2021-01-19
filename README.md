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

