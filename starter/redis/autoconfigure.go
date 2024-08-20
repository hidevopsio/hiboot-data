// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redis

import (
	"context"
	"fmt"
	"github.com/hidevopsio/hiboot-data/utils"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/redis/go-redis/v9"
	"time"
)

const Profile = "redis"

type Client struct {
	at.Scope `value:"request"`

	*redis.Client
}

type configuration struct {
	at.AutoConfiguration

	prop   *properties
	client *Client
}

func newConfiguration(prop *properties) *configuration {
	return &configuration{prop: prop}
}

func init() {
	app.Register(newConfiguration, new(properties))
}

func (c *configuration) Client() (cli *Client, err error) {
	var redisCli *redis.Client
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(c.prop.Timeout))
	defer cancel()
	// return cached if it is healthy
	if c.client != nil {
		_, err = c.client.Ping(ctx).Result()
		if err == nil {
			// If the connection is still alive, return the existing connection
			cli = c.client
			return
		}
		log.Warnf("lost connection to redis, attempting to reconnect...")
		c.client = nil
	}

	// create new connection if it is unhealthy
	log.Infof("create a new database connection to %v:%v db: %v", c.prop.Host, c.prop.Port, c.prop.DB)
	cli = new(Client)
	password := c.prop.Password
	if c.prop.Config.Decrypt {
		password = utils.Decrypt(password, c.prop.Config.DecryptKey)
	}
	redisCli = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", c.prop.Host, c.prop.Port),
		Password: password,
		DB:       c.prop.DB,
	})
	if err != nil {
		log.Errorf("failed to connect to redis server: %v", err)
		return
	}
	_, err = redisCli.Ping(ctx).Result()
	if err == nil {
		log.Infof("redis %v:%v is connected", c.prop.Host, c.prop.Port)
		cli.Client = redisCli
		c.client = cli
	}

	return
}
