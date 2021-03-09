package classpath

import(
	"archive/zip"	//提供读取和写入ZIP压缩包的操作
	"errors"
	"io/ioutil"	//IO工具包，提供一些IO操作
	"path/filepath"
)
/**
压缩文件形式的类路径
 */
type ZipEntry struct {
	absPath string	//存放JAR/ZIP文件的绝对路径
}

func newZipEntry(path string) *ZipEntry {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return &ZipEntry{absPath}
}

//从ZIP文件中提取class文件
func (self *ZipEntry) readClass(className string) ([]byte, Entry, error) {
	//OpenReader will open the Zip file specified by name and return a ReadCloser.
	r, err := zip.OpenReader(self.absPath)
	if err != nil {
		return nil, nil, err
	}
	//defer类似于Java中的finally语句块，在函数返回之前，或者说在return语句后执行。
	defer r.Close()
	//遍历压缩包内的文件
	for _, f := range r.File {
		//如果找到了className对应的class文件，就打开并读取内容，然后返回
		if f.Name == className {
			//打开文件
			rc, err := f.Open()
			if err != nil {//打开失败，返回错误
				return nil, nil, err
			}
			defer rc.Close()
			//读取文件所有数据，返回字节数组
			data, err := ioutil.ReadAll(rc)
			if err != nil {//读取失败
				return nil, nil, err
			}
			return data, self, nil
		}
	}
	return nil, nil, errors.New("class not found: " + className)

}

func (self *ZipEntry) String() string {
	return self.absPath
}
