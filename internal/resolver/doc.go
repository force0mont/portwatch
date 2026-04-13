// Package resolver maps open socket inodes to the owning process by walking
// the /proc filesystem.
//
// Given an inode number obtained from /proc/net/tcp or /proc/net/udp, Lookup
// scans every numeric PID directory under /proc, reads the file-descriptor
// symlinks in /proc/<pid>/fd, and returns the first ProcessInfo whose socket
// inode matches.
//
// The resolved process name is read from /proc/<pid>/comm and falls back to
// "unknown" when the file is unavailable (e.g. the process exited between the
// inode scan and the comm read).
//
// Usage:
//
//	r := resolver.New()
//	info, err := r.Lookup(inodeNumber)
//	if err == nil {
//		fmt.Printf("port owned by pid=%d name=%s\n", info.PID, info.Name)
//	}
package resolver
