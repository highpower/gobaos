package main

import (
	"flag"
	mqtt "github.com/eclipse/paho.golang/paho"
	"gobaos/pkg/ft12"
)

func runWith(conn *ft12.Conn) error {
	opts := mqtt.ClientConfig{
		ClientID: "baos-mqtt",
	}
	_ = opts
	return nil
}

func main() {
	var conf string
	flag.StringVar(&conf, "conf", "", "config file")

	_ = runWith
}
