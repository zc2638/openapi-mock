package rabbit

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"log"
)

/**
 * Created by zc on 2019-11-28.
 */
var conn *amqp.Connection

func init() {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
}

type Channel struct {
	produce *amqp.Channel
	consume map[string]*amqp.Channel
}

func (c *Channel) newProduce(name string) error {
	var err error
	if c.produce == nil {
		c.produce, err = conn.Channel()
	}
	return err
}

func (c *Channel) newConsume(name string) error {
	var err error
	if c.consume == nil {
		c.consume = make(map[string]*amqp.Channel)
		c.consume[name], err = conn.Channel()
	}
	return err
}

func (c *Channel) Receive(queue string) (<-chan amqp.Delivery, error) {

	if err := c.newConsume(queue); err != nil {
		return nil, err
	}
	ch := c.consume[queue]
	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return ch.Consume(
		q.Name, // queue
		"mock", // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}

func (c *Channel) Send(exchange, queue string, data interface{}) error {
	if err := c.newProduce(queue); err != nil {
		return err
	}
	ch := c.produce
	defer ch.Close()
	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = ch.Publish(
		exchange, // 交换器名称
		q.Name,   // 路由键
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body: body,
		});
	if err != nil {
		return errors.New("failed to publish a message")
	}

	log.Printf(" [x] Sent %s", body)
	return nil
}