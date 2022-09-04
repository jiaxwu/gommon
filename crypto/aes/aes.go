package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/jiaxwu/gommon/crypto"
)

// AES+CBC+PKCS5Padding加密，添加一些便捷方法
type Cipher struct {
	block cipher.Block
	iv    []byte
}

func New(key []byte) *Cipher {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return &Cipher{
		block: block,
		iv:    key[:block.BlockSize()],
	}
}

// 加密
func (c *Cipher) Encrypt(src []byte) []byte {
	encrypter := cipher.NewCBCEncrypter(c.block, c.iv)
	src = crypto.PKCS5Padding(src, c.BlockSize())
	dst := make([]byte, len(src))
	encrypter.CryptBlocks(dst, src)
	return dst
}

// 加密到Base64
func (c *Cipher) EncryptToBase64(src string) string {
	dst := c.Encrypt([]byte(src))
	return base64.RawStdEncoding.EncodeToString(dst)
}

// 解密
func (c *Cipher) Decrypt(src []byte) []byte {
	decrypter := cipher.NewCBCDecrypter(c.block, c.iv)
	dst := make([]byte, len(src))
	decrypter.CryptBlocks(dst, src)
	return crypto.PKCS5Trimming(dst)
}

// 从Base64解密
func (c *Cipher) DecryptFromBase64(src string) (string, error) {
	dst, err := base64.RawStdEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	return string(c.Decrypt(dst)), nil
}

func (c *Cipher) BlockSize() int {
	return c.block.BlockSize()
}
