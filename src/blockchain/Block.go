package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
	"encoding/hex"
)

type Block struct {
	Timestamp    int64
	Transactions []*Transaction
	PreBlockHash []byte
	Hash         []byte
	Nonce        int
}

func (b *Block) InfoMap() map[string]interface{} {
	info := make(map[string]interface{})
	info["Hash"] = hex.EncodeToString(b.Hash)
	info["PreBlockHash"] = hex.EncodeToString(b.PreBlockHash)

	var Transactions []map[string]interface{}
	for _, tx := range b.Transactions {
		Transactions = append(Transactions, tx.InfoMap())
	}
	info["Transactions"] = Transactions
	return info
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}
	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		panic(err)
	}
	return &block
}
