package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/mqtt"
)

func main() {
	// webview.Open("Minimal webview example",
	// 	"http://localhost:8080", 1500, 900, true)
	bot()
}

func bot() {
	mqttAdaptor := mqtt.NewAdaptor(os.Args[1], "pinger")
	mqttAdaptor.SetAutoReconnect(true)

	holaDriver := mqtt.NewDriver(mqttAdaptor, "hola")
	helloDriver := mqtt.NewDriver(mqttAdaptor, "hello")

	work := func() {
		helloDriver.On(mqtt.Data, func(data interface{}) {
			fmt.Println("hello")
		})

		holaDriver.On(mqtt.Data, func(data interface{}) {
			fmt.Println("hola")
		})

		data := []byte("o")
		gobot.Every(1*time.Second, func() {
			helloDriver.Publish(data)
		})

		gobot.Every(5*time.Second, func() {
			holaDriver.Publish(data)
		})
	}

	robot := gobot.NewRobot("mqttBot",
		[]gobot.Connection{mqttAdaptor},
		[]gobot.Device{helloDriver, holaDriver},
		work,
	)

	robot.Start()
}
