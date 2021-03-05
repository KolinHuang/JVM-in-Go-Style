package classfile

type ConstantNameAndTypeInfo struct {
	cp ConstantPool
	nameIndex	uint16	//字段或方法名索引
	descriptorIndex	uint16	//字段或方法描述符索引
}

func(self *ConstantNameAndTypeInfo) readInfo(reader *ClassReader){
	self.nameIndex = reader.readUint16()
	self.descriptorIndex = reader.readUint16()
}