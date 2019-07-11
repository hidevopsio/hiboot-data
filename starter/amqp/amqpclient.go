package amqp

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"hidevops.io/hiboot/pkg/log"
	"time"
)

type Channel struct {
	*amqp.Channel
}

type properties struct {
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Host      string `json:"host"`
	QueueName string `json:"queueName"`
	Exchange  string `json:"exchange"`
	SleepTime int64 `json:"sleepTime" default:"3*1e9"`
}

type AmqpClient interface {
	Connect() (cha *amqp.Channel, err error)
	Close()
}

func Connect(p *properties) (chn *Channel, err error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", p.Username, p.Password, p.Host, p.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Errorf("Failed to connect to RabbitMQ :%v", err)
		return
	}
	chn = new(Channel)
	chn.Channel, err = conn.Channel()

	return chn, err
}

func (chn *Channel) Receive(queueName string) {
	for {
		msg, ok, _ := chn.Channel.Get(queueName, true)
		if !ok {
			time.Sleep(3*1e9)
			continue
		}
		//err = s.channel.Ack(msg.DeliveryTag, false)
		b := BytesToString(&(msg.Body))
		fmt.Printf("receve msg is :%s\n", *b)
	}
}

func (chn *Channel) ReceiveFanout(queueName, exchange string) (*string, error) {
	for {
		_, err := chn.QueueDeclare(queueName, false, false,
			false, false, nil)
		if err != nil {
			return nil, err
		}
		err = chn.QueueBind(queueName, "", exchange, false, nil)
		if err != nil {
			return nil, err
		}
		msg, ok, err := chn.Get(queueName, true)
		if !ok {
			time.Sleep(3*1e9)
			continue
		}

		//err = s.channel.Ack(msg.DeliveryTag, false)
		b := BytesToString(&(msg.Body))
		return b, nil
	}


}


func (chn *Channel) PublishDirect(exchange, queueName, mgsConnect, key string) error {
	//type : 交换器类型 DIRECT("direct"), FANOUT("fanout"), TOPIC("topic"), HEADERS("headers");
	//durable: 是否持久化,durable设置为true表示持久化,反之是非持久化
	err := chn.ExchangeDeclare(exchange, "direct", false, false, false, false, nil)
	if err != nil {
		return err
	}
	_, err = chn.QueueDeclare(queueName, false, false,
		false, false, nil)

	err = chn.QueueBind(queueName, key, exchange, false, nil)

	err = chn.Publish(exchange, key, false, false, amqp.Publishing{
		ContentType: "text/plain", Body: []byte(mgsConnect),
	})
	return err
}

func (chn *Channel) PublishFanout(exchange, mgsConnect string) error {
	//type : 交换器类型 DIRECT("direct"), FANOUT("fanout"), TOPIC("topic"), HEADERS("headers");
	//durable: 是否持久化,durable设置为true表示持久化,反之是非持久化
	err := chn.ExchangeDeclare(exchange, "fanout", false, false, false, false, nil)
	if err != nil {
		return err
	}
	err = chn.Publish(exchange, "", false, false, amqp.Publishing{
		ContentType: "text/plain", Body: []byte(mgsConnect),
	})
	return err
}

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}