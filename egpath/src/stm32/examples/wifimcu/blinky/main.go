package main

import (
	"delay"

	"stm32/hal/gpio"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var led gpio.Pin

func init() {
	system.Setup96(26)
	systick.Setup(2e6)

	gpio.A.EnableClock(false)
	led = gpio.A.Pin(4)

	cfg := gpio.Config{
		Mode:   gpio.Out,
		Driver: gpio.OpenDrain,
		Speed:  gpio.Low,
	}
	led.Setup(&cfg)
}

func main() {
	for {
		led.Clear()
		delay.Millisec(50)
		led.Set()
		delay.Millisec(950)
	}
}
