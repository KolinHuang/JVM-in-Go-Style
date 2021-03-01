package classfile

type AttributeInfo interface {
	readInfo(reader *ClassReader)
}

//读取属性表
func readAttributes(reader *ClassReader, cp ConstantPool) []AttributeInfo {
	attributesCount := reader.readUint16()//属性表计数
	attributes := make([]AttributeInfo, attributesCount)//根据属性表计数创建属性表数组
	for i := range attributes {
		attributes[i] = readAttribute(reader, cp)
	}
	return attributes
}
//读取单个属性
func readAttribute(reader *ClassReader, cp ConstantPool) AttributeInfo {
	attrNameIndex := reader.readUint16()//读取属性名索引
	attrName := cp.getUtf8(attrNameIndex)//从常量池中找到属性名
	attrLen := reader.readUint32()//读取属性长度
	attrInfo := newAttributeInfo(attrName, attrLen, cp)//创建具体的属性实例
	attrInfo.readInfo(reader)
	return attrInfo
}
//根据属性名创建属性实例
func newAttributeInfo(attrName string, attrLen uint32, cp ConstantPool) AttributeInfo {
	switch attrName {
	case "Code":
		return &CodeAttribute{cp : cp}//
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

