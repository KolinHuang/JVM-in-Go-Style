package classfile

import "fmt"

type ConstantPool []ConstantInfo

//读取所有的常量
func readConstantPool(reader *ClassReader) ConstantPool {
	//获取常量池计数
	cpCount := int(reader.readUint16())
	cp := make([]ConstantInfo, cpCount)
	//
	for i := 1; i < cpCount; i++ {
		cp[i] = readConstantInfo(reader, cp)
		//如果常量池遇到了这两种类型，就让计数+1，因为这两个类型需要占用两个entry
		switch cp[i].(type) {
		case *ConstantLongInfo, *ConstantDoubleInfo:
			i++
		}
	}
	return cp
}
//按索引查找常量
func (self ConstantPool) getConstantInfo(index uint16) ConstantInfo {
	if cpInfo := self[index]; cpInfo != nil {
		return cpInfo
	}
	panic(fmt.Errorf("Invalid constant pool index: %v!", index))
}

//返回名称和描述符
func (self ConstantPool) getNameAndType(index uint16) (string, string) {
	ntInfo := self.getConstantInfo(index).(*ConstantNameAndTypeInfo)
	name := self.getUtf8(ntInfo.nameIndex)
	_type := self.getUtf8(ntInfo.descriptorIndex)
	return name, _type
}

func (self ConstantPool) getClassName(index uint16) string {
	classInfo := self.getConstantInfo(index).(*ConstantClassInfo)
	return self.getUtf8(classInfo.nameIndex)
}

//根据字符引用索引找到字符
func (self ConstantPool) getUtf8(index uint16) string{
	utf8Info := self.getConstantInfo(index).(*ConstantUtf8Info)
	return utf8Info.str
}
