package adapter

import (
	"github.com/ethereum/go-ethereum/common"
)

type Listener struct {
	Client          *Client
	ContractAddress common.Address
	Contract        *DocumentAudit
}
