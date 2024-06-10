package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()
	assert.Equal(t, len(privKey.Bytes()), privKeyLen)
	pubKey := privKey.Public()
	assert.Equal(t, len(pubKey.Bytes()), pubKeyLen)
}

func TestGeneratePrivateKeyFromString(t *testing.T) {
	var (
		seed = "e85e1b917a219887209c4aac27566d65749c8f0b53a71b813eb32eb0550fac91"
		expectedPrivKey = "e85e1b917a219887209c4aac27566d65749c8f0b53a71b813eb32eb0550fac91c2cbfb2aa521ca95ec5b9509adbe3b11680260dbbfdf03648c848a25ef0543aa"
	)
	privKey := GeneratePrivateKeyFromString(seed)
	assert.Equal(t, privKey.String(), expectedPrivKey)

}


func TestPrivateKeySign(t *testing.T){
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	msg := []byte("hello")
	signature := privKey.Sign(msg)

	assert.True(t, signature.Verify(pubKey, msg))

	assert.False(t, signature.Verify(pubKey, []byte("world")))

	invalidPrivateKey := GeneratePrivateKey()
	assert.False(t, signature.Verify(invalidPrivateKey.Public(), msg))
}

func TestPublicKeyToAddress(t *testing.T){
	privKeyLen := GeneratePrivateKey()
	pubKey := privKeyLen.Public()
	address := pubKey.Address()
	assert.Equal(t, len(address.Bytes()), addressLen)
	fmt.Println(address)
}