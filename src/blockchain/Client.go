package blockchain

import (
	"encoding/json"
)

type Client struct {
}

func (cli *Client) CreateBlockChain(address string) {
	bc := CreateBlockChain(address)
	defer bc.Close()

	UTXOSet := UTXOSet{bc}
	UTXOSet.ReIndex()
}

func (cli *Client) PrintChain() string {
	bc := GetBlockChain()
	defer bc.Close()
	info := bc.InfoMap()
	b, _ := json.Marshal(info)
	return string(b)
}

func (cli *Client) GetBalance(address string) int {
	bc := GetBlockChain()
	defer bc.Close()
	UTXOSet := UTXOSet{bc}
	balance := 0
	UTXOs := UTXOSet.FindUTXO(GetPubKeyHash(address))
	for _, out := range UTXOs {
		balance += out.Value
	}
	return balance
}

func (cli *Client) Send(from, to string, amount int) {
	bc := GetBlockChain()
	defer bc.Close()
	UTXOSet := UTXOSet{bc}
	tx := NewUTXOTransaction(from, to, amount, bc)
	newBlock := bc.AddBlock([]*Transaction{tx})
	UTXOSet.Update(newBlock)
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
