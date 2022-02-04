package main

import (
	"crypto/rsa"
	"crypto/rand"
	"crypto/md5"
	"crypto"
	"fmt"
)

func main(){

	//生成秘钥对
	priv,_ := rsa.GenerateKey(rand.Reader,2048)

	//消息
	msg := []byte("hello baby")
	h :=md5.New()
	h.Write(msg)
	result :=h.Sum(nil)

	opts:= &rsa.PSSOptions{SaltLength:rsa.PSSSaltLengthAuto,Hash:crypto.MD5}

	// 签名
	sig,_ := rsa.SignPSS(rand.Reader,priv,crypto.MD5,result,opts)


	//msg1 := []byte("hello baby1")
	//h1 :=md5.New()
	//h1.Write(msg1)
	//result1 :=h1.Sum(nil)

	//获取公钥
	pub := &priv.PublicKey

	err := rsa.VerifyPSS(pub,crypto.MD5,result,sig,opts)
	if err!=nil{
		fmt.Println("验证失败")
	}else{
		fmt.Println("验证成功")
	}


}