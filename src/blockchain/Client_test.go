package blockchain

import (
	"testing"
)

const foo = "1A5m7hwXtbzYPctXNwMheR6KPov74yUSdw"
const bar = "1KpLM7HEuvS854URNYByYYXWQXrVXsYcEY"

func TestCreateBlockChain(t *testing.T) {
	cli := Client{}
	cli.CreateBlockChain(foo)
	t.Log("Done")
}

func TestPrintChain(t *testing.T) {
	cli := Client{}
	t.Log(cli.PrintChain())
}

func TestCreateWallet(t *testing.T) {
	cli := Client{}
	address := cli.CreateWallet()
	t.Log(address)
}

func TestListAddresses(t *testing.T)  {
	cli := Client{}
	t.Log(cli.ListAddresses())
}

func TestSend(t *testing.T) {
	cli := Client{}
	cli.Send(foo, bar, 5)
	t.Log("Success")
}

func TestGetBalance(t *testing.T) {
	cli := Client{}
	balance := cli.GetBalance(foo)
	t.Log(balance)
}