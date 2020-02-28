package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/log"
	"time"
)

type Properties struct {
	at.ConfigurationProperties `value:"mongo"`
	Port                       int    `json:"port" default:"5672"`
	Host                       string `json:"host" default:"127.0.0.1"`
	Username                   string `json:"username" default:""`
	Password                   string `json:"password" default:"password"`
	Timeout                    string `json:"timeout" default:"5s"`
}

type Client struct {
	*mongo.Client
}

func newClient() (client *Client) {
	return new(Client)
}

func (c *Client) Connect(p *Properties) (err error) {
	duration, err := time.ParseDuration(p.Timeout)
	if err != nil {
		log.Errorf("dataSource parse duration failed: %v", err)
		return err
	}
	mongoUrl := ""
	if p.Username == "" {
		mongoUrl = fmt.Sprintf("mongodb://%s:%d", p.Host, p.Port)
	} else {
		mongoUrl = fmt.Sprintf("mongodb://%s:%s@%s:%d", p.Username, p.Password, p.Host, p.Port)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoUrl))
	if err != nil {
		log.Errorf("mongo connection error host:%v, port:%v, username :%v, password :%v", p.Host, p.Port, p.Username, p.Password)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), duration)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Errorf("mongo client error:%v", err)
		return
	}
	c.Client = client
	log.Infof("mongo connection success host:%v, port:%v, username :%v, password :%v", p.Host, p.Port, p.Username, p.Password)
	return
}
