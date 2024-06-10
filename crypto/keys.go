package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io"
)

const (
	privKeyLen = 64
	pubKeyLen  = 32
	seedLen  = 32
	addressLen = 20
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

func GeneratePrivateKeyFromString(s string) PrivateKey{
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	if len(b) != seedLen {
		panic("invalid private key length")
	}

	return GeneratePrivateKeyFromSeed(b)
}

func GeneratePrivateKeyFromSeed(seed []byte) PrivateKey {
	if len(seed) != seedLen {
		panic("invalid seed length")
	}
	return PrivateKey{key: ed25519.NewKeyFromSeed(seed)}
}

func GeneratePrivateKey() PrivateKey {
	seed := make([]byte, seedLen)
	_, err := io.ReadFull(rand.Reader,seed)
	if err != nil {
		panic(err)
		}
	return PrivateKey{key: ed25519.NewKeyFromSeed(seed)}
}

func (p PrivateKey) Bytes() []byte {
	return p.key
}

func (p PrivateKey) String() string{
	return hex.EncodeToString(p.key)
}

func (p *PrivateKey) Sign(message []byte) Signature {
	return Signature{ed25519.Sign(p.key, message)}
}

func (p *PrivateKey) Public() PublicKey {
	b := make([]byte, pubKeyLen)
	copy(b, p.key[32:])
	return PublicKey{key: b}
}


type PublicKey struct {
	key ed25519.PublicKey
}

func (p *PublicKey) Bytes() []byte {
	return p.key
}

func (p *PublicKey) Address() Address{
	return Address{p.key[len(p.key)-addressLen:]}
}

type Signature struct {
	value []byte
}

func (s *Signature) Bytes() []byte {
	return s.value
}

func (s *Signature) Verify(pubKey PublicKey, message []byte) bool {
	return ed25519.Verify(pubKey.key, message, s.value)
}

type Address struct {
	value []byte
}

func (a Address) String() string {
	return hex.EncodeToString(a.value)
}

func (a Address) Bytes() []byte {
	return a.value
}