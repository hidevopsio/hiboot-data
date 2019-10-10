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
	"hidevops.io/hiboot-data/examples/bolt/entity"
	"hidevops.io/hiboot-data/examples/bolt/service"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
	"net/http"
)

//hi: RestController
type UserController struct {
	at.RestController
	at.RequestMapping `value:"/user"`

	userService *service.UserService
}

func init() {
	app.Register(newUserController)
}

// newUserController inject userService automatically
func newUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// Post /user
func (c *UserController) Post(at struct{at.PostMapping `value:"/"`}, user *entity.User) (model.Response, error) {
	err := c.userService.AddUser(user)
	response := new(model.BaseResponse)
	response.SetData(user)
	return response, err
}

// Get /user/id/{id}
func (c *UserController) GetById(at struct{at.GetMapping `value:"/{id}"`}, id string) (model.Response, error) {
	user, err := c.userService.GetUser(id)
	response := new(model.BaseResponse)
	if err != nil {
		response.SetCode(http.StatusNotFound)
	} else {
		response.SetData(user)
	}
	return response, err
}

// Delete /user/id/{id}
func (c *UserController) DeleteById(at struct{at.DeleteMapping `value:"/{id}"`}, id string) (response model.Response, err error) {
	err = c.userService.DeleteUser(id)
	response = new(model.BaseResponse)
	return
}
