package classfile


type MemberInfo struct {
	cp	ConstantPool	//常量池引用
	accessFlags	uint16	//访问标志
	nameIndex	uint16	//简单名称索引
	describetorIndex	uint16	//描述符索引
	attributes	[]AttributeInfo	//属性表集合
}
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

func (self *MemberInfo) AccessFlags() uint16 {
	return self.accessFlags
}

//根据简单名称索引读取简单名称
func (self *MemberInfo) Name() string {
	return self.cp.getUtf8(self.nameIndex)
}
//根据描述符索引读取描述符
func (self *MemberInfo) Descriptor() string {
	return self.cp.getUtf8(self.describetorIndex)
}

func (self *MemberInfo) CodeAttribute() *CodeAttribute{
	for _, attrInfo := range self.attributes {
		switch attrInfo.(type) {
		case *CodeAttribute:
			return attrInfo.(*CodeAttribute)
		}
	}
	return nil
}

func (self *MemberInfo) ConstantValueAttribute() *ConstantValueAttribute {
	for _, attrInfo := range self.attributes {
		switch attrInfo.(type) {
		case *ConstantValueAttribute:
			return attrInfo.(*ConstantValueAttribute)
		}
	}
	return nil
}