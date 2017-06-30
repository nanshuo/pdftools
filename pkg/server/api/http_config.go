package api

import "time"

type HttpConfig struct {
	StaticPath string
	StaticSystemPath string
	Addr       string
	Timeout    time.Duration
	Log        string
}
