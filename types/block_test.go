package types

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/IvanTarjan/P2PgRPC/crypto"
	"github.com/IvanTarjan/P2PgRPC/util"
	"github.com/stretchr/testify/assert"
)

func TestSignBlock(t *testing.T) {
	block := util.RandomBlock()
	privKey := crypto.GeneratePrivateKey()
	pubKey := privKey.Public()
	signature := SignBlock(&privKey, block)
	assert.Equal(t, len(signature.Bytes()), 64)
	assert.True(t, signature.Verify(pubKey, HashBlock(block)))
}

func TestHashBlock(t *testing.T) {
	block := util.RandomBlock()
	hash := HashBlock(block)
	fmt.Println(hex.EncodeToString(hash))
	assert.Equal(t, len(hash), 32)
}