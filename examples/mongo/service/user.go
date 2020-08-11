// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licensmongo/LICENSE-2.0
//
// Unlmongos required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either exprmongos or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"hidevops.io/hiboot-data/examples/mongo/entity"
	"hidevops.io/hiboot-data/starter/mongo"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/utils/idgen"
	"log"
	"time"
)

type UserService interface {
	AddUser(user *entity.User) (newUser *entity.User, err error)
	GetUser(id string) (user *entity.User, err error)
	GetAll() (user *[]entity.User, err error)
	DeleteUser(id string) (err error)
}

type userServiceImpl struct {
	client *mongo.Client
}

func init() {
	// register UserServiceImpl
	app.Register(newUserService)
}

// will inject gorm.Repository that configured in hidevops.io/hiboot-data/starter/gorm
func newUserService(client *mongo.Client) UserService {
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
	db := s.client.Database("test").Collection("numbers")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := db.InsertOne(ctx, user)
	fmt.Sprintln("insert res :", res.InsertedID)
	newUser = user
	return
}

func (s *userServiceImpl) GetUser(id string) (user *entity.User, err error) {
	db := s.client.Database("test").Collection("numbers")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.D{{"id", id}}
	res := db.FindOne(ctx, filter)
	user = &entity.User{}
	err = res.Decode(user)
	if err != nil {
		return nil, err
	}
	return
}

func (s *userServiceImpl) GetAll() (users *[]entity.User, err error) {
	var users1 []entity.User
	user := entity.User{}
	db := s.client.Database("test").Collection("numbers")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.D{}
	res, err := db.Find(ctx, filter)
	for res.Next(context.TODO()) {
		err := res.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		users1 = append(users1, user)
	}
	users = &users1
	return
}

func (s *userServiceImpl) DeleteUser(id string) (err error) {
	db := s.client.Database("test").Collection("numbers")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.D{{"id", id}}
	_, err = db.DeleteOne(ctx, filter)
	return
}
