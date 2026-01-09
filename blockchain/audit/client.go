// blockchain/audit/client.go
package audit

import (
	"errors"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	// Replace with actual path to generated bindings
	"doc-tracker/blockchain/audit/contract"
)

type Client struct {
	Eth      *ethclient.Client
	Contract *contract.Contract
	Auth     *bind.TransactOpts
}

func NewClient() (*Client, error) {
	rpc := os.Getenv("RPC_URL")
	if rpc == "" {
		return nil, errors.New("RPC_URL not set")
	}

	cli, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	privKey := os.Getenv("PRIVATE_KEY")
	if privKey == "" {
		return nil, errors.New("PRIVATE_KEY not set")
	}

	privKey = strings.TrimPrefix(privKey, "0x")
	priv, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}

	chainID, ok := new(big.Int).SetString(os.Getenv("CHAIN_ID"), 10)
	if !ok {
		return nil, errors.New("invalid CHAIN_ID")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(priv, chainID)
	if err != nil {
		return nil, err
	}

	addr := common.HexToAddress(os.Getenv("AUDIT_CONTRACT_ADDRESS"))
	ctr, err := contract.NewContract(addr, cli)
	if err != nil {
		return nil, err
	}

	return &Client{
		Eth:      cli,
		Contract: ctr,
		Auth:     auth,
	}, nil
}

func Verify(height uint64, hash string) (bool, error) {
	client, err := NewClient()
	if err != nil {
		return false, err
	}

	return client.Contract.Verify(
		nil,
		big.NewInt(int64(height)),
		common.HexToHash(hash),
	)
}

func StringToHash(s string) common.Hash {
	return common.HexToHash(s)
}
