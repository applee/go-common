package crypto

import (
	"crypto/rand"
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	key := make([]byte, 16)
	rand.Read(key)
	str := "just a test"
	encrypted, err := AesEncrypt([]byte(str), key)
	if err != nil {
		t.Fatal(err)
	}
	b, err := AesDecrypt(encrypted, key)
	if string(b) != str {
		t.Fail()
	}
}
