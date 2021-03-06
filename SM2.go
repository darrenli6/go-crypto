package main

import (
	"crypto/rand"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
	"io/ioutil"
)

// SM2 为非对称加密 基于椭圆加密,该算法已公开,签名速度与秘钥生成速度快于RSA

// 公钥密码

// mod与git版本的问题
/*
 go.mod 是将pack存放到 /go/pkg/ 下面
*/

func main1() {
	//1.生成密钥对
	privateKey, err := sm2.GenerateKey()
	if err != nil {
		fmt.Println("秘钥生成失败!")
		return
	}

	//2.从私钥中取出公钥
	publicKey := &privateKey.PublicKey

	//3.公钥加密
	msg := []byte("darren")
	encrypt_msg, err := publicKey.Encrypt(msg)
	if err != nil {
		fmt.Println("公钥加密失败!")
		return
	} else {
		fmt.Printf("%x\n", encrypt_msg)
	}

	//4.私钥解密
	decrypt_msg, err := privateKey.Decrypt(encrypt_msg)
	if err != nil {
		fmt.Println("私钥解密失败!")
		return
	} else {
		fmt.Println(string(decrypt_msg))
	}

}

// 验签名
func main() {

	//err := WriteKekPairToFile("private.pem", "public.pem", []byte("123456"))
	//
	//if err != nil {
	//	fmt.Println("秘钥对写入文件失败!")
	//}

	privateKey, publicKey, err := ReadKeyPairFromFile("private.pem", "public.pem", []byte("123456"))
	if err != nil {
		fmt.Println("读取失败")
		return
	}

	file, err := ioutil.ReadFile("D:\\Evc\\test.mp4")
	if err != nil {
		fmt.Println("文件读取失败")
		return
	}

	// 签名
	sig_msg, err := privateKey.Sign(rand.Reader, file, nil)
	if err != nil {
		fmt.Println("签名失败")
		return
	}

	file, err = ioutil.ReadFile("D:\\Evc\\2.txt")
	if err != nil {
		fmt.Println("文件读取失败")
		return
	}
	// 验证签名
	flag := publicKey.Verify(file, sig_msg)
	if flag {
		fmt.Println("验签成功")
	} else {
		fmt.Println("验签失败")
	}
}

// 从文件中读取公钥和私钥

func ReadKeyPairFromFile(privateKeyPath, publicKeyPath string, password []byte) (*sm2.PrivateKey, *sm2.PublicKey, error) {

	privatekey, err := sm2.ReadPrivateKeyFromPem(privateKeyPath, password)

	if err != nil {
		fmt.Println("私钥文件读取失败")
		return nil, nil, err
	}

	//读取公钥文件

	publicKey, err := sm2.ReadPublicKeyFromPem(publicKeyPath, nil)

	if err != nil {
		fmt.Println("公钥文件读取失败")
		return nil, nil, err
	}

	return privatekey, publicKey, nil

}

//生成的公钥和私钥,写入文件中
func WriteKekPairToFile(privateKeyPath, pulicKeyPath string, password []byte) error {

	//1 生成秘钥对
	privateKey, err := sm2.GenerateKey()
	if err != nil {
		fmt.Println("秘钥对生成失败")
		return err
	}

	// 将私钥写入文件
	flag, err := sm2.WritePrivateKeytoPem(privateKeyPath, privateKey, password)

	if !flag || err != nil {
		fmt.Println("私钥写入文件失败")
		return err
	}

	// 获取公钥
	publicKey := privateKey.Public().(*sm2.PublicKey)
	// publicKey := privateKey.PublicKey
	// 将公钥写入文件
	flag, err = sm2.WritePublicKeytoPem(pulicKeyPath, publicKey, nil)
	if !flag || err != nil {
		fmt.Println("公钥写入文件失败")
		return err
	}

	return nil

}
