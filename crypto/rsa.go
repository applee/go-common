package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"io"
	"math/big"
)

const (
	PKCS1 = iota
	PKCS8
)

func GetPrivKey(privKey []byte, format int) (priv *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(privKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	if format == PKCS1 {
		priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return
		}
	} else {
		var iface interface{}
		iface, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return
		}
		priv = iface.(*rsa.PrivateKey)
	}
	return
}

func GetPubKey(pubKey []byte) (pub *rsa.PublicKey, err error) {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	pub = pubInterface.(*rsa.PublicKey)
	return
}

// 加密
func RsaEncrypt(origData, pubKey []byte) ([]byte, error) {
	pub, err := GetPubKey(pubKey)
	if err != nil {
		return nil, err
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext, privKey []byte, privKeyFormat int) ([]byte, error) {
	priv, err := GetPrivKey(privKey, privKeyFormat)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

// 签名
func Sign(privKey []byte, hash crypto.Hash, hashed []byte, privKeyFormat int) ([]byte, error) {
	priv, err := GetPrivKey(privKey, privKeyFormat)
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(nil, priv, hash, hashed)
}

// 验签
func Verify(pubKey []byte, hash crypto.Hash, hashed, signature []byte) error {
	pub, err := GetPubKey(pubKey)
	if err != nil {
		return err
	}
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

// 私钥加密
func RsaEncryptPrivate(origData, privKey []byte, privKeyFormat int) (
	[]byte, error) {
	priv, err := GetPrivKey(privKey, privKeyFormat)
	if err != nil {
		return nil, err
	}
	tLen := len(origData)
	k := (priv.N.BitLen() + 7) / 8
	if k < tLen+11 {
		return nil, rsa.ErrMessageTooLong
	}
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-tLen:k], origData)
	m := new(big.Int).SetBytes(em)
	c, err := decrypt(rand.Reader, priv, m)
	if err != nil {
		return nil, err
	}
	copyWithLeftPad(em, c.Bytes())
	return em, nil
}

// 公钥解密
func RsaDecryptPublic(origData, pubKey []byte) ([]byte, error) {
	pub, err := GetPubKey(pubKey)
	if err != nil {
		return nil, err
	}
	k := (pub.N.BitLen() + 7) / 8
	if k != len(origData) {
		return nil, rsa.ErrVerification
	}
	m := new(big.Int).SetBytes(origData)
	if m.Cmp(pub.N) > 0 {
		return nil, rsa.ErrVerification
	}
	m.Exp(m, big.NewInt(int64(pub.E)), pub.N)

	d := leftPad(m.Bytes(), k)

	return leftUnPad(d)
}

func GenKeyPair(bits, format int, priv, pub io.Writer) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	var derStream []byte
	if format == PKCS1 {
		derStream = x509.MarshalPKCS1PrivateKey(privateKey)
	} else {
		derStream, _ = MarshalPKCS8PrivateKey(privateKey)
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	err = pem.Encode(priv, block)
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
	err = pem.Encode(pub, block)
	if err != nil {
		return err
	}
	return nil
}

var bigZero = big.NewInt(0)
var bigOne = big.NewInt(1)

// decrypt performs an RSA decryption, resulting in a plaintext integer. If a
// random source is given, RSA blinding is used.
func decrypt(random io.Reader, priv *rsa.PrivateKey, c *big.Int) (m *big.Int, err error) {
	// TODO(agl): can we get away with reusing blinds?
	if c.Cmp(priv.N) > 0 {
		err = rsa.ErrDecryption
		return
	}
	if priv.N.Sign() == 0 {
		return nil, rsa.ErrDecryption
	}

	var ir *big.Int
	if random != nil {
		// Blinding enabled. Blinding involves multiplying c by r^e.
		// Then the decryption operation performs (m^e * r^e)^d mod n
		// which equals mr mod n. The factor of r can then be removed
		// by multiplying by the multiplicative inverse of r.

		var r *big.Int

		for {
			r, err = rand.Int(random, priv.N)
			if err != nil {
				return
			}
			if r.Cmp(bigZero) == 0 {
				r = bigOne
			}
			var ok bool
			ir, ok = modInverse(r, priv.N)
			if ok {
				break
			}
		}
		bigE := big.NewInt(int64(priv.E))
		rpowe := new(big.Int).Exp(r, bigE, priv.N) // N != 0
		cCopy := new(big.Int).Set(c)
		cCopy.Mul(cCopy, rpowe)
		cCopy.Mod(cCopy, priv.N)
		c = cCopy
	}

	if priv.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, priv.D, priv.N)
	} else {
		// We have the precalculated values needed for the CRT.
		m = new(big.Int).Exp(c, priv.Precomputed.Dp, priv.Primes[0])
		m2 := new(big.Int).Exp(c, priv.Precomputed.Dq, priv.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, priv.Primes[0])
		}
		m.Mul(m, priv.Precomputed.Qinv)
		m.Mod(m, priv.Primes[0])
		m.Mul(m, priv.Primes[1])
		m.Add(m, m2)

		for i, values := range priv.Precomputed.CRTValues {
			prime := priv.Primes[2+i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}

	if ir != nil {
		// Unblind.
		m.Mul(m, ir)
		m.Mod(m, priv.N)
	}

	return
}

// copyWithLeftPad copies src to the end of dest, padding with zero bytes as
// needed.
func copyWithLeftPad(dest, src []byte) {
	numPaddingBytes := len(dest) - len(src)
	for i := 0; i < numPaddingBytes; i++ {
		dest[i] = 0
	}
	copy(dest[numPaddingBytes:], src)
}

// modInverse returns ia, the inverse of a in the multiplicative group of prime
// order n. It requires that a be a member of the group (i.e. less than n).
func modInverse(a, n *big.Int) (ia *big.Int, ok bool) {
	g := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)
	g.GCD(x, y, a, n)
	if g.Cmp(bigOne) != 0 {
		// In this case, a and n aren't coprime and we cannot calculate
		// the inverse. This happens because the values of n are nearly
		// prime (being the product of two primes) rather than truly
		// prime.
		return
	}

	if x.Cmp(bigOne) < 0 {
		// 0 is not the multiplicative inverse of any element so, if x
		// < 1, then x is negative.
		x.Add(x, n)
	}

	return x, true
}

// leftPad returns a new slice of length size. The contents of input are right
// aligned in the new slice.
func leftPad(input []byte, size int) (out []byte) {
	n := len(input)
	if n > size {
		n = size
	}
	out = make([]byte, size)
	copy(out[len(out)-n:], input)
	return
}

func leftUnPad(d []byte) ([]byte, error) {
	if d[0] != 0 {
		return nil, rsa.ErrDecryption
	}
	if d[1] != 0 && d[1] != 1 {
		return nil, rsa.ErrDecryption
	}
	var i = 2
	for ; i < len(d); i++ {
		if d[i] == 0 {
			break
		}
	}
	i++
	if i == len(d) {
		return nil, nil
	}
	return d[i:], nil
}
