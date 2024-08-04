package types

import (
	"crypto/sha256"
	"github.com/IvanTarjan/P2PgRPC/crypto"
	"github.com/IvanTarjan/P2PgRPC/proto"
	pb "google.golang.org/protobuf/proto"
)

// HashBlock returns a SHA256 of the header.
func HashBlock(block *proto.Block) []byte {
	return HashHeader(block.Header)
}

func HashHeader(header *proto.Header) []byte {
	b, err := pb.Marshal(header)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}

func SignBlock(pk *crypto.PrivateKey, block *proto.Block) *crypto.Signature{
	signature := pk.Sign(HashBlock(block))
	return &signature
}