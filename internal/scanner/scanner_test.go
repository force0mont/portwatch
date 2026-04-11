package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeProcNet creates a temporary /proc/net directory with fake tcp/udp files.
func fakeProcNet(t *testing.T, tcpContent, udpContent string) string {
	t.Helper()
	dir := t.TempDir()
	netDir := filepath.Join(dir, "net")
	require.NoError(t, os.MkdirAll(netDir, 0755))

	require.NoError(t, os.WriteFile(filepath.Join(netDir, "tcp"), []byte(tcpContent), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(netDir, "udp"), []byte(udpContent), 0644))
	return dir
}

const tcpHeader = "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n"
const udpHeader = "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode\n"

func TestScan_TCP_ListeningPort(t *testing.T) {
	// 0.0.0.0:22 in hex little-endian = 00000000:0016, state 0A = LISTEN
	tcpContent := tcpHeader + "   0: 00000000:0016 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1\n"
	udpContent := udpHeader

	procPath := fakeProcNet(t, tcpContent, udpContent)
	s := NewWithProcPath(procPath)

	ports, err := s.Scan()
	require.NoError(t, err)
	require.Len(t, ports, 1)
	assert.Equal(t, 22, ports[0].Port)
	assert.Equal(t, TCP, ports[0].Protocol)
	assert.Equal(t, "0.0.0.0", ports[0].Address)
}

func TestScan_SkipsNonListeningTCP(t *testing.T) {
	// state 01 = ESTABLISHED, should be skipped
	tcpContent := tcpHeader + "   0: 00000000:0016 00000000:0000 01 00000000:00000000 00:00000000 00000000     0        0 12345 1\n"
	udpContent := udpHeader

	procPath := fakeProcNet(t, tcpContent, udpContent)
	s := NewWithProcPath(procPath)

	ports, err := s.Scan()
	require.NoError(t, err)
	assert.Empty(t, ports)
}

func TestScan_UDP_Port(t *testing.T) {
	tcpContent := tcpHeader
	// 0.0.0.0:53 in hex = 00000000:0035
	udpContent := udpHeader + "   0: 00000000:0035 00000000:0000 07 00000000:00000000 00:00000000 00000000   101        0 67890 2\n"

	procPath := fakeProcNet(t, tcpContent, udpContent)
	s := NewWithProcPath(procPath)

	ports, err := s.Scan()
	require.NoError(t, err)
	require.Len(t, ports, 1)
	assert.Equal(t, 53, ports[0].Port)
	assert.Equal(t, UDP, ports[0].Protocol)
}

func TestScan_LoopbackAddress(t *testing.T) {
	// 127.0.0.1:8080 in hex little-endian = 0100007F:1F90, state 0A = LISTEN
	tcpContent := tcpHeader + "   0: 0100007F:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 11111 1\n"
	udpContent := udpHeader

	procPath := fakeProcNet(t, tcpContent, udpContent)
	s := NewWithProcPath(procPath)

	ports, err := s.Scan()
	require.NoError(t, err)
	require.Len(t, ports, 1)
	assert.Equal(t, 8080, ports[0].Port)
	assert.Equal(t, TCP, ports[0].Protocol)
	assert.Equal(t, "127.0.0.1", ports[0].Address)
}

func TestParseHexAddr(t *testing.T) {
	tests := []struct {
		hex     string
		wantIP  string
		wantPort int
		wantErr bool
	}{
		{"00000000:0016", "0.0.0.0", 22, false},
		{"0100007F:1F90", "127.0.0.1", 8080, false},
		{"invalid", "", 0, true},
	}
	for _, tt := range tests {
		ip, port, err := parseHexAddr(tt.hex)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tt.wantIP, ip)
			assert.Equal(t, tt.wantPort, port)
		}
	}
}
