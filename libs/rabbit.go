package libs

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"

	"github.com/cloudsrc/api.awaymail.v1.go/config"
)

type RabbitClient struct {
	AMQP        *amqp.Connection
	NotifyClose chan *amqp.Error
}

func newRabbitclient(config config.Config) *RabbitClient {
	var rmqpConnection = fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		config.Rabbit.User,
		config.Rabbit.Password,
		config.Rabbit.Host,
		config.Rabbit.Port,
		config.Rabbit.Vhost,
	)
	rmqConn, err := amqp.Dial(rmqpConnection)
	if err != nil {
		panic(err)
	}
	return &RabbitClient{
		AMQP:        rmqConn,
		NotifyClose: rmqConn.NotifyClose(make(chan *amqp.Error)),
	}
}

func (r *RabbitClient) Ping() (err error) {
	select { //check connection
	case err = <-r.NotifyClose:
		return err
	default:
		return nil
	}
}

func (r *RabbitClient) Publish(qName string, data interface{}) error {
	ch, err := r.AMQP.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(qName, false, true, false, false, nil)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(data)
	err = ch.Publish("", q.Name, false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	if err != nil {
		fmt.Println("Publish err: ", err.Error())
		return err
	}
	return nil
}

func (r *RabbitClient) Consume(qName string, f func([]byte) error) error {
	ch, err := r.AMQP.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(qName, false, true, false, false, nil)
	if err != nil {
		return err
	}
	if err := ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return err
	}
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			if err := f(d.Body); err != nil {
				var errBody interface{}
				json.Unmarshal(d.Body, &errBody)
				r.Publish(qName+"_error", errBody)
			}
			d.Ack(false)
		}
	}()
	<-forever
	return nil
}

func (r *RabbitClient) Close() (err error) {
	return r.AMQP.Close()
}
