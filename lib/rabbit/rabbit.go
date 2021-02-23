package rabbit

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
)

type MsgClient struct {
	conn *amqp.Connection
}

func (m *MsgClient) ConnectToBroker(connectionString string) error {
	var err error
	m.conn, err = amqp.Dial(connectionString)
	if err != nil {
		return errors.New("rabbmit 连接错误: " + err.Error())
	}
	return nil
}

func (m *MsgClient) Publish(body []byte, exchangeName string, exchangeType string, queueName string) error {
	if m.conn == nil {
		return errors.New("amp conn连接错误")
	}
	ch, err := m.conn.Channel() // Get a channel from the connection
	if err != nil {
		return errors.New("MQ打开管道失败：" + err.Error())
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		true,         // noWait
		nil,          // arguments
	)
	if err != nil {
		return errors.New("MQ注册交换机失败：" + err.Error())
	}
	_, err = ch.QueueDeclare( // Declare a queue that will be created if not exists with some args
		queueName, // our queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		true,      // no-wait
		nil,       // arguments
	)
	if err != nil {
		return errors.New("MQ注册队列失败：" + err.Error())
	}

	//队列绑定
	err = ch.QueueBind(
		queueName,    // name of the queue
		exchangeName, // bindingKey
		exchangeName, // sourceExchange
		true,         // noWait
		nil,          // arguments
	)
	if err != nil {
		return errors.New("绑定队列失败：" + err.Error())
	}

	err = ch.Publish( // Publishes a message onto the queue.
		exchangeName, // exchange
		exchangeName, // routing key      q.Name
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Body: body, // Our JSON body as []byte
		})
	fmt.Printf("A message was sent: %v", string(body))

	if err != nil {
		return errors.New("消息推送失败：" + err.Error())
	}

	return nil
}

func (m *MsgClient) PublishOnQueue(body []byte, queueName string) error {
	if m.conn == nil {
		return errors.New("amp conn连接错误")
	}
	ch, err := m.conn.Channel() // Get a channel from the connection
	if err != nil {
		return errors.New("MQ打开管道失败：" + err.Error())
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare( // Declare a queue that will be created if not exists with some args
		queueName, // our queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		true,      // no-wait
		nil,       // arguments
	)
	if err != nil {
		return errors.New("MQ注册队列失败：" + err.Error())
	}

	// Publishes a message onto the queue.
	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body, // Our JSON body as []byte
		})
	fmt.Printf("A message was sent: %v", string(body))

	if err != nil {
		return errors.New("消息推送失败：" + err.Error())
	}

	return nil
}

func (m *MsgClient) GetSubscribe(exchangeName string, exchangeType string, consumerName string, queueName string) (<-chan amqp.Delivery, error) {
	ch, err := m.conn.Channel()
	if err != nil {
		return nil, errors.New("MQ打开管道失败：" + err.Error())
	}

	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		true,         // noWait
		nil,          // arguments
	)
	if err != nil {
		return nil, errors.New("MQ注册交换机失败：" + err.Error())
	}

	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err = ch.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		true,      // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, errors.New("MQ注册队列失败：" + err.Error())
	}
	err = ch.QueueBind(
		queueName,    // name of the queue
		exchangeName, // bindingKey
		exchangeName, // sourceExchange
		true,         // noWait
		nil,          // arguments
	)
	if err != nil {
		return nil, errors.New("绑定队列失败：" + err.Error())
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	msgs, err := ch.Consume(
		queueName,    // queue
		consumerName, // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		true,         // no-wait
		nil,          // args
	)
	if err != nil {
		return nil, errors.New("获取消费通道异常：" + err.Error())
	}
	return msgs, nil
}

func (m *MsgClient) GetSubscribeToQueue(queueName string, consumerName string) (<-chan amqp.Delivery, error) {
	ch, err := m.conn.Channel()
	if err != nil {
		return nil, errors.New("MQ打开管道失败：" + err.Error())
	}

	queue, err := ch.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		true,      // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, errors.New("MQ注册队列失败" + err.Error())
	}

	msgs, err := ch.Consume(
		queue.Name,   // queue
		consumerName, // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return nil, errors.New("获取消费通道异常" + err.Error())
	}

	return msgs, nil
}

func (m *MsgClient) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
}

//延时队列相关
func (m *MsgClient) PublishToDelayed(body []byte, exchangeName string, exchangeType string, routingKey string, seconds int) error {
	if m.conn == nil {
		return errors.New("amp conn连接错误")
	}
	ch, err := m.conn.Channel() // Get a channel from the connection
	if err != nil {
		return errors.New("MQ打开管道失败：" + err.Error())
	}
	defer ch.Close()

	args := amqp.Table{
		"x-delayed-type": exchangeType,
	}
	err = ch.ExchangeDeclare(
		exchangeName,        // name of the exchange
		"x-delayed-message", // type
		true,                // durable
		false,               // delete when complete
		false,               // internal
		true,                // noWait
		args,                // arguments
	)
	if err != nil {
		return errors.New("MQ注册交换机失败：" + err.Error())
	}

	err = ch.Publish( // Publishes a message onto the queue.
		exchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Body:    body, // Our JSON body as []byte
			Headers: amqp.Table{"x-delay": seconds * 1000},
		})

	if err != nil {
		return errors.New("消息推送失败：" + err.Error())
	}
	return nil
}

// 从延时队列中获取消息
func (m *MsgClient) GetSubscribeOnDelayed(exchangeName string, exchangeType string, queueName string, consumerName string, routingKey string) (<-chan amqp.Delivery, error) {
	ch, err := m.conn.Channel()
	if err != nil {
		return nil, errors.New("MQ打开管道失败：" + err.Error())
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"x-delayed-message",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-delayed-type": exchangeType,
		})
	if err != nil {
		return nil, errors.New("MQ注册交换机失败：" + err.Error())
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, errors.New("MQ注册队列失败：" + err.Error())
	}

	err = ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
	if err != nil {
		return nil, errors.New("MQ绑定队列失败：" + err.Error())
	}

	msgs, err := ch.Consume(
		q.Name,       // queue
		consumerName, // consumer
		false,        // auto ack
		false,        // exclusive
		false,        // no local
		false,        // no wait
		nil,          // args
	)
	if err != nil {
		return nil, errors.New("获取消费通道异常：" + err.Error())
	}

	return msgs, nil
}
