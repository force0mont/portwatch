// Package scanner provides functionality to detect open TCP/UDP ports on the host.
package scanner

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// Protocol represents a network protocol.
type Protocol string

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
)

// PortInfo holds information about a single open port.
type PortInfo struct {
	Port     int
	Protocol Protocol
	Address  string
	PID      int
}

// Scanner reads open ports from the system.
type Scanner struct {
	procPath string
}

// New creates a Scanner with the default /proc path.
func New() *Scanner {
	return &Scanner{procPath: "/proc"}
}

// NewWithProcPath creates a Scanner with a custom /proc path (useful for testing).
func NewWithProcPath(path string) *Scanner {
	return &Scanner{procPath: path}
}

// Scan returns all currently open ports on the system.
func (s *Scanner) Scan() ([]PortInfo, error) {
	var ports []PortInfo

	tcpPorts, err := s.parseProcNet("tcp", TCP)
	if err != nil {
		return nil, fmt.Errorf("scanning tcp ports: %w", err)
	}
	ports = append(ports, tcpPorts...)

	udpPorts, err := s.parseProcNet("udp", UDP)
	if err != nil {
		return nil, fmt.Errorf("scanning udp ports: %w", err)
	}
	ports = append(ports, udpPorts...)

	return ports, nil
}

func (s *Scanner) parseProcNet(proto string, protocol Protocol) ([]PortInfo, error) {
	path := fmt.Sprintf("%s/net/%s", s.procPath, proto)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ports []PortInfo
	scanner := bufio.NewScanner(f)
	scanner.Scan() // skip header line

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		// state 0A = LISTEN for TCP, any state for UDP
		state := fields[3]
		if protocol == TCP && state != "0A" {
			continue
		}
		ip, port, err := parseHexAddr(fields[1])
		if err != nil {
			continue
		}
		ports = append(ports, PortInfo{
			Port:     port,
			Protocol: protocol,
			Address:  ip,
		})
	}
	return ports, scanner.Err()
}

func parseHexAddr(hexAddr string) (string, int, error) {
	parts := strings.Split(hexAddr, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid address: %s", hexAddr)
	}
	ipHex := parts[0]
	portHex := parts[1]

	ipInt, err := strconv.ParseUint(ipHex, 16, 32)
	if err != nil {
		return "", 0, err
	}
	ipBytes := []byte{
		byte(ipInt & 0xFF),
		byte((ipInt >> 8) & 0xFF),
		byte((ipInt >> 16) & 0xFF),
		byte((ipInt >> 24) & 0xFF),
	}
	ip := net.IP(ipBytes).String()

	port, err := strconv.ParseInt(portHex, 16, 32)
	if err != nil {
		return "", 0, err
	}
	return ip, int(port), nil
}
