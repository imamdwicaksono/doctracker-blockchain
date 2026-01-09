// blockchain/audit/anchor.go
package audit

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type BlockPayload struct {
	Height    uint64
	BlockHash string
}

func (c *Client) Anchor(height uint64, blockHash string) error {
	tx, err := c.Contract.AnchorBlock(
		c.Auth,
		big.NewInt(int64(height)),
		common.HexToHash(blockHash),
	)
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(context.Background(), c.Eth, tx)
	return err
}
