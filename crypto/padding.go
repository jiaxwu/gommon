package crypto

import (
	"math"
)

// 填充
func PKCS5Padding(src []byte, blockSize int) []byte {
	if blockSize > math.MaxUint8 {
		panic("too large block size")
	}
	srcLen := len(src)
	paddingLen := blockSize - srcLen%blockSize
	dst := make([]byte, srcLen+paddingLen)
	copy(dst, src)
	for i := len(src); i < len(dst); i++ {
		dst[i] = byte(paddingLen)
	}
	return dst
}

// 移除填充
func PKCS5Trimming(src []byte) []byte {
	paddingLen := src[len(src)-1]
	return src[:len(src)-int(paddingLen)]
}

// 填充目标长度
func PKCS5DstLen(srcLen, blockSize int) int {
	paddingLen := blockSize - srcLen%blockSize
	return srcLen + paddingLen
}
