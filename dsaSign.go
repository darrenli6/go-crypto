package main

import (
	"crypto/dsa"
	"crypto/rand"
	"fmt"
)

//直接对消息加密
func main(){
  // Paramter代表私钥的参数

  var param dsa.Parameters

  //GenerateParameters 函数随机的设置合法的参数到param
  // 根据第三个参数就决定了L和N长度,长度越长,加密程度越高
  dsa.GenerateParameters(&param,rand.Reader,dsa.L2048N256)

  var priv dsa.PrivateKey

  priv.Parameters=param

  // 生成秘钥对
  dsa.GenerateKey(&priv,rand.Reader)

  message := []byte("hello darren")
  // 使用私钥对消息直接签名

  r,s, _:=dsa.Sign(rand.Reader,&priv,message)

  pub :=priv.PublicKey

  //message=[]byte("123")
  if dsa.Verify(&pub,message,r,s){
  	fmt.Println("验证通过")

  }else{
  	fmt.Println("验证不通过!")
  }

}
