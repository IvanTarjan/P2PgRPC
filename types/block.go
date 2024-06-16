package types

import (
	"crypto/sha256"

	"github.com/IvanTarjan/P2PgRPC/crypto"
	"github.com/IvanTarjan/P2PgRPC/proto"
	pb "google.golang.org/protobuf/proto"
)

// HashBlock returns a SHA256 of the header.
func hashBlock(block *proto.Block) []byte {
	b, err := pb.Marshal(block)
	if err != nil {
		panic(err)
	}
	hash:= sha256.Sum256(b)
	return hash[:]
}

func SignBlock(pk *crypto.PrivateKey, block *proto.Block) *crypto.Signature{
	signature := pk.Sign(hashBlock(block))
	return &signature
}