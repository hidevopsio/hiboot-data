package mongo

import (
	"hidevops.io/hiboot/pkg/log"
	"testing"
)

func TestClient_Connect(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	conf := newConfiguration(&Properties{
		Host:     "127.0.0.1",
		Port:     1111,
	})

	conf.Client()
}
