package utils

import "bytes"

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

//5
//12345 60000
func ZeroPadding(data []byte, blockSize int) []byte {
	//求出最后一个分组需要填充的字节数
	padding := blockSize - len(data)%blockSize
	//准备填充数据
	padText := bytes.Repeat([]byte{byte(0)}, padding)
	return append(data, padText...)
}

func ZeroUnPadding(data []byte) []byte {
	//true：去除右侧数据
	return bytes.TrimRightFunc(data, func(r rune) bool {
		return r == rune(0)
	})
}
