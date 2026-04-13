// Package topology classifies observed network ports into structural groups
// based on protocol (tcp/udp) and address class (loopback, private, public).
//
// A Topology is rebuilt on each scan cycle via Build, which atomically
// replaces all groups. Consumers call Groups to obtain a point-in-time
// snapshot safe for independent use without holding any lock.
//
// Address classification follows RFC 5735 / RFC 4193:
//   - Loopback: 127.0.0.0/8 and ::1
//   - Private:  10/8, 172.16/12, 192.168/16, fc00::/7
//   - Public:   everything else
//
// Topology is safe for concurrent use.
package topology
