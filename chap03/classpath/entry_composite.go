package classpath

import(
	"errors"
	"strings"
)
/**
多个Entry构成的类路径
 */
type CompositeEntry []Entry
//将路径列表参数按分隔符分成小路径，然后将每个小路径转化为具体的Entry实例
func newCompositeEntry(pathList string) CompositeEntry {
	compositeEntry := []Entry{}
	for _, path := range strings.Split(pathList, pathListSepatator) {
		//调用Entry接口中的newEntry函数
		entry := newEntry(path)
		compositeEntry = append(compositeEntry, entry)
	}
	return compositeEntry
}

func (self CompositeEntry) readClass(className string) ([]byte, Entry, error) {
	for _, entry := range self {
		data, from, err := entry.readClass(className)
		if err == nil {
			return data, from, nil
		}
	}
	return nil, nil, errors.New("class not found: " + className)
}
//调用每一个子路径的String方法，然后把得到的字符串用路径分隔符拼接起来即可
func (self CompositeEntry) String() string {
	strs := make([]string, len(self))

	for i, entry := range self {
		strs[i] = entry.String()
	}
	return strings.Join(strs, pathListSepatator)
}