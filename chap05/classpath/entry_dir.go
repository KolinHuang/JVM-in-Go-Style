package classpath

import(
	"io/ioutil"
	"path/filepath"
)
/**
目录形式的类路径
 */
type DirEntry struct {
	absDir	string
}

//返回目录形式类路径项实例
func newDirEntry(path string) *DirEntry {
	//Abs函数返回路径的绝对路径表示
	absDir, err := filepath.Abs(path)
	if err != nil {
		//当调用了panic函数，此函数的执行将会停止
		panic(err)
	}
	//返回DirEntry实例
	return &DirEntry{absDir}
}

func (self *DirEntry) readClass(className string) ([]byte, Entry, error){
	//把目录和class文件名拼成一个完整的路径
	fileName := filepath.Join(self.absDir, className)
	//读取class文件的内容
	data, err := ioutil.ReadFile(fileName)
	return data, self, err
}

func (self *DirEntry) String() string {
	return self.absDir
}
