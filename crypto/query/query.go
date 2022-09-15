package query

import (
	"encoding/base64"
	"errors"

	"github.com/jiaxwu/gommon/crypto"
	"github.com/jiaxwu/gommon/crypto/aes"
)

// 不满一个切片
var ErrNotFullSlice = errors.New("not full slice")

// 支持模糊查询的加密
// 基于AES+CBC+PKCS5Padding
type Cipher struct {
	sliceSize int // 切片长度
	c         *aes.Cipher
}

func New(key []byte, sliceSize int) *Cipher {
	if sliceSize <= 0 {
		panic("slice size must be greater than 0")
	}
	return &Cipher{
		sliceSize: sliceSize,
		c:         aes.New(key),
	}
}

// 加密
func (c *Cipher) Encrypt(src []byte) ([]byte, error) {
	return c.encrypt(src, false)
}

// 加密到Base64
func (c *Cipher) EncryptToBase64(src string) (string, error) {
	dst, err := c.encrypt([]byte(src), true)
	if err != nil {
		return "", err
	}
	return string(dst), nil
}

func (c *Cipher) encrypt(src []byte, isBase64 bool) ([]byte, error) {
	srcLen := len(src)
	if srcLen < c.sliceSize {
		return nil, ErrNotFullSlice
	}
	sliceCnt := srcLen - c.sliceSize + 1           // 切片数量
	paddingSliceLen := c.paddingSliceLen(isBase64) // 填充后切片长度
	dstLen := sliceCnt * paddingSliceLen           // 目标串长度
	dst := make([]byte, dstLen)
	for i := 0; i+c.sliceSize <= srcLen; i++ {
		encrypt := c.c.Encrypt(src[i : i+c.sliceSize])
		if isBase64 {
			base64.RawStdEncoding.Encode(dst[paddingSliceLen*i:], encrypt)
		} else {
			copy(dst[paddingSliceLen*i:], encrypt)
		}
	}
	return dst, nil
}

// 解密
func (c *Cipher) Decrypt(src []byte) ([]byte, error) {
	return c.decrypt(src, false)
}

// 从Base64解密
func (c *Cipher) DecryptFromBase64(src string) (string, error) {
	dst, err := c.decrypt([]byte(src), true)
	return string(dst), err
}

func (c *Cipher) decrypt(src []byte, isBase64 bool) ([]byte, error) {
	paddingSliceLen := c.paddingSliceLen(isBase64) // 填充后切片长度
	srcLen := len(src)
	if srcLen < paddingSliceLen || srcLen%paddingSliceLen != 0 {
		return nil, ErrNotFullSlice
	}
	sliceCnt := srcLen / paddingSliceLen // 切片数量
	dstLen := sliceCnt + c.sliceSize - 1 // 目标串长度
	dst := make([]byte, dstLen)
	var decoded []byte
	if isBase64 {
		decoded = make([]byte, base64.RawStdEncoding.DecodedLen(paddingSliceLen))
	}
	for i := 0; i < sliceCnt; i++ {
		if isBase64 {
			_, err := base64.RawStdEncoding.Decode(decoded, src[paddingSliceLen*i:paddingSliceLen*(i+1)])
			if err != nil {
				return nil, err
			}
			copy(dst[i:], c.c.Decrypt(decoded))
		} else {
			copy(dst[i:], c.c.Decrypt(src[paddingSliceLen*i:paddingSliceLen*(i+1)]))
		}
	}
	return dst, nil
}

// 填充后切片长度
func (c *Cipher) paddingSliceLen(isBase64 bool) int {
	paddingSliceLen := crypto.PKCS5DstLen(c.sliceSize, c.c.BlockSize())
	if isBase64 {
		paddingSliceLen = base64.RawStdEncoding.EncodedLen(paddingSliceLen)
	}
	return paddingSliceLen
}
