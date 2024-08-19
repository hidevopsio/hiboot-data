package amqp

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestAmqp(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	conf := newConfiguration(&Properties{
		Username: "user",
		Password: "password",
		Host:     "127.0.0.1",
		Port:     1111,
	})

	ch := conf.Channel()
	var c *Channel
	assert.Equal(t, c, ch)

}
