package helper

import "time"

var (
	BotID   int64
	BotName string
	Version = "v0.2.0 (Go Rewrite)"
	Update  = time.Now().Format("2006-01-02 15:04:05")
)
