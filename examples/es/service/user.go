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
	"fmt"
	"github.com/olivere/elastic/v6"
	"hidevops.io/hiboot-data/examples/es/entity"
	"hidevops.io/hiboot-data/starter/es"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/utils/idgen"
	"log"
)

type UserService interface {
	AddUser(user *entity.User) (newUser *entity.User, err error)
	GetUser(id string) (user *entity.User, err error)
	GetAll() (user *[]entity.User, err error)
	DeleteUser(id string) (err error)
}

type userServiceImpl struct {
	client *es.Client
}

func init() {
	// register UserServiceImpl
	app.Register(newUserService)
}

// will inject gorm.Repository that configured in hidevops.io/hiboot-data/starter/gorm
func newUserService(client *es.Client) UserService {
	return &userServiceImpl{
		client: client,
	}
}

func (s *userServiceImpl) AddUser(user *entity.User) (newUser *entity.User, err error) {
	if user.Id == "" {
		id, err := idgen.NextString()
		if err != nil {
			return nil, err
		}
		user.Id = id
	}
	_, err = s.client.Index().Index("test").Type("test").Id(user.Id).BodyJson(user).Do(context.Background())
	newUser = user
	return
}

func (s *userServiceImpl) GetUser(id string) (user *entity.User, err error) {
	esResponse, err := s.client.Get().Index("test").Type("test").Id(id).Do(context.Background())
	if err != nil {
		// Handle Error
		return
	}
	err = json.Unmarshal(*esResponse.Source, &user)
	return
}

func (s *userServiceImpl) GetAll() (users *[]entity.User, err error) {
	var users1 []entity.User
	//err = s.repository.Find(users).Error()
	ctx := context.Background()
	query := elastic.NewBoolQuery()
	res, err := s.client.Search().Index("test").Type("test").Size(10).Query(query).Do(ctx)
	if err != nil {
		log.Println("err:", err)
		return
	}

	if res.TotalHits() > 0 {
		log.Printf("Found a total of %d indice\n", res.TotalHits())
		for _, hit := range res.Hits.Hits {
			var user entity.User
			err := json.Unmarshal(*hit.Source, &user)
			if err != nil {
				log.Printf("%v", err)
				return nil, err
			}
			fmt.Printf("user: %v", user)
			users1 = append(users1, user)
		}
	} else {
		log.Println("没有查询到数据")
	}
	users = &users1
	return
}

func (s *userServiceImpl) DeleteUser(id string) (err error) {
	_, err = s.client.Delete().Index("test").Type("test").Id(id).Do(context.Background())
	return
}
