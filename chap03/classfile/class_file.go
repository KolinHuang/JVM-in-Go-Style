package classfile

import "fmt"

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
	cf.read(cr)
	return
}

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

func (self *ClassFile) readAndCheckMagic(reader *ClassReader) {
	magic := reader.readUint32()
	if magic != 0xCAFEBABE {
		panic("java.lang.ClassFormatError: miss cafebabe!")
	}
}
func (self *ClassFile) readAndCheckVersion(reader *ClassReader) {
	self.minorVersion = reader.readUint16()
	self.majorVersion = reader.readUint16()
	switch self.majorVersion {
	case 45:
		return
	case 46, 47, 48, 49, 50, 51, 52:
		if self.minorVersion == 0 {
			return
		}
	}

	panic("java.lang.UnsupportedClassVersionError! ")
}



func (self *ClassFile) MinorVersion() uint16 {
	return self.minorVersion
}

func (self *ClassFile) MajorVersion() uint16 {
	return self.majorVersion
}

func (self *ClassFile) ConstantPool() ConstantPool {
	return self.constantPool
}

func (self *ClassFile) AccessFlags() uint16 {
	return self.accessFlags
}

func (self *ClassFile) Fileds() []*MemberInfo {
	return self.fields
}

func(self *ClassFile) Methods() []*MemberInfo {
	return self.methods
}

//从常量池中查找类名
func (self *ClassFile) ClassName() string {
	return self.constantPool.getClassName(self.thisClass)
}
//从常量池中查找父类名
func (self *ClassFile) SuperClassName() string {
	if self.superClass > 0 {
		return self.constantPool.getClassName(self.superClass)
	}
	return ""
}
//从常量池中查找接口名
func (self *ClassFile) InterfaceNames() []string {
	interfaceNames := make([]string, len(self.interfaces))
	for i, cpIdx := range self.interfaces {
		interfaceNames[i] = self.constantPool.getClassName(cpIdx)
	}
	return interfaceNames
}

