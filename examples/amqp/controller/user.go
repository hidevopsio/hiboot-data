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
	"hidevops.io/hiboot-data/examples/amqp/service"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/at"
	"hidevops.io/hiboot/pkg/model"
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
func (c *UserController) Publish(at struct{at.PostMapping `value:"/publish"`}) (model.Response, error) {
	err := c.userService.PublishFanout()
	return nil, err
}

func (c *UserController) Push(at struct{at.PostMapping `value:"/push"`}) (model.Response, error) {
	err := c.userService.Publish()
	return nil, err
}

func (c *UserController) Receive(at struct{at.PostMapping `value:"/receive"`}) (model.Response, error) {
	c.userService.ReceiveFanout()
	return nil, nil
}

func (c *UserController) Receive1(at struct{at.PostMapping `value:"/receive1"`}) (model.Response, error) {
	c.userService.ReceiveFanout3()
	return nil, nil
}

func (c *UserController) Create(at struct{at.PostMapping `value:"/create"`}) (model.Response, error) {
	err := c.userService.Create()
	return nil, err
}

func (c *UserController) Create1(at struct{at.PostMapping `value:"/create1"`}) (model.Response, error) {
	err := c.userService.Create1()
	return nil, err
}