package es

import (
	"fmt"
	"github.com/olivere/elastic/v6"
	"hidevops.io/hiboot/pkg/at"
	"time"
)

type Properties struct {
	at.ConfigurationProperties `value:"es"`
	Port                       int    `json:"port" default:"5672"`
	Host                       string `json:"host" default:"127.0.0.1"`
	Username                   string `json:"username" default:"elastic"`
	Password                   string `json:"password" default:"password"`
}

type Client struct {
	*elastic.Client
}

func newClient() (client *Client) {
	return new(Client)
}

func (c *Client) Connect(p *Properties) (err error) {
	esUrl := fmt.Sprintf("http://%s:%d", p.Host, p.Port)
	client, err := elastic.NewSimpleClient(
		elastic.SetURL(esUrl),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetURL(esUrl),
		elastic.SetBasicAuth(p.Username, p.Password),
		)
	if err != nil {
		return
	}
	c.Client = client
	return
}
