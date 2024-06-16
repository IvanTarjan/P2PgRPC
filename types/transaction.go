package types

import (
	"crypto/sha256"

	"github.com/IvanTarjan/P2PgRPC/crypto"
	"github.com/IvanTarjan/P2PgRPC/proto"
	pb "google.golang.org/protobuf/proto"
)

func SignTransaction(pk *crypto.PrivateKey, tx *proto.Transaction) *crypto.Signature{
	signature := pk.Sign(HashTransaction(tx))
	return &signature
}

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}

func VerifyTransaction(tx *proto.Transaction) bool {
	for _, input := range tx.Inputs{
		var (
			sig = crypto.SignatureFromBytes(input.Signature)
			pubKey = crypto.PublicKeyFromBytes(input.PublicKey)
		)
		// TODO: make sure we don't run into problems 
		//after verification cause we have set the signature to nil
		input.Signature = nil
		if !sig.Verify(pubKey, HashTransaction(tx)){
			return false
		}
	}
	return true
}