package api

import (
	"testing"
)

func TestServe(t *testing.T) {

	config := &HttpConfig{
		Addr: ":8080",
		StaticSystemPath: "C:\\Users\\viruser.v-desktop\\Downloads",
		StaticPath: "/storage",
	}

	_api, _ := NewApiHttp(config)
	_api.Serve()
}