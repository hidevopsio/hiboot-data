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
	"github.com/hidevopsio/hiboot-data/examples/gorm/entity"
	"github.com/hidevopsio/hiboot-data/starter/gorm"
	"github.com/hidevopsio/hiboot-data/starter/redis"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/utils/idgen"
	"strconv"
)

type UserService struct {
	at.Scope `value:"request"`
	db       *gorm.DB
	cache    *redis.Client
}

func init() {
	// register UserService
	app.Register(newUserService)
}

// will inject gorm.Repository that configured in hiboot-data/starter/gorm
func newUserService(db *gorm.DB, redisClient *redis.Client) *UserService {
	_ = db.AutoMigrate(&entity.User{})
	return &UserService{
		db:    db,
		cache: redisClient,
	}
}

func (s *UserService) AddUser(user *entity.User) (err error) {
	if user == nil {
		return errors.New("user is not allowed nil")
	}
	if user.Id == 0 {
		user.Id, _ = idgen.Next()
	}
	err = s.db.Create(user).Error
	if err != nil {
		return
	}
	err = s.cacheUser(user)
	return
}

func (s *UserService) cacheUser(user *entity.User) (err error) {
	var userJSON []byte
	userJSON, err = json.Marshal(user)
	if err != nil {
		return
	}
	err = s.cache.Set(context.Background(), strconv.FormatUint(user.Id, 16), string(userJSON), 0).Err()
	return
}

func (s *UserService) GetUser(id uint64) (user *entity.User, err error) {
	user = &entity.User{}
	var res string
	res, err = s.cache.Get(context.Background(), strconv.FormatUint(id, 16)).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		err = s.db.Where("id = ?", id).First(user).Error
		if err != nil {
			return
		}
		err = s.cacheUser(user)
	}

	return
}

func (s *UserService) GetAll() (users *[]entity.User, err error) {
	users = &[]entity.User{}
	err = s.db.Find(users).Error
	return
}

func (s *UserService) DeleteUser(id uint64) (err error) {
	err = s.db.Where("id = ?", id).Delete(entity.User{}).Error
	return
}
