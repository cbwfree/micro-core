package fn

import (
	"os"
)

// ExistDir 检查目录是否存在
func ExistDir(dir string) bool {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return true
	}
	return false
}

// MkDir 创建目录
func MkDir(dir string, perm ...os.FileMode) error {
	var p os.FileMode
	if len(perm) > 0 {
		p = perm[0]
	} else {
		p = 0755
	}
	return os.MkdirAll(dir, p)
}

func InitFolder(folder string) error {
	if !ExistDir(folder) {
		if err := MkDir(folder); err != nil {
			return err
		}
	}
	return nil
}
