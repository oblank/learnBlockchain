package blockchain

import (
	"log"
	"encoding/json"
)

type Client struct {
}

func (cli *Client) CreateBlockChain(address string) {
	bc := CreateBlockChain(address)
	bc.Close()
}

func (cli *Client) PrintChain() string {
	bc := GetBlockChain()
	defer bc.Close()
	info := bc.InfoMap()
	b, _ := json.Marshal(info)
	return string(b)
}

func (cli *Client) GetBalance(address string) int {
	if !ValidateAddress(address) {
		log.Panic("Address is not valid")
	}
	bc := GetBlockChain()
	defer bc.Close()
	balance := 0
	UTXOs := bc.FindUTXO(GetPubKeyHash(address))
	for _, out := range UTXOs {
		balance += out.Value
	}
	return balance
}

func (cli *Client) Send(from, to string, amount int) {
	if !ValidateAddress(from) {
		log.Panic("sender address is not valid")
	}
	if !ValidateAddress(to) {
		log.Panic("receiver address is not valid")
	}

	bc := GetBlockChain()
	defer bc.Close()
	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.AddBlock([]*Transaction{tx})
}

func (cli *Client) CreateWallet() string {
	wallets, _ := GetWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()
	return address
}

func (cli *Client) ListAddresses() []string {
	wallets, _ := GetWallets()
	return wallets.GetAddresses()
}
