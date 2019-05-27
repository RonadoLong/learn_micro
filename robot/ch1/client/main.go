package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/mqtt"
)

func main() {
	mqttClient()
}

func mqttClient() {
	mqttAdaptor := mqtt.NewAdaptor("tcp://0.0.0.0:1883", "pinger")

	work := func() {
		mqttAdaptor.On("hello", func(msg mqtt.Message) {
			fmt.Println(msg)
		})
		mqttAdaptor.On("hola", func(msg mqtt.Message) {
			fmt.Println(msg)
		})
	}

	robot := gobot.NewRobot("mqttBot",
		[]gobot.Connection{mqttAdaptor},
		work,
	)

	_ = robot.Start()
}
