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
	channel *amqp.Channel
}

func init() {
	app.Register(newUserService)
}

// will inject BoltRepository that configured in hidevops.io/hiboot/pkg/starter/data/bolt
func newUserService(channel *amqp.Channel) *UserService {
	return &UserService{
		channel: channel,
	}
}

const (
	mgsConnect = "hello world"
	exchange   = "test11"
	queueName  = "Test"
)

func (s *UserService) PublishDirect() error {
	err := s.channel.PublishDirect(exchange, queueName, "hello", "info")
	return err
}

func (s *UserService) PublishFanout() error {
	err := s.channel.PublishFanout(exchange, "hello")
	return err
}

//
//func (s *UserService) Receive() {
//	for {
//		msg, ok, err := s.channel.Get(queueName, true)
//
//		if !ok {
//			fmt.Println("do not get msg")
//			time.Sleep(3*1e9)
//			continue
//		}
//
//		//err = s.channel.Ack(msg.DeliveryTag, false)
//		log.Infof("err :%v", err)
//
//		b := BytesToString(&(msg.Body))
//		fmt.Printf("receve msg is :%s\n", *b)
//	}
//
//
//}

func (s *UserService) ReceiveFanout() error {
	go func() {
		for {
			c, _ := s.channel.ReceiveFanout("test22", exchange)
			if c != nil {
				log.Infof("cha: %s", *c)
			}
			time.Sleep(5 * time.Second)

		}
	}()
	return nil

}

func (s *UserService) ReceiveFanout3() error {
	go func() {
		for {
			c, _ := s.channel.ReceiveFanout("test22222", exchange)
			if c != nil {
				log.Infof("cha: %s", *c)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	return nil

}

//func (s *UserService) ReceiveFanout1() {
//	for {
//		_, err := s.channel.QueueDeclare("test1", false, false,
//			false, false, nil)
//
//		err = s.channel.QueueBind("test1", "", exchange, false, nil)
//		msg, ok, err := s.channel.Get("test1", true)
//		if !ok {
//			fmt.Println("do not get msg1")
//			time.Sleep(3*1e9)
//			continue
//		}
//
//		//err = s.channel.Ack(msg.DeliveryTag, false)
//		log.Infof("err :%v", err)
//
//		b := BytesToString(&(msg.Body))
//		fmt.Printf("receve msg is1 :%s\n", *b)
//	}
//
//
//}
//
//func BytesToString(b *[]byte) *string {
//	s := bytes.NewBuffer(*b)
//	r := s.String()
//	return &r
//}
