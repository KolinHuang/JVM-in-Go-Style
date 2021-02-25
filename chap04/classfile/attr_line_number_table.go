package classfile

type LineNumberTableAttribute struct {
	lineNumberTable []*LineNumberTableEntry
}

type LineNumberTableEntry struct {
	startPC	uint16
	lineNumber	uint16
}
//读取属性表数据
func (self *LineNumberTableAttribute) readInfo(reader *ClassReader){
	lineNumberTableLength := reader.readUint16()
	self.lineNumberTable = make([] *LineNumberTableEntry, lineNumberTableLength)
	for i := range self.lineNumberTable {
		self.lineNumberTable[i] = &LineNumberTableEntry{
			startPC:    reader.readUint16(),
			lineNumber: reader.readUint16(),
		}
	}
}
