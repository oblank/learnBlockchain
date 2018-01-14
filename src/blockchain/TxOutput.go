package blockchain

import (
	"bytes"
	"encoding/hex"
)

type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

func NewTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock(address)
	return txo
}

func (out *TxOutput) Lock(address string) {
	out.PubKeyHash = GetPubKeyHash(address)
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func (out *TxOutput) InfoMap() map[string]interface{} {
	info := make(map[string]interface{})
	info["Value"] = out.Value
	info["PubKeyHash"] = hex.EncodeToString(out.PubKeyHash)
	return info
}
