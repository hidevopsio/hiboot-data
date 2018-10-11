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

package fake

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type Repository struct {
	mock.Mock
}

func (e *Repository) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	args := e.Called(nil, key)
	return args[0].(*clientv3.PutResponse), args.Error(1)
}

func (e *Repository) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	args := e.Called(nil, key)
	return args[0].(*clientv3.GetResponse), args.Error(1)
}

func (e *Repository) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	args := e.Called(nil, key)
	return args[0].(*clientv3.DeleteResponse), args.Error(1)
}

func (e *Repository) Compact(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error) {
	return nil, nil
}

func (e *Repository) Do(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error) {
	return clientv3.OpResponse{}, nil
}

func (e *Repository) Txn(ctx context.Context) clientv3.Txn {
	return nil
}
