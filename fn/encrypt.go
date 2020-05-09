// 加密处理
package fn

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
)

// MD5 使用MD5对数据签名 (长度32)
func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// Sha256 使用Sha256对数据签名 (长度64)
func Sha256(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// UUID4 生成唯一ID
func UUID() string {
	return uuid.New().String()
}
