package classfile

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

type ConstantInfo interface {
	//读取常量信息，由具体的常量结构体实现
	readInfo(reader *ClassReader)
}
//先读出tag值，再调用newConstantInfo创建具体的常量，最后调用常量的readInfo()方法读取常量信息
func readConstantInfo(reader *ClassReader, cp ConstantPool) ConstantInfo {
	tag := reader.readUint8()
	c := newConstantInfo(tag, cp)
	c.readInfo(reader)
	return c
}

func newConstantInfo(tag uint8, cp ConstantPool) ConstantInfo {
	switch tag {
	case CONSTANT_Integer:
		return &ConstantIntegerInfo{}
	case CONSTANT_Float:
		return &ConstantFloatInfo{}
	case CONSTANT_Long:
		return &ConstantLongInfo{}
	case CONSTANT_Double:
		return &ConstantDoubleInfo{}
	case CONSTANT_Utf8:
		return &ConstantUtf8Info{}
	case CONSTANT_String:
		return &ConstantStringInfo{cp: cp}
	case CONSTANT_Class:
		return &ConstantClassInfo{cp: cp}
	case CONSTANT_Fieldref:
		return &ConstantFieldrefInfo{ConstantMemberrefInfo{cp: cp}}
	case CONSTANT_Methodref:
		return &ConstantMethodrefInfo{ConstantMemberrefInfo{cp: cp}}
	case CONSTANT_InterfaceMethodref:
		return &ConstantInterfaceMethodrefInfo{ConstantMemberrefInfo{cp: cp}}
	case CONSTANT_NameAndType:
		return &ConstantNameAndTypeInfo{}
	case CONSTANT_MethodType:
		return &ConstantMethodTypeInfo{}
	case CONSTANT_MethodHandle:
		return &ConstantMethodHandleInfo{}
	case CONSTANT_InvokeDynamic:
		return &ConstantInvokeDynamicInfo{}
	default:
		panic("java.lang.ClassFormatError: constant pool tag!")
	}


}