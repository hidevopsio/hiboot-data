package mongo

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/hiboot/pkg/at"
	"testing"
)

func TestConnect(t *testing.T) {
	pro := &Properties{
		ConfigurationProperties: at.ConfigurationProperties{},
		Port:                    27017,
		Host:                    "10.10.10.10",
		Timeout:                 "1s",
	}
	client := newClient()
	err := client.Connect(pro)
	assert.Equal(t, "context deadline exceeded", err.Error())
}
