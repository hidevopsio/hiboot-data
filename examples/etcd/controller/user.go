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
	"hidevops.io/hiboot-data/examples/etcd/entity"
	"hidevops.io/hiboot-data/examples/etcd/service"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	"hidevops.io/hiboot/pkg/utils/copier"
	"hidevops.io/hiboot/pkg/utils/idgen"
	"net/http"
)

type userRequest struct {
	at.RequestBody

	Id       string `json:"id"`
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Age      uint   `json:"age" validate:"gte=0,lte=130"`
	Gender   uint   `json:"gender" validate:"gte=0,lte=2"`
}

// RestController
type userController struct {
	at.RestController

	userService service.UserService
}

func init() {
	app.Register(newUserController)
}

// Init inject userService automatically
func newUserController(userService service.UserService) *userController {
	return &userController{
		userService: userService,
	}
}

// Post POST /user
func (c *userController) Post(request *userRequest) (response model.Response, err error) {
	var user entity.User
	response = new(model.BaseResponse)
	copier.Copy(&user, request)

	id, err := idgen.NextString()
	if request.Id != "" {
		id = request.Id
	}

	if err == nil {
		err = c.userService.AddUser(id, &user)
		response.SetData(user)
	}
	return response, err
}

// GetById GET /id/{id}
func (c *userController) GetById(id string) (response model.Response, err error) {
	user, err := c.userService.GetUser(id)
	response = new(model.BaseResponse)
	if err != nil {
		response.SetCode(http.StatusNotFound)
	} else {
		response.SetData(user)
	}
	return
}

// DeleteById DELETE /id/{id}
func (c *userController) DeleteById(id string) (response model.Response, err error) {
	err = c.userService.DeleteUser(id)
	response = new(model.BaseResponse)
	return
}
