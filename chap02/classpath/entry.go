package classpath

import(
"os"
"strings"
)

//定义路径分隔符常量，":"
const pathListSepatator = string(os.PathListSeparator)
type Entry interface {
	//负责寻找和加载class文件
	readClass(className string) ([]byte, Entry, error)
	//负责返回变量的字符串表示，相当于Java的toString()方法
	String() string
}

func newEntry(path string) Entry {
	if strings.Contains(path, pathListSepatator) {
		return newCompositeEntry(path)
	}

	if strings.HasSuffix(path, "*") {
		return newWildcardEntry(path)
	}

	if strings.HasSuffix(path, ".jar") ||
		strings.HasSuffix(path, ".JAR") ||
		strings.HasSuffix(path, ".zip") ||
		strings.HasSuffix(path, ".ZIP") {
		return newZipEntry(path)
	}

	return newDirEntry(path)
}