package main

import mqtt "github.com/eclipse/paho.golang/paho"

func main() {
	opts := mqtt.ClientConfig{
		ClientID: "baos-mqtt",
	}
	_ = opts
}
