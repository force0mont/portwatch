// Package config provides loading, validation, and default handling for
// portwatch configuration files.
//
// Configuration files are JSON documents with the following structure:
//
//	{
//	  "interval_seconds": 10,
//	  "log_format": "json",
//	  "rules": [
//	    { "port": 22,  "protocol": "tcp", "action": "allow" },
//	    { "port": 80,  "protocol": "tcp", "action": "allow" },
//	    { "port": 443, "protocol": "tcp", "action": "allow" }
//	  ]
//	}
//
// Fields omitted from the file receive sensible defaults via DefaultConfig.
package config
