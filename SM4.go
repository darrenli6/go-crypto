package main

import (
	"bytes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"github.com/tjfoc/gmsm/sm4"
)

//填充最后一个分组
//src:待填充的数据  blockSize:每组大小
func PaddingText(src []byte, blockSize int) []byte {
	//求出最后一个分组需要填充的字节数
	padding := blockSize - len(src)%blockSize
	//准备填充数据
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	//拼接数据
	nextText := append(src, padText...)
	return nextText
}

//去除尾部填充数据
func UnPaddingText(src []byte) []byte {
	length := len(src)
	number := int(src[length-1])
	newText := src[:length-number]
	return newText
}

// sm4是分组加密算法
func EncryptSm4(src, key []byte) []byte {

	// 1 创建加密块
	block, err := sm4.NewCipher(key)
	if err != nil {
		panic(err)
	}

	//2 填充数据
	src = PaddingText(src, block.BlockSize())

	//3 设置加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])

	//4 加密

	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)
	return dst
}

func DecryptSm4(src, key []byte) []byte {
	//1 创建解密的块
	block, err := sm4.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// 2设置解密模式

	blockMode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])

	// 3 解密

	dst := make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)
	// 4 去除尾部填充
	src = UnPaddingText(dst)

	return src

}

func main() {

	key := []byte("1234567891234560")

	msg := []byte("北京")

	encrypt_msg := EncryptSm4(msg, key)

	fmt.Println(" encrypto msg :", hex.EncodeToString(encrypt_msg))

	descript_msg := DecryptSm4(encrypt_msg, key)

	fmt.Println("decrypt_msg = ", string(descript_msg))
}
