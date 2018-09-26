package fake

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type Repository struct {
	mock.Mock
}

func (r *Repository) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	args := r.Called(nil, key)
	return args[0].(*clientv3.PutResponse), args.Error(1)
}

func (r *Repository) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	args := r.Called(nil, key)
	return args[0].(*clientv3.GetResponse), args.Error(1)
}

func (r *Repository) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	args := r.Called(nil, key)
	return args[0].(*clientv3.DeleteResponse), args.Error(1)
}

func (r *Repository) Compact(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error) {
	return nil, nil
}

func (r *Repository) Do(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error) {
	return clientv3.OpResponse{}, nil
}

func (r *Repository) Txn(ctx context.Context) clientv3.Txn {
	return nil
}
