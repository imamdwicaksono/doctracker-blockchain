package p2p

var Peers = []string{
	"http://localhost:3001",
	"http://localhost:3002",
}

func GetPeers() []string {
	return Peers
}
