package main

import (
	mq "learn_micro/mq/demo2"
	"log"
	"time"
)

var (
	MQURL = "amqp://mq_dev:abc123@ec2-13-230-248-136.ap-northeast-1.compute.amazonaws.com:5672/test"
)

func main() {
	m, err := mq.New(MQURL).Open()
	if err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return
	}
	defer m.Close()

	c, err := m.Consumer("test-consume")
	if err != nil {
		log.Printf("[ERROR] Create consumer failed, %v\n", err)
		return
	}
	defer c.Close()

	exb := []*mq.ExchangeBinds{
		&mq.ExchangeBinds{
			Exch: mq.DefaultExchange("exch.unitest", mq.ExchangeDirect),
			Bindings: []*mq.Binding{
				&mq.Binding{
					RouteKey: "route.unitest1",
					Queues: []*mq.Queue{
						mq.DefaultQueue("queue.unitest1"),
					},
				},
				&mq.Binding{
					RouteKey: "route.unitest2",
					Queues: []*mq.Queue{
						mq.DefaultQueue("queue.unitest2"),
					},
				},
			},
		},
	}

	msgC := make(chan mq.Delivery, 1)
	defer close(msgC)

	c.SetExchangeBinds(exb).SetMsgCallback(msgC).SetQos(20)
	//c.SetMsgCallback(msgC)
	//c.SetQos(2)
	if err = c.Open(); err != nil {
		log.Printf("[ERROR] Open failed, %v\n", err)
		return
	}

	var i = 0
	for msg := range msgC {
		//log.Printf("Tag(%d) Body: %s Route: %s\n", msg.DeliveryMode, string(msg.Body), 		msg.RoutingKey)
		i++
		_ = msg.Ack(true)

		//if i%5 == 0 {
		//	_ = msg.Reject(true)
		//}
		//log.Info("Consumer receive msg `%s`", string(msg))
		time.Sleep(time.Millisecond * 1000)
	}

}
