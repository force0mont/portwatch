// Package scanner implements host port scanning by reading from the Linux
// /proc/net filesystem. It supports TCP and UDP protocols and exposes
// a simple Scan() method that returns a slice of PortInfo structs
// describing each open/listening port found on the system.
//
// Usage:
//
//	s := scanner.New()
//	ports, err := s.Scan()
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, p := range ports {
//		fmt.Printf("%s %s:%d\n", p.Protocol, p.Address, p.Port)
//	}
//
// The scanner reads /proc/net/tcp and /proc/net/udp directly, making it
// dependency-free and suitable for use in a lightweight daemon context.
package scanner
