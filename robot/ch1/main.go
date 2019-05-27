package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/keyboard"
	"gobot.io/x/gobot/platforms/mqtt"
)

func main() {

	mqttAdaptor := mqtt.NewAdaptor("tcp://0.0.0.0:1883", "pinger")

	work := func() {
		mqttAdaptor.On("hello", func(msg mqtt.Message) {
			fmt.Println(msg)
		})
		mqttAdaptor.On("hola", func(msg mqtt.Message) {
			fmt.Println(msg)
		})
		data := []byte("o")
		gobot.Every(1*time.Second, func() {
			mqttAdaptor.Publish("hello", data)
		})
		gobot.Every(5*time.Second, func() {
			mqttAdaptor.Publish("hola", data)
		})
	}

	robot := gobot.NewRobot("mqttBot",
		[]gobot.Connection{mqttAdaptor},
		work,
	)

	_ = robot.Start()
}

func mqttClient() {

}

func mqttServer() {
	mqttAdaptor := mqtt.NewAdaptor("tcp://0.0.0.0:1883", "pinger")
	work := func() {
		data := []byte("o")
		gobot.Every(1*time.Second, func() {
			mqttAdaptor.Publish("hello", data)
		})
		gobot.Every(5*time.Second, func() {
			mqttAdaptor.Publish("hola", data)
		})
	}

	robot := gobot.NewRobot("mqttBot",
		[]gobot.Connection{mqttAdaptor},
		work,
	)
	_ = robot.Start()
}

func keyboardFn() {
	keys := keyboard.NewDriver()
	work := func() {
		_ = keys.On(keyboard.Key, func(data interface{}) {
			key := data.(keyboard.KeyEvent)

			if key.Key == keyboard.A {
				fmt.Println("A pressed!")
			} else {
				fmt.Println("keyboard event!", key, key.Char)
			}
		})
	}

	robot := gobot.NewRobot("keyboardbot",
		[]gobot.Connection{},
		[]gobot.Device{keys},
		work,
	)

	robot.Start()
}
