package classfile

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

func (self *CodeAttribute) readInfo(reader *ClassReader) {
	self.maxStack = reader.readUint16()
	self.maxLocal = reader.readUint16()
	codeLength := reader.readUint32()
	self.code = reader.readBytes(codeLength)
	self.exceptionTable = readExceptionTable(reader)
	self.attributes = readAttributes(reader,self.cp)
}

func readExceptionTable(reader *ClassReader) []*ExceptionTableEntry{
	exceptionTableLength := reader.readUint16()
	exceptionTable := make([] *ExceptionTableEntry, exceptionTableLength)
	for i := range exceptionTable {
		exceptionTable[i] = &ExceptionTableEntry{
			startPC:   reader.readUint16(),
			endPC:     reader.readUint16(),
			handlerPC: reader.readUint16(),
			catchType: reader.readUint16(),
		}
	}
	return exceptionTable
}

func (self *CodeAttribute) MaxLocals() uint16{
	return self.maxLocal
}

func (self *CodeAttribute) MaxStack() uint16{
	return self.maxStack
}

func (self *CodeAttribute) Code() []byte{
	return self.code
}