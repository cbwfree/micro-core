package fn

import (
	"os"
)

// ExistFile 检查目录/文件是否存在
func FileExist(folder string) bool {
	if _, err := os.Stat(folder); !os.IsNotExist(err) {
		return true
	}
	return false
}

// MkDir 创建目录
func Mkdir(folder string, perm ...os.FileMode) error {
	if FileExist(folder) {
		return nil
	}

	var p os.FileMode
	if len(perm) > 0 {
		p = perm[0]
	} else {
		p = 0755
	}

	return os.MkdirAll(folder, p)
}
