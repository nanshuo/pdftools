package server

import "github.com/jiusanzhou/pdf2html/pkg/server/backend"

type Config struct {
	BackendUrl  string
	RedisPrefix string

	HttpAddr string

	StaticDir string

	TmpDir string

	backend backend.Backend
}
