package common

import (
	"bytes"
	"crypto/aes"
)

var Key = []byte("abcdabcdabcdabcdabcdabcdabcdabcd")

// ECB模式加密

func ECBEncrypt(src, key []byte) ([]byte, error) {

	block, _ := aes.NewCipher(key)
	src = PKCS5Padding(src, block.BlockSize())
	//制作成16的倍数
	var dst []byte
	tmpData := make([]byte, block.BlockSize())
	for index := 0; index < len(src); index += block.BlockSize() {
		block.Encrypt(tmpData, src[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}
	return dst, nil
}

// ECB模式解密
func ECBDecrypt(crypted, key []byte) ([]byte, error) {

	block, _ := aes.NewCipher(key)

	var dst []byte
	tmpData := make([]byte, block.BlockSize())

	for index := 0; index < len(crypted); index += block.BlockSize() {
		block.Decrypt(tmpData, crypted[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}
	dst, _ = PKCS5UnPadding(dst)
	//制作成原长度
	return dst, nil
}

// PKCS5填充
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 去除PKCS5填充
func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)], nil
}
