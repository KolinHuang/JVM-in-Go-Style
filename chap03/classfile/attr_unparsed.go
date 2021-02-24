package classfile

type UnparsedAttribute struct {
	name string
	length uint32
	info []byte
}

func (self *UnparsedAttribute) readInfo(reader *ClassReader){
	//读取length长度的字节放入info中作为属性值
	self.info = reader.readBytes(self.length)
}
