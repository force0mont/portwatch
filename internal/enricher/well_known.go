package enricher

// wellKnown maps "proto/port" strings to IANA service names for the most
// commonly encountered ports. Extend this table as needed.
var wellKnown = map[string]string{
	"tcp/21":   "ftp",
	"tcp/22":   "ssh",
	"tcp/23":   "telnet",
	"tcp/25":   "smtp",
	"tcp/53":   "domain",
	"udp/53":   "domain",
	"tcp/80":   "http",
	"tcp/110":  "pop3",
	"tcp/143":  "imap",
	"tcp/443":  "https",
	"tcp/465":  "smtps",
	"tcp/587":  "submission",
	"tcp/993":  "imaps",
	"tcp/995":  "pop3s",
	"tcp/3306": "mysql",
	"tcp/5432": "postgresql",
	"tcp/6379": "redis",
	"tcp/8080": "http-alt",
	"tcp/8443": "https-alt",
	"tcp/27017": "mongodb",
}
