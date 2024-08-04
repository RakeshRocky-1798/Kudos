package cron

import (
	"os"
)

func LocalHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}
