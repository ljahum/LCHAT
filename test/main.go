package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

type Box struct {
	A int    `json:"A"`
	B string `json:"B"`
}
type Comment struct {
	Id      int64
	Name    string
	Content string
	Mail    string
	Time    string
}
type StatusFlow struct {
	//payload和sign默认b64传输
	Status  int    `json:"Status"`
	Payload string `json:"Payload"` // b64
}

func main() {
	sqlStr := "SELECT * FROM mail_table WHERE `to` = '" + "user1" + "' ORDER BY id DESC LIMIT 20;"
	println(sqlStr)

}

// pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

func AesEncrypt(data []byte, key []byte) ([]byte, error) {
	//创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//判断加密快的大小
	blockSize := block.BlockSize()
	//填充
	encryptBytes := pkcs7Padding(data, blockSize)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

// AesDecrypt 解密
func AesDecrypt(data []byte, key []byte) []byte {
	//创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	//初始化解密数据接收切片
	crypted := make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(crypted, data)
	//去除填充
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil
	}
	return crypted
}

// EncryptByAes Aes加密 后 base64 再加
func EncryptecbMode_withPadding(data []byte, PwdKey []byte) []byte {
	res, err := AesEncrypt(data, PwdKey)
	if err != nil {
		return nil
	}
	return []byte(base64.StdEncoding.EncodeToString(res))
}

// DecryptByAes Aes 解密
func DecryptecbMode_withUnpadding(data []byte, PwdKey []byte) []byte {
	dataByte, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil
	}
	return AesDecrypt(dataByte, PwdKey)
}
