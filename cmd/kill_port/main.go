package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	// Set up proper flag usage
	portPtr := flag.Int("port", 0, "Port number to kill process on (required)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s --port <portnumber>\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Validate port
	if *portPtr == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Port number validation
	if *portPtr < 1 || *portPtr > 65535 {
		fmt.Fprintf(os.Stderr, "Error: Port number must be between 1 and 65535\n")
		os.Exit(1)
	}

	// Kill process
	err := killProcessOnPort(*portPtr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error killing process on port %d: %v\n", *portPtr, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully killed process on port %d\n", *portPtr)
}

func killProcessOnPort(port int) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Find PID using netstat and kill with taskkill
		findCmd := exec.Command("netstat", "-ano")
		findOut, err := findCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to find process: %v", err)
		}

		lines := strings.Split(string(findOut), "\n")
		for _, line := range lines {
			if strings.Contains(line, fmt.Sprintf(":%d", port)) {
				parts := strings.Fields(line)
				if len(parts) > 4 {
					pid := parts[len(parts)-1]
					cmd = exec.Command("taskkill", "/F", "/PID", pid)
					break
				}
			}
		}
	case "darwin", "linux", "freebsd", "openbsd":
		// Find PID using lsof and kill with kill
		findCmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port), "-t")
		findOut, err := findCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to find process: %v", err)
		}

		pid := strings.TrimSpace(string(findOut))
		if pid != "" {
			cmd = exec.Command("kill", "-9", pid)
		}
	default:
		return fmt.Errorf("unsupported platform")
	}

	if cmd == nil {
		return fmt.Errorf("no process found on port %d", port)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
