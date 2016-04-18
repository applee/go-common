package crypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"os"
)

const (
	PKCS1 = iota
	PKCS8
)

// 加密
func RsaEncrypt(origData, pubKey []byte) ([]byte, error) {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext, privKey []byte, privKeyFormat int) ([]byte, error) {
	block, _ := pem.Decode(privKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	var (
		priv *rsa.PrivateKey
		err  error
	)
	if privKeyFormat == PKCS1 {
		priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	} else {
		privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		priv = privInterface.(*rsa.PrivateKey)
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

// 签名
func Sign(privKey []byte, hash crypto.Hash, hashed []byte, privKeyFormat int) ([]byte, error) {
	block, _ := pem.Decode(privKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	var (
		priv *rsa.PrivateKey
		err  error
	)
	if privKeyFormat == PKCS1 {
		priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	} else {
		privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		priv = privInterface.(*rsa.PrivateKey)
	}

	return rsa.SignPKCS1v15(nil, priv, hash, hashed)
}

//验签
func Verify(pubKey []byte, hash crypto.Hash, hashed, signature []byte) error {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return errors.New("public key error!")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA1, hashed, signature)
}

// PKCS8 私钥
func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) ([]byte, error) {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)

	return asn1.Marshal(info)
}

func GenKeyPair(bits int, privKeyFormat int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	var derStream []byte
	if privKeyFormat == PKCS1 {
		derStream = x509.MarshalPKCS1PrivateKey(privateKey)
	} else {
		derStream, _ = MarshalPKCS8PrivateKey(privateKey)
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}
