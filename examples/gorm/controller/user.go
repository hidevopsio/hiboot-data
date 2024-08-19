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
	"github.com/hidevopsio/hiboot-data/examples/gorm/entity"
	"github.com/hidevopsio/hiboot-data/examples/gorm/service"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/model"
	"net/http"
)

// RestController
type userController struct {
	at.RestController
}

func init() {
	app.Register(newUserController)
}

// newUserController inject userService automatically
func newUserController() *userController {
	return &userController{}
}

// Post create /user
func (c *userController) Post(request *entity.User, userService *service.UserService) (model.Response, error) {
	err := userService.AddUser(request)
	response := new(model.BaseResponse)
	response.SetData(request)
	return response, err
}

// GetById GET /id/{id}
func (c *userController) GetById(id uint64, userService *service.UserService) (response model.Response, err error) {
	user, err := userService.GetUser(id)
	response = new(model.BaseResponse)
	if err != nil {
		response.SetCode(http.StatusNotFound)
	} else {
		response.SetData(user)
	}
	return
}

// GetAll GET /id/{id}
func (c *userController) GetAll(userService *service.UserService) (response model.Response, err error) {
	users, err := userService.GetAll()
	response = new(model.BaseResponse)
	response.SetData(users)
	return
}

// DeleteById DELETE /id/{id}
func (c *userController) DeleteById(id uint64, userService *service.UserService) (response model.Response, err error) {
	err = userService.DeleteUser(id)
	response = new(model.BaseResponse)
	return
}
