package main

import (
	"fmt"
	mq "learn_micro/mq/demo2"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

	// 使用不同的producer并发publish
	for i := 0; i < 1; i++ {
		go func(idx int) {
			p, err := m.Producer(strconv.Itoa(i))
			if err != nil {
				log.Printf("[ERROR] Create producer failed, %v\n", err)
				return
			}
			if err = p.SetExchangeBinds(exb).Confirm(true).Open(); err != nil {
				log.Printf("[ERROR] Open failed, %v\n", err)
				return
			}

			// 使用同一个producer并发publish
			for j := 0; j < 10; j++ {
				go func(v int) {
					for {
						v++
						var routeKey = ""
						if v%2 == 0 {
							routeKey = "route.unitest2"
						} else {
							routeKey = "route.unitest1"
						}
						msg := mq.NewPublishMsg([]byte(fmt.Sprintf(`{"name":"zwf-%d"}`, v)))
						err = p.Publish("exch.unitest", routeKey, msg)
						if err != nil {
							log.Printf("[ERROR] %s\n", err.Error())
						}
						//log.Info("Producer(%d) state:%d, err:%v\n", i, p.State(), err)
					}
				}(j)
				time.Sleep(1 * time.Millisecond * 1000)
			}

		}(i)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

}
