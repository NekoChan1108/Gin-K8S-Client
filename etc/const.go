package etc

import (
	"path/filepath"
	"runtime"
)

const YamlName = "config"

var YamlPath string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	// 获取当前文件所在目录
	currentDir := filepath.Dir(filename)
	// 构造etc目录的绝对路径
	YamlPath = filepath.Join(currentDir)
}
