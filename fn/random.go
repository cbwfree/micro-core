package fn

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	mRand "math/rand"
	"time"
)

func init() {
	mRand.Seed(time.Now().UnixNano())
}

// GenRandStr 生成随机字符串
func GenRandStr(length int) (string, error) {
	b := make([]byte, length/2)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return "", errors.New("could not successfully read from the system CSPRNG")
	}
	return hex.EncodeToString(b), nil
}

// GetRandom 在指定范围内生成随机数
func GetRandom(min, max int64) int64 {
	return min + mRand.Int63n(max-min)
}
