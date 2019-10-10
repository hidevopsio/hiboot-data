package etcd

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot-data/starter/etcd/fake"
	"hidevops.io/hiboot/pkg/log"
	"testing"
)

func TestEtcd(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	conf := new(etcdConfiguration)

	t.Run("should create instance named etcdClient", func(t *testing.T) {
		conf.Properties = &Properties{
			DialTimeout:    5,
			RequestTimeout: 10,
			Endpoints:      []string{"172.16.10.470:2379"},
			Cert: cert{CertFile: "config/certs/etcd.pem",
				KeyFile:       "config/certs/etcd-key.pem",
				TrustedCAFile: "config/certs/ca.pem"},
		}
		client := conf.Client()
		assert.Equal(t, (*Client)(nil), client)

	})

	client := new(Client)
	client.Client = new(clientv3.Client)
	t.Run("should not create instance named etcdRepository", func(t *testing.T) {
		repo := conf.Repository(client)
		assert.Equal(t, nil, repo)
	})

	t.Run("should create instance named etcdRepository", func(t *testing.T) {
		client.KV = new(fake.Repository)
		repo := conf.Repository(client)
		assert.Equal(t, client.KV, repo)
	})
}
