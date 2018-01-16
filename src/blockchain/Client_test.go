package blockchain

import (
	"testing"
)

const foo = "1Bknnztugm7d8s813CFvUcJHhzezUUDPYM"
const bar = "1JHukhuBKcaWFkX5Mm8JS4uK3paMP9P7Gu"

func TestCreateWallet(t *testing.T) {
	cli := Client{}
	address := cli.CreateWallet()
	t.Log(address)
}

func TestListAddresses(t *testing.T) {
	cli := Client{}
	t.Log(cli.ListAddresses())
}

func TestCreateBlockChain(t *testing.T) {
	cli := Client{}
	cli.CreateBlockChain(foo)
	t.Log("Done")
}

func TestPrintChain(t *testing.T) {
	cli := Client{}
	t.Log(cli.PrintChain())
}

func TestSend(t *testing.T) {
	cli := Client{}
	cli.Send(bar, foo, 1)
	t.Log("Success")
}

func TestGetBalance(t *testing.T) {
	cli := Client{}
	balance := cli.GetBalance(bar)
	t.Log(balance)
}
