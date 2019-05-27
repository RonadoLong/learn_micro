package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"time"
)

func main() {
	r := raspi.NewAdaptor()
	led := gpio.NewLedDriver(r, "7")

	work := func() {
		gobot.Every(1*time.Second, func() {
			err := led.Toggle()
			if err != nil {
				fmt.Println(err.Error())
			}
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{r},
		[]gobot.Device{led},
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
}
