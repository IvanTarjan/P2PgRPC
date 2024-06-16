package types

import (
	"testing"

	"github.com/IvanTarjan/P2PgRPC/crypto"
	"github.com/IvanTarjan/P2PgRPC/proto"
	"github.com/IvanTarjan/P2PgRPC/util"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	senderPrivKey := crypto.GeneratePrivateKey()
	receiverPrivKey := crypto.GeneratePrivateKey()
	receiverAddress := receiverPrivKey.Public().Address()
	input := &proto.TxInput{
		PrevTxHash: util.RandomHash(),
		PrevOutIndex: 0,
		PublicKey: senderPrivKey.Public().Bytes(),
	}

	output1 := &proto.TxOutput{
		Amount: 5,
		Address: receiverAddress.Bytes(),
	}
	output2 := &proto.TxOutput{
		Amount: 95,
		Address: senderPrivKey.Public().Address().Bytes(),
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs: []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output1, output2},
	}

	signature := SignTransaction(&senderPrivKey, tx)
	input.Signature = signature.Bytes()

	assert.True(t, VerifyTransaction(tx))
}