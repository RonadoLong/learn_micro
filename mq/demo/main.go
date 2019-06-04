package main

import (
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
		return
	}
}

func main() {
	conn, err := amqp.Dial("amqp://mq_dev:abc123@ec2-13-230-248-136.ap-northeast-1.compute.amazonaws.com/test")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	var exchange = "test-ctp"
	err = ch.ExchangeDeclare(
		exchange,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a ExchangeDeclare")

	q, err := ch.QueueDeclare(
		"queue.ctp", // name
		true,        // durable  是否持久化
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(q.Name, "route.ctp", exchange, false, nil)
	failOnError(err, "Failed to declare a QueueBind")

	confirmChan := make(chan amqp.Confirmation, 1)
	ch.NotifyPublish(confirmChan)

	returnChan := make(chan amqp.Return, 1)
	ch.NotifyReturn(returnChan)

	closeChan := make(chan *amqp.Error, 1)
	ch.NotifyClose(closeChan)

	_ = ch.Confirm(false)

	forever := make(chan bool)
	go func() {
		for d := range confirmChan {
			log.Printf("confirmChan call: %+v \n", d)
		}
	}()

	go func() {
		for r := range returnChan {
			log.Printf("returnChan call: %+v \n", r)
		}
	}()

	//go func() {
	//	var index = 1
	//	for {
	//		index++
	//		mID := strconv.FormatInt(int64(index), 10)
	//		body := fmt.Sprintf("Hello World ------ %d !", index)
	//		err = ch.Publish(
	//			exchange,     // exchange
	//			"route.ctp", // routing key
	//			true,  // mandatory
	//			false,  // immediate
	//			amqp.Publishing{
	//				MessageId:    mID,
	//				DeliveryMode: 2,
	//				ContentType:  "text/plain",
	//				Body:         []byte(body),
	//			})
	//		failOnError(err, "Failed to publish a message")
	//	}
	//}()
	//
	//go func() {
	//	var index = 1000000
	//	for {
	//		index++
	//		mID := strconv.FormatInt(int64(index), 10)
	//		body := fmt.Sprintf("Hello World ------ %d !", index)
	//		err = ch.Publish(
	//			exchange,     // exchange
	//			"route.ctp", // routing key
	//			false,  // mandatory
	//			false,  // immediate
	//			amqp.Publishing{
	//				MessageId:    mID,
	//				DeliveryMode: 2,
	//				ContentType:  "text/plain",
	//				Body:         []byte(body),
	//			})
	//		failOnError(err, "Failed to publish a message")
	//	}
	//}()

	//go func() {
	//	var index = 100000000
	//	for {
	//		index++
	//		mID := strconv.FormatInt(int64(index), 10)
	//
	//		body := fmt.Sprintf("Hello World ------ %d !", index)
	//		err = ch.Publish(
	//			"",     // exchange
	//			q.Name, // routing key
	//			false,  // mandatory
	//			false,  // immediate
	//			amqp.Publishing{
	//				MessageId: mID,
	//				ContentType: "text/plain",
	//				Body:        []byte(body),
	//			})
	//		failOnError(err, "Failed to publish a message")
	//	}
	//}()

	<-forever

}
