package es

import (
	"fmt"
	"github.com/olivere/elastic/v6"
	"hidevops.io/hiboot/pkg/at"
)

type Properties struct {
	at.ConfigurationProperties `value:"es"`
	Port                       int    `json:"port" default:"5672"`
	Host                       string `json:"host" default:"127.0.0.1"`
}

type Client struct {
	*elastic.Client
}

func newClient() (client *Client) {
	return new(Client)
}

func (c *Client) Connect(p *Properties) (err error) {
	esUrl := fmt.Sprintf("http://%s:%d", p.Host, p.Port)
	client, err := elastic.NewSimpleClient(elastic.SetURL(esUrl))
	if err != nil {
		return
	}
	c.Client = client
	return
}
