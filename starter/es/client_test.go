package es

import (
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestConnect(t *testing.T) {
	pro := &Properties{
		ConfigurationProperties: at.ConfigurationProperties{},
		Port:                    0,
		Host:                    "",
	}
	client := newClient()
	err := client.Connect(pro)
	assert.Equal(t, nil, err)
}
