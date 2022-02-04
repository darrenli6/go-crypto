package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// 椭圆加密
func main(){

	//生成秘钥对
	privateKey, _ :=ecdsa.GenerateKey(elliptic.P384(),rand.Reader)

	msg := []byte("hello")
	digest :=sha256.Sum256(msg)

	// 签名
	r,s, _ := ecdsa.Sign(rand.Reader,privateKey,digest[:])

	pub := privateKey.PublicKey

	flag := ecdsa.Verify(&pub,digest[:],r,s)
	if flag{
		fmt.Println("验证成功")
	}else{
		fmt.Println("验证失败")
	}


}
