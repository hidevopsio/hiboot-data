package amqp

import (
	"github.com/magiconair/properties/assert"
	"hidevops.io/hiboot/pkg/log"
	"testing"
)

func TestAmqp(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	conf := newConfiguration()
	conf.Properties = &properties{
		Username: "user",
		Password: "password",
		Host:     "127.0.0.1",
		Port:     1111,
	}
	ch := conf.Channel()
	var c *Channel
	assert.Equal(t, c, ch)

}
