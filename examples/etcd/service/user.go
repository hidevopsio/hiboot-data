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

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/hidevopsio/hiboot-data/examples/etcd/entity"
	"github.com/hidevopsio/hiboot-data/starter/etcd"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/log"
	"time"
)

// UserService is the interface for userService
type UserService interface {
	AddUser(id string, user *entity.User) (err error)
	GetUser(id string) (user *entity.User, err error)
	DeleteUser(id string) (err error)
}

type userService struct {
	repository etcd.Repository
	watcher    etcd.Watcher
}

func init() {
	// register UserServiceImpl
	app.Component(newUserService)
}

// will inject BoltRepository that configured in github.com/hidevopsio/hiboot/pkg/starter/data/bolt
func newUserService(repository etcd.Repository, watcher etcd.Watcher) UserService {
	svc := &userService{
		repository: repository,
		watcher:    watcher,
	}
	svc.Watch("/user/")
	return svc
}

func (s *userService) AddUser(id string, user *entity.User) (err error) {
	if user == nil {
		return errors.New("user is not allowed nil")
	}
	userBuf, _ := json.Marshal(user)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := s.repository.Put(ctx, id, string(userBuf))
	if err != nil {
		fmt.Println("failed to put data to etcd, err:", err)
		return err
	}

	log.Debug(res)

	return nil
}

func (s *userService) GetUser(id string) (user *entity.User, err error) {
	user = &entity.User{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := s.repository.Get(ctx, id)
	if err != nil {
		log.Debugf("failed to get data from etcd, err: %v", err)
		return nil, err
	}

	if resp.Count == 0 {
		return nil, errors.New("record not found")
	}

	if err = json.Unmarshal(resp.Kvs[0].Value, &user); err != nil {
		log.Debugf("failed to unmarshal data, err: %v", err)
		return nil, err
	}

	return
}

func (s *userService) DeleteUser(id string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = s.repository.Delete(ctx, id)
	return
}

func (s *userService) Watch(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 360*time.Second)
	wch := s.watcher.Watch(ctx, key, clientv3.WithPrefix())

	go func() {
		for resp := range wch {
			for _, ev := range resp.Events {
				log.Debugf("WATCH %s %q : %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
		cancel()
	}()
}
