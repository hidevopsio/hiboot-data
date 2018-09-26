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
	"encoding/json"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	_ "github.com/erikstmartin/go-testdb"
	"github.com/hidevopsio/hiboot-data/examples/etcd/entity"
	"github.com/hidevopsio/hiboot-data/starter/etcd/fake"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/idgen"
	"github.com/stretchr/testify/assert"
	"testing"
)

var fakeUser = entity.User{
	Id:       "",
	Name:     "Bill Gates",
	Username: "billg",
	Password: "3948tdaD",
	Email:    "bill.gates@microsoft.com",
	Age:      60,
	Gender:   1,
}

func newID(t *testing.T, path string) string {
	id, err := idgen.NextString()
	fakeUser.Id = id
	assert.Equal(t, nil, err)
	return path + id
}

func TestUserCrud(t *testing.T) {
	fakeRepository := new(fake.Repository)
	fakeWatcher := new(fake.Watcher)
	userService := newUserService(fakeRepository, fakeWatcher)

	id := newID(t, "/UserAddedEvent/")
	t.Run("should return error if user is nil", func(t *testing.T) {
		err := userService.AddUser(id, (*entity.User)(nil))
		assert.NotEqual(t, nil, err)
	})

	response := new(clientv3.PutResponse)

	fakeRepository.On("Put", nil, id).Return(response, nil)
	t.Run("should add user", func(t *testing.T) {
		err := userService.AddUser(id, &fakeUser)
		assert.Equal(t, nil, err)
	})

	simulationErr := errors.New("simulation err")
	id = newID(t, "/user/")
	fakeRepository.On("Put", nil, id).Return((*clientv3.PutResponse)(nil), simulationErr)
	t.Run("should add user", func(t *testing.T) {
		err := userService.AddUser(id, &fakeUser)
		assert.Equal(t, err, simulationErr)
	})

	recordNotFound := errors.New("record not found")
	id = newID(t, "/user/")
	fakeRepository.On("Get", nil, id).Return((*clientv3.GetResponse)(nil), recordNotFound)
	t.Run("should generate user id", func(t *testing.T) {
		//u := &entity.User{}
		_, err := userService.GetUser(id)
		log.Debug("Error %v", err)
		assert.NotEqual(t, recordNotFound, nil)
	})

	fakeUserBuf, _ := json.Marshal(&fakeUser)
	getRes := new(clientv3.GetResponse)
	kv := &mvccpb.KeyValue{
		Key:   []byte("test"),
		Value: fakeUserBuf,
	}
	getRes.Kvs = append(getRes.Kvs, kv)
	id = newID(t, "/user/")
	fakeRepository.On("Get", nil, id).Return(getRes, nil)
	t.Run("should generate user id", func(t *testing.T) {
		//u := &entity.User{}
		var err error
		err = nil
		_, err = userService.GetUser(id)
		log.Debug("Error %v", err)
		assert.NotEqual(t, getRes, nil)
	})

	getRes = new(clientv3.GetResponse)
	kv = &mvccpb.KeyValue{
		Key:   []byte("test"),
		Value: []byte("test"),
	}
	getRes.Kvs = append(getRes.Kvs, kv)
	id = newID(t, "/user/")
	fakeRepository.On("Get", nil, id).Return(getRes, nil)
	t.Run("should generate user id", func(t *testing.T) {
		//u := &entity.User{}
		var err error
		err = nil
		_, err = userService.GetUser(id)
		log.Debug("Error %v", err)
		assert.NotEqual(t, getRes, nil)
	})
	id = newID(t, "/user/")
	fakeRepository.On("Delete", nil, id).Return((*clientv3.DeleteResponse)(nil), nil)
	t.Run("should delete user", func(t *testing.T) {
		err := userService.DeleteUser(id)
		assert.Equal(t, nil, err)
	})
}

func TestDependencyInjection(t *testing.T) {
	testApp := web.NewTestApplication(t).(app.ApplicationContext)

	svc := testApp.GetInstance("userService")
	if svc != nil {
		t.Run("should not get record that does not exist", func(t *testing.T) {
			s := svc.(UserService)
			id, err := idgen.NextString()
			assert.Equal(t, nil, err)

			_, err = s.GetUser(id)
			assert.NotEqual(t, nil, err)
		})
	}
}
