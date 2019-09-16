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
	"time"
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
	mgsConnect = "hello world"
	exchange   = "test23223"
	queueName  = "Test"
)

func (s *UserService) Create() error {
	shn := s.newChannel()
	err := shn.CreateFanout(queueName, exchange)
	return err
}

func (s *UserService) PublishDirect() error {
	shn := s.newChannel()
	err := shn.PublishDirect(exchange, queueName, mgsConnect, "info")
	return err
}

func (s *UserService) PublishFanout() error {
	shn := s.newChannel()
	err := shn.PublishFanout(exchange, "hello")
	return err
}

func (s *UserService) ReceiveFanout() error {
	go func() {
		c := 1
		shn := s.newChannel()
		defer shn.Close()
		chas, err := shn.Receive(queueName)
		if err != nil {
			log.Infof("err: %s", err)
		}
		for cha := range chas {
			log.Debugf("cha :%v", *amqp.BytesToString(&(cha.Body)))
			cha.Ack(false)
			c ++
			log.Debugf("cha :%v", c)
			if c == 5 {
				return
			}
		}
	}()
	return nil
}

func (s *UserService) ReceiveFanout3() error {
	go func() {
		for {
			shn := s.newChannel()
			c, _ := shn.ReceiveFanout("test22222", exchange)
			if c != nil {
				log.Infof("cha: %s", *c)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	return nil
}
