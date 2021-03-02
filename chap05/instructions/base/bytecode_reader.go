package base

//字节码读取器
type BytecodeReader struct {
	code []byte//存放字节码
	pc int//记录读取到了哪个字节
}


func (self *BytecodeReader) Reset(code []byte, pc int){
	self.code = code
	self.pc = pc
}
//读8比特，也就是一个字节
func (self *BytecodeReader) ReadUint8() uint8{
	data := self.code[self.pc]
	self.pc++
	return data
}

func (self *BytecodeReader) ReadInt8() int8{
	return int8(self.ReadUint8())
}
//连续读取两字节
func (self *BytecodeReader) ReadUint16() uint16{
	byte1 := uint16(self.ReadUint8())
	byte2 := uint16(self.ReadUint8())
	return (byte1 << 8) | byte2
}

func (self *BytecodeReader) ReadInt16() int16{
	return int16(self.ReadUint16())
}

//连续读取四字节
func (self *BytecodeReader) ReadInt32() int32{
	byte1 := int32(self.ReadUint8())
	byte2 := int32(self.ReadUint8())
	byte3 := int32(self.ReadUint8())
	byte4 := int32(self.ReadUint8())
	return (byte1 << 24) | (byte2 << 16) | (byte3 << 8) | byte4
}

func (self *BytecodeReader) ReadInt32s(n int32) []int32{
	ints := make([]int32, n)
	for i := range ints{
		ints[i] = self.ReadInt32()
	}
	return ints
}

func (self *BytecodeReader) SkipPadding() {
	for self.pc % 4 != 0 {//如果PC值不是4的倍数，就继续SKip padding
		self.ReadUint8()	//越过一个字节
	}
}

