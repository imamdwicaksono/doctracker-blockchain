package adapter

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	Eth     *ethclient.Client
	Auth    *bind.TransactOpts
	ChainID *big.Int
}

func NewClient(
	rpcURL string,
	privateKey *ecdsa.PrivateKey,
) (*Client, error) {

	eth, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	chainID := big.NewInt(2026)

	auth, err := bind.NewKeyedTransactorWithChainID(
		privateKey,
		chainID,
	)
	if err != nil {
		return nil, err
	}

	auth.GasPrice = big.NewInt(0)

	return &Client{
		Eth:     eth,
		Auth:    auth,
		ChainID: chainID,
	}, nil
}
