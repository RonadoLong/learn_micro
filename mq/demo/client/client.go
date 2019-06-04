package main

import (
	"context"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var (
	RUrl        = "amqp://mq_dev:abc123@ec2-13-230-248-136.ap-northeast-1.compute.amazonaws.com/test"
	CTPqueue    = "queue.ctp"
	CTPConsumer = "ctp-consumer"
	CTPRouteKey = "route.ctp"
	exchange    = "test-exchange"
)

func main() {

	msgC := make([]chan amqp.Delivery, 20)

	forever := make(chan bool)
	for i := range msgC {
		log.Println(i)
		var msgChan = make(chan amqp.Delivery, 10)
		go worker(msgChan)
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

func worker(msgC chan amqp.Delivery) {
	consumer := NewConsumer(RUrl, exchange, CTPqueue, CTPRouteKey, 20)
	err := consumer.SetQps(10).SetMsgCallBack(msgC).Connect()
	if err != nil {
		panic(err)
	}
	var index = 0
	for {
		select {
		case d, ok := <-msgC:
			if ok {
				index++
				if index == 10000000 {
					//consumer.CloseChannel()
					consumer.cancelFunc()
					log.Println("模拟断开=======")
				}
				//log.Printf("Received a message: %s", d.Body)
				_ = d.Ack(true) // 手动ACK
			}
		}
	}
}

const (
	reconnectDelay = 5 * time.Second // 连接断开后多久重连
	resendDelay    = 5 * time.Second // 消息发送失败后，多久重发
	resendTime     = 3               // 消息重发次数
)

type Consumer struct {
	Url                       string
	Exchange, Queue, RouteKey string
	closeChan                 chan *amqp.Error
	isConnect                 bool
	ch                        *amqp.Channel
	prefetch                  int
	callBack                  chan<- amqp.Delivery
	ctx                       context.Context
	cancelFunc                context.CancelFunc
	conn                      *amqp.Connection
	ReceiverNum               int
}

func NewConsumer(Url, Exchange, Queue, RouteKey string, ReceiverNum int) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Consumer{
		Url:         Url,
		Exchange:    Exchange,
		Queue:       Queue,
		RouteKey:    RouteKey,
		closeChan:   make(chan *amqp.Error, 1),
		ctx:         ctx,
		cancelFunc:  cancel,
		ReceiverNum: 1,
	}
}

func (c *Consumer) Connect() error {
	var err error
	c.conn, err = amqp.Dial(c.Url)
	if err != nil {
		failOnError(err, "Failed to Connect to RabbitMQ")
		return err
	}
	for i := 0; i < c.ReceiverNum; i++ {
		go c.receiverMsg()
	}

	return nil
}

func (c *Consumer) receiverMsg() {
	var err error
	ch, err := c.conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel")
	}

	//c.ch.NotifyClose(c.closeChan)
	_ = ch.Qos(c.prefetch, 0, false)

	_ = ch.ExchangeDeclare(
		c.Exchange,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	)
	q, err := ch.QueueDeclare(
		c.Queue, // name
		true,    // durable  是否持久化
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		failOnError(err, "Failed to declare a queue")
	}
	_ = ch.QueueBind(q.Name, c.RouteKey, exchange, false, nil)
	delivery, _ := ch.Consume(
		c.Queue,     // queue
		CTPConsumer, // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	for {
		select {
		case res, ok := <-delivery:
			if ok {
				c.callBack <- res
			}
		case <-c.ctx.Done():
			log.Println("退出=============")
			return
		}
	}
}

func (c *Consumer) listenerClose() {
	for {
		res, ok := <-c.closeChan
		if ok {
			log.Println("reconnect ", res)
			c.cancelFunc()
		}
	}
}

func (c *Consumer) SetMsgCallBack(callBack chan<- amqp.Delivery) *Consumer {
	c.callBack = callBack
	return c
}

func (c *Consumer) SetQps(prefetch int) *Consumer {
	c.prefetch = prefetch
	return c
}

func (c *Consumer) CloseChannel() {
	c.ch.Close()
}
