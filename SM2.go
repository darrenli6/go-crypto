package main

import (
	"github.com/tjfoc/gmsm/sm2"
	"fmt"
)

// SM2 为非对称加密 基于椭圆加密,该算法已公开,签名速度与秘钥生成速度快于RSA

// 公钥密码

//func main1(){
//
//	//1 生成秘钥对
//	privateKey, err :=sm2.GenerateKey()
//
//	if err !=nil{
//		fmt.Println("秘钥对生成失败")
//		return
//	}
//
//    // 2从私钥中取出公钥
//
//    publicKey :=&privateKey.PublicKey
//
//    // 3 用公钥加密
//    msg := []byte("darren")
//    encrypto_msg,err :=publicKey.Encrypt( msg)
//	if err !=nil{
//		fmt.Println("公钥加密失败")
//		return
//	}else{
//
//		fmt.Printf("%x\n",encrypto_msg)
//	}
//
//
//	// 4 对私钥解密
//
//	decrypto_msg,err :=privateKey.Decrypt(encrypto_msg)
//
//	if err !=nil{
//		fmt.Println("私钥解密失败")
//		return
//	}else{
//		fmt.Printf("%x\n",decrypto_msg)
//	}
//
//
//
//}


// 验签名
func main(){


	err:=WriteKekPairToFile("private.pem","public.pem",[]byte("123456"))

	if err!=nil{
		fmt.Println("秘钥对写入文件失败!")
	}


}


//生成的公钥和私钥,写入文件中
func WriteKekPairToFile(privateKeyPath, pulicKeyPath string,password []byte) error{

	//1 生成秘钥对
	privateKey, err :=sm2.GenerateKey()
	if err !=nil{
		fmt.Println("秘钥对生成失败")
		return err
	}

    // 将私钥写入文件
	flag,err :=sm2.WritePrivateKeytoPem(privateKeyPath,privateKey,password)

    if !flag  || err!=nil{
    	fmt.Println("私钥写入文件失败")
    	return err
	}

	// 获取公钥
	publicKey := privateKey.Public().(*sm2.PublicKey)
   // publicKey := privateKey.PublicKey
   // 将公钥写入文件
	flag,err =sm2.WritePublicKeytoPem(pulicKeyPath,publicKey,nil)
	if !flag  || err!=nil{
		fmt.Println("公钥写入文件失败")
		return err
	}




	return nil

}