package p2p

type Message struct {
	Type string      `json:"type"` // new_block, get_chain, etc.
	Data interface{} `json:"data"`
}
