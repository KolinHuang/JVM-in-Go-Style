package classpath

import(
	"os"
	"path/filepath"
	"strings"
)
//WildcardEntry实际上也是CompositeEntry类型的
func newWildcardEntry(path string) CompositeEntry {
	baseDir := path[: len(path) - 1]	//删除"*"号
	compositeEntry := []Entry{}
	//匿名函数
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//通配符类路径不能递归匹配子目录下的JAR文件
		if info.IsDir() && path != baseDir {
			return filepath.SkipDir
		}
		//找出JAR文件
		if strings.HasSuffix(path, ".jar") ||
			strings.HasSuffix(path,".JAR") {
			jarEntry := newZipEntry(path)
			compositeEntry = append(compositeEntry, jarEntry)
		}
		return nil
	}
	//遍历baseDir，将所有JAR文件创建为ZipEntry，放入数组compositeEntry中
	//Walk walks the file tree rooted at root, calling walkFn for each file
	//or directory in the tree, including root.
	filepath.Walk(baseDir, walkFn)
	return compositeEntry
}