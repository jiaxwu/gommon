package crypto

import (
	"fmt"
	"math"
)

// 填充
func PKCS5Padding(src []byte, blockSize int) []byte {
	if blockSize > math.MaxUint8 {
		panic("too large block size")
	}
	paddingLen := blockSize - len(src)%blockSize
	dst := make([]byte, len(src)+paddingLen)
	copy(dst, src)
	for i := len(src); i < len(dst); i++ {
		dst[i] = byte(paddingLen)
	}
	return dst
}

// 移除填充
func PKCS5Trimming(src []byte) []byte {
	fmt.Println(len(src))
	paddingLen := src[len(src)-1]
	return src[:len(src)-int(paddingLen)]
}
