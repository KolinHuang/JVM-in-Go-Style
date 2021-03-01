package classfile

import "encoding/binary"

//[]byte的包装类型
type ClassReader struct {
	data []byte
}

//读取u1
func (self *ClassReader) readUint8() uint8 {
	val := self.data[0]
	//跳过1字节
	self.data = self.data[1:]
	return val
}

//读取u2
func (self *ClassReader) readUint16() uint16 {
	//用binary包以大端格式读取2字节数据
	val := binary.BigEndian.Uint16(self.data)
	self.data = self.data[2:]
	return val
}

//读取u4
func (self *ClassReader) readUint32() uint32 {
	//用binary包以大端格式读取4字节数据
	val := binary.BigEndian.Uint32(self.data)
	self.data = self.data[4:]
	return val
}

//读取u8
func (self *ClassReader) readUint64() uint64 {
	//用binary包以大端格式读取4字节数据
	val := binary.BigEndian.Uint64(self.data)
	self.data = self.data[8:]
	return val
}

//读取u2集合
func (self *ClassReader) readUint16s() []uint16 {
	n := self.readUint16()//集合长度
	s := make([] uint16, n)
	for i := range s {
		s[i] = self.readUint16()
	}
	return s
}

//读取指定数量的字节
func (self *ClassReader) readBytes(n uint32) []byte {
	bytes := self.data[:n]
	self.data = self.data[n:]
	return bytes
}

