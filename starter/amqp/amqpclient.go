package amqp

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"hidevops.io/hiboot/pkg/log"
)

type Channel struct {
	*amqp.Channel
}

type properties struct {
	Port      int    `json:"port" default:"5672"`
	Username  string `json:"username" default:"guest"`
	Password  string `json:"password" default:"guest"`
	Host      string `json:"host" default:"127.0.0.1"`
	QueueName string `json:"queueName" default:"my-queue"`
	Exchange  string `json:"exchange" default:"my-exchange"`
	SleepTime int64  `json:"sleepTime" default:"3*1e9"`
}

type AmqpClient interface {
	Connect() (cha *amqp.Channel, err error)
	Close()
}

func newChannel() (chn *Channel) {
	return new(Channel)
}

func (chn *Channel) Connect(p *properties) (err error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", p.Username, p.Password, p.Host, p.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Errorf("Failed to connect to RabbitMQ :%v", err)
		return
	}
	chn.Channel, err = conn.Channel()
	return err
}

func (chn *Channel) Receive(queueName string) (<-chan amqp.Delivery, error) {
	return chn.Channel.Consume(queueName,
		"",
		false,
		false,
		false,
		false,
		nil)

}

func (chn *Channel) ReceiveFanout(queueName, exchange string) (*string, error) {
	msgs, err := chn.Channel.Consume(queueName,
		"",
		false,
		false,
		false,
		false,
		nil)
	/*	msg, ok, err := chn.Get(queueName, true)
		if !ok {
			return nil, err
		}*/
	//err = s.channel.Ack(msg.DeliveryTag, false)
	if err != nil {
		return nil, err
	}
	// 使用callback消费数据
	for msg := range msgs {
		b := BytesToString(&(msg.Body))
		fmt.Sprintf("heee: %v", b)

		// 确认收到本条消息, multiple必须为false
		msg.Ack(false)
	}
	return nil, nil

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
	err := chn.Publish(exchange, "", false, false, amqp.Publishing{
		ContentType: "text/plain", Body: []byte(mgsConnect),
	})
	return err
}

func (chn *Channel) Push(exchange, key, expiration, mgsConnect string) error {
	err := chn.Publish(exchange, key, false, false, amqp.Publishing{
		ContentType: "text/plain", Body: []byte(mgsConnect), Expiration: expiration,
	})
	return err
}

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}

//CreateFanout 创建Fanout类型的队列
func (chn *Channel) CreateFanout(queueName, exchange string) error {
	//type : 交换器类型 DIRECT("direct"), FANOUT("fanout"), TOPIC("topic"), HEADERS("headers");
	//durable: 是否持久化,durable设置为true表示持久化,反之是非持久化
	err := chn.ExchangeDeclare(exchange, "fanout", false, false, false, false, nil)
	_, err = chn.QueueDeclare(queueName, false, false,
		false, false, nil)
	if err != nil {
		return err
	}
	err = chn.QueueBind(queueName, "", exchange, false, nil)
	return err
}

func (chn *Channel) Create(queueName, exchange, key, kind string) error {
	//type : 交换器类型 DIRECT("direct"), FANOUT("fanout"), TOPIC("topic"), HEADERS("headers");
	//durable: 是否持久化,durable设置为true表示持久化,反之是非持久化
	err := chn.ExchangeDeclare(exchange, kind, false, false, false, false, nil)
	_, err = chn.QueueDeclare(queueName, false, false,
		false, false, nil)
	if err != nil {
		return err
	}
	err = chn.QueueBind(queueName, key, exchange, false, nil)
	return err
}
