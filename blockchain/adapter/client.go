package adapter

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	Eth  *ethclient.Client
	Auth *bind.TransactOpts
}

func NewClient() (*Client, error) {
	rpc := "http://127.0.0.1:8545"

	eth, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(
		"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
	)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(
		privateKey,
		big.NewInt(31337),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		Eth:  eth,
		Auth: auth,
	}, nil
}
