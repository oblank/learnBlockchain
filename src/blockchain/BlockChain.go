package blockchain

import (
	"github.com/boltdb/bolt"
	"encoding/hex"
	"os"
	"log"
	"bytes"
	"errors"
	"crypto/ecdsa"
)

const blockBucket = "blocks"

const god = "God"

const dbFile = "hello.db"

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

func CreateBlockChain(address string) *BlockChain {
	if !ValidateAddress(address) {
		log.Panic("Address is not valid")
	}
	if dbExists() {
		log.Panic("Blockchain already exists.")
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	PanicIfError(err)
	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTx(address, god)
		genesis := NewGenesisBlock(cbtx)
		b, err := tx.CreateBucket([]byte(blockBucket))
		if err != nil {
			return err
		}
		if err = b.Put(genesis.Hash, genesis.Serialize()); err != nil {
			return err
		}
		if err = b.Put([]byte("l"), genesis.Hash); err != nil {
			return err
		}
		tip = genesis.Hash
		return nil
	})
	return &BlockChain{tip, db}
}
func GetBlockChain() *BlockChain {
	if !dbExists() {
		log.Panic("No existing blockchain found. Create one first.")
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	PanicIfError(err)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		tip = b.Get([]byte("l"))
		return nil
	})
	PanicIfError(err)
	return &BlockChain{tip, db}
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func (bc *BlockChain) FindUTXO(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)
	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

func (bc *BlockChain) AddBlock(transactions []*Transaction) {
	var lastHash []byte
	for _, tx := range transactions {
		if !bc.VerifyTransaction(tx) {
			log.Panic("invalid transaction")
		}
	}
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	PanicIfError(err)
	newBlock := NewBlock(transactions, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		PanicIfError(err)
		err = b.Put([]byte("l"), newBlock.Hash)
		PanicIfError(err)
		bc.tip = newBlock.Hash
		return nil
	})
	PanicIfError(err)
}

func (bc *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()
	for bci.HasNext() {
		block := bci.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if in.UseKey(pubKeyHash) {
						inTXID := hex.EncodeToString(in.Txid)
						spentTXOs[inTXID] = append(spentTXOs[inTXID], in.Vout)
					}
				}
			}
		}
	}
	return unspentTXs
}

func (bc *BlockChain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOutputs
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()
	for bci.HasNext() {
		block := bci.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
	}
	return Transaction{}, errors.New("transaction is not found")
}

func (bc *BlockChain) SignTransaction(tx *Transaction, priKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)
	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		PanicIfError(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(priKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	preTXs := make(map[string]Transaction)
	for _, vin := range tx.Vin {
		preTX, err := bc.FindTransaction(vin.Txid)
		PanicIfError(err)
		preTXs[hex.EncodeToString(preTX.ID)] = preTX
	}
	return tx.Verify(preTXs)
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{bc.tip, bc.db}
}

func (bc *BlockChain) InfoMap() map[string]interface{} {
	info := make(map[string]interface{})
	var blocks []map[string]interface{}
	bci := bc.Iterator()
	for bci.HasNext() {
		block := bci.Next()
		blocks = append(blocks, block.InfoMap())
	}
	info["blocks"] = blocks
	return info
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func (bc *BlockChain) Close() {
	bc.db.Close()
}

type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (i *BlockChainIterator) Next() *Block {
	var block *Block
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	PanicIfError(err)
	i.currentHash = block.PreBlockHash
	return block
}

func (i *BlockChainIterator) HasNext() bool {
	return len(i.currentHash) != 0
}
