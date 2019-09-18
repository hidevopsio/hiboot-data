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
	"hidevops.io/hiboot/pkg/app/web"
	"hidevops.io/hiboot/pkg/model"
)

//hi: RestController
type UserController struct {
	web.Controller
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
func (c *UserController) PostPublish() (model.Response, error) {
	err := c.userService.PublishFanout()
	return nil, err
}

func (c *UserController) PostPush() (model.Response, error) {
	err := c.userService.Publish()
	return nil, err
}

func (c *UserController) PostReceive() (model.Response, error) {
	c.userService.ReceiveFanout()
	return nil, nil
}

func (c *UserController) PostReceive1() (model.Response, error) {
	c.userService.ReceiveFanout3()
	return nil, nil
}

func (c *UserController) PostCreate() (model.Response, error) {
	err := c.userService.Create()
	return nil, err
}

func (c *UserController) PostCreate1() (model.Response, error) {
	err := c.userService.Create1()
	return nil, err
}