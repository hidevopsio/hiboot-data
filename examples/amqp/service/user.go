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
	"hidevops.io/hiboot-data/starter/amqp"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
)

type UserService struct {
	newChannel amqp.NewChannel
}

func init() {
	app.Register(newUserService)
}

// will inject BoltRepository that configured in hidevops.io/hiboot/pkg/starter/data/bolt
func newUserService(newChannel amqp.NewChannel) *UserService {
	return &UserService{
		newChannel: newChannel,
	}
}

const (
	exchange = "test1"
)


func (s *UserService) PublishDirect() (err error) {
	chn := s.newChannel()
	defer chn.Close()
	err = chn.PublishDirect("", "test1", "hello", "info")
	return
}

func (s *UserService) PublishFanout() (err error) {
	chn := s.newChannel()
	defer chn.Close()
	err = chn.PublishFanout(exchange, "hello")
	return
}


func (s *UserService) ReceiveFanout() error {
	chn := s.newChannel()
	defer chn.Close()
	go func() {
		for {
			c, err := chn.ReceiveFanout("test2", exchange)
			log.Infof("cha: %s,  err: %v", *c, err)
			chn.Close()
		}
	}()
	return nil
}

func (s *UserService) ReceiveFanout3() error {
	chn := s.newChannel()
	go func() {
		for {
			c, err := chn.ReceiveFanout("test1", exchange)
			log.Infof("cha: %s,  err: %v", *c, err)
			chn.Close()
		}
	}()
	return nil
}