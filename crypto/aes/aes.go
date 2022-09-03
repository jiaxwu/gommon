package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/jiaxwu/gommon/crypto"
)

// AES+CBC加密，添加一些便捷方法
type Cipher struct {
	block     cipher.Block
	encrypter cipher.BlockMode
	decrypter cipher.BlockMode
}

func New(key []byte) *Cipher {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypter := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	decrypter := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	return &Cipher{
		block:     block,
		encrypter: encrypter,
		decrypter: decrypter,
	}
}

// 加密
func (c *Cipher) Encrypt(src []byte) []byte {
	src = crypto.PKCS5Padding(src, c.block.BlockSize())
	dst := make([]byte, len(src))
	c.encrypter.CryptBlocks(dst, src)
	return dst
}

// 加密到Base64
func (c *Cipher) EncryptToBase64(src string) string {
	b := []byte(src)
	dst := c.Encrypt([]byte(b))
	return base64.RawStdEncoding.EncodeToString(dst)
}

// 解密
func (c *Cipher) Decrypt(src []byte) []byte {
	dst := make([]byte, len(src))
	c.decrypter.CryptBlocks(dst, src)
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
