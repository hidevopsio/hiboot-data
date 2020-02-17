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

package controller

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"hidevops.io/hiboot-data/examples/mongo/entity"
	"hidevops.io/hiboot-data/examples/mongo/service/mocks"
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/log"
	"hidevops.io/hiboot/pkg/utils/idgen"
	"net/http"
	"testing"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestCrdRequest(t *testing.T) {

	mockUserService := new(mocks.UserService)
	userController := newUserController(mockUserService)
	testApp := web.NewTestApp(userController).Run(t)

	id, err := idgen.NextString()
	assert.Equal(t, nil, err)

	testUser := &entity.User{
		Id:       id,
		Name:     "Bill Gates",
		Username: "billg",
		Password: "3948tdaD",
		Email:    "bill.gates@microsoft.com",
		Age:      60,
		Gender:   1,
	}

	// first, call mocks.UserService.AddUser
	mockUserService.On("AddUser", testUser).Return(testUser, nil)
	// then run the test that will call UserService.AddUser
	t.Run("should add user with POST request", func(t *testing.T) {
		// First, let's Post User
		testApp.Post("/user").
			WithJSON(testUser).
			Expect().Status(http.StatusOK)
	})

	mockUserService.On("GetUser", id).Return(testUser, nil)
	t.Run("should get user with GET request", func(t *testing.T) {
		// Then Get User
		// e.g. GET /user/id/123456
		testApp.Get("/user/id/{id}").
			WithPath("id", id).
			Expect().Status(http.StatusOK)
	})

	mockUserService.On("GetAll").Return(&[]entity.User{*testUser}, nil)
	t.Run("should get user with GET request", func(t *testing.T) {
		// Then Get User
		// e.g. GET /user/id/123456
		testApp.Get("/user/all").
			Expect().Status(http.StatusOK)
	})

	// assert that the expectations were met
	mockUserService.AssertExpectations(t)

	unknownId, err := idgen.NextString()
	assert.Equal(t, nil, err)
	mockUserService.On("GetUser", unknownId).Return((*entity.User)(nil), errors.New("not found"))

	t.Run("should return 404 if trying to find a record that does not exist", func(t *testing.T) {
		// Then Get User
		testApp.Get("/user/id/{id}").
			WithPath("id", unknownId).
			Expect().Status(http.StatusNotFound)
	})

	// assert that the expectations were met
	mockUserService.AssertExpectations(t)

	mockUserService.On("DeleteUser", id).Return(nil)
	t.Run("should delete the record with DELETE request", func(t *testing.T) {
		// Finally Delete User
		testApp.Delete("/user/id/{id}").
			WithPath("id", id).
			Expect().Status(http.StatusOK)
	})
}
