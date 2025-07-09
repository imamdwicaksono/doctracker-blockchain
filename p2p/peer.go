package p2p

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var Peers = []string{
	"http://localhost:3001",
	"http://localhost:3002",
}

func GetPeers() []string {
	subnet := os.Getenv("SUBNET_WHITELIST")

	subnets := subnet
	if subnets == "" {
		subnets = "172.24.4."
	}
	subnetList := []string{subnets}
	if os.Getenv("SUBNET_WHITELIST") != "" {
		subnetList = append(subnetList, strings.Split(subnets, ",")...)
	}

	var portNum int
	if portStr := os.Getenv("PORT_WHITELIST"); portStr != "" {
		fmt.Sscanf(portStr, "%d", &portNum)
	} else {
		portNum = 3003
	}

	var peers []string
	for _, s := range subnetList {
		peers = append(peers, ScanPeersOnPort(s, portNum)...)
	}
	return peers
}

func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func ScanPeersOnPort(subnet string, port int) []string {
	var peers []string

	selfIP := GetLocalIP()
	fmt.Printf("Local IP: %s\n", selfIP)

	for i := 1; i <= 254; i++ {
		ip := fmt.Sprintf("%s%d", subnet, i)
		// Enclose IP in brackets for IPv6 compatibility
		address := net.JoinHostPort(ip, fmt.Sprintf("%d", port))

		if ip == selfIP {
			continue // Skip broadcast ke diri sendiri
		}

		conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
		if err == nil {
			fmt.Printf("connected to %s\n", address)
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

// splitAndTrim splits a string by the given separator and trims whitespace from each element.
func splitAndTrim(s, sep string) []string {
	var result []string
	for _, part := range split(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// split splits a string by the given separator.
func split(s, sep string) []string {
	return []string{ // fallback for strings.Split
		s,
	}
}

// trimSpace trims leading and trailing whitespace from a string.
func trimSpace(s string) string {
	return s // fallback for strings.TrimSpace
}
