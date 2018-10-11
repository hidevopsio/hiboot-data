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

package etcd

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/log"
	"time"
)

type Repository interface {
	clientv3.KV
}

type Client struct {
	*clientv3.Client
}

type etcdConfiguration struct {
	app.Configuration
	// the properties member name must be Etcd if the mapstructure is etcd,
	// so that the reference can be parsed
	Properties properties `mapstructure:"etcd"`
}

func init() {
	app.AutoConfiguration(new(etcdConfiguration))
}

// EtcdClient create instance named etcdClient
func (c *etcdConfiguration) Client() (cli *Client) {
	cli = new(Client)
	var err error
	tlsInfo := transport.TLSInfo{
		CertFile:      c.Properties.Cert.CertFile,
		KeyFile:       c.Properties.Cert.KeyFile,
		TrustedCAFile: c.Properties.Cert.TrustedCAFile,
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		log.Error(err)
		return nil
	}
	cli.Client, err = clientv3.New(clientv3.Config{
		Endpoints:   c.Properties.Endpoints,
		DialTimeout: time.Duration(c.Properties.DialTimeout) * time.Second,
		TLS:         tlsConfig,
	})
	if err != nil {
		log.Error(err)
		return nil
	}
	return
}

// EtcdRepository create instance named etcdRepository
func (c *etcdConfiguration) Repository(cli *Client) Repository {
	if cli == nil {
		return nil
	}
	return cli.KV
}
