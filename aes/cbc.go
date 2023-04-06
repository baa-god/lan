package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

type CBC struct {
	Text string
	// 必须是 16(AES-128)、24(AES-192) 或 32(AES-256) 位的 AES 密钥
	Secret string
}

// Encrypt AES-CBC 加密
func (c CBC) Encrypt() (ciphertext string, err error) {
	block, err := aes.NewCipher([]byte(c.Secret))
	if err != nil {
		return
	}

	blockSize := len(c.Secret)
	padding := blockSize - len(c.Text)%blockSize // 填充字节
	if padding == 0 {
		padding = blockSize
	}

	// 填充 padding 个 byte(padding) 到 plaintext
	plaintext := c.Text + string(bytes.Repeat([]byte{byte(padding)}, padding))
	cipherBytes := make([]byte, aes.BlockSize+len(plaintext))

	// 初始向量 iv 为随机的 16 位字符串 (必须是16位)
	// 解密需要用到这个相同的 iv，因此将它包含在密文的开头。
	iv := cipherBytes[:aes.BlockSize]
	if _, err = rand.Read(iv); err != nil { // 将同时写到 cipherBytes 的开头
		return
	}

	if len(iv) != block.BlockSize() {
		err = errors.New("IV length must equal block size")
		return
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherBytes[aes.BlockSize:], []byte(plaintext))

	ciphertext = base64.StdEncoding.EncodeToString(cipherBytes)
	return
}

// Decrypt AES-CBC 解密
func (c CBC) Decrypt(ciphertext string) (plaintext string, err error) {
	block, err := aes.NewCipher([]byte(c.Secret))
	if err != nil {
		return
	}

	ciphercode, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return
	}

	iv := ciphercode[:aes.BlockSize]        // 密文的前 16 个字节为 iv
	ciphercode = ciphercode[aes.BlockSize:] // 正式密文

	if len(iv) != block.BlockSize() {
		err = errors.New("IV length must equal block size")
		return
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphercode, ciphercode)

	plaintext = string(ciphercode) // ↓ 减去 padding
	plaintext = plaintext[:len(plaintext)-int(plaintext[len(plaintext)-1])]

	return
}
