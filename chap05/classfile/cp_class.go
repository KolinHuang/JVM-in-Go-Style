package classfile

type ConstantClassInfo struct {
	cp	ConstantPool
	nameIndex	uint16
}

//读取全限定名常量项的索引
func(self *ConstantClassInfo) readInfo(reader *ClassReader){
	self.nameIndex = reader.readUint16()
}

//读取索引对应的字符串
func (self *ConstantClassInfo) Name() string{
	return self.cp.getUtf8(self.nameIndex)
}
