package p2p

import (
	"fmt"
	"net"
	"os"
	"time"
)

var Peers = []string{
	"http://localhost:3001",
	"http://localhost:3002",
}

func GetPeers() []string {
	return ScanPeersOnPort(3003)
	return Peers
}

func ScanPeersOnPort(port int) []string {
	subnet := "172.24.4." // sesuaikan dengan jaringanmu
	var peers []string

	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s%d", subnet, i)
		address := fmt.Sprintf("%s:%d", ip, port)

		conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
		if err == nil {
			peers = append(peers, address)
			conn.Close()
		}
	}
	return peers
}

func WritePeerListToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	for _, peer := range GetPeers() {
		if _, err := file.WriteString(peer + "\n"); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return nil
}
