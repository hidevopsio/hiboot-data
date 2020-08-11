package mongo

import (
	"hidevops.io/hiboot/pkg/app"
)

const Profile = "mongo"

type configuration struct {
	app.Configuration
	// the properties member name must be amqp if the mapstructure is amqp,
	// so that the reference can be parsed
	Properties *Properties
}

func newConfiguration(properties *Properties) *configuration {
	return &configuration{Properties: properties}
}

func init() {
	app.Register(newConfiguration, new(Properties))
}

// Repository method name must be unique
func (c *configuration) Client() *Client {
	client := newClient()
	err := client.Connect(c.Properties)
	if err != nil {
		return nil
	}
	return client
}
