package main

import (
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"
	"io"
	"net/url"
	"os"
	"time"
)

func run(name string) error {
	c, err := newConfig(name)
	if err != nil {
		return err
	}
	conn, err := ft12.NewConn(c.Device, ft12.WithTimeout(50*time.Millisecond))
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()
	return runWith(conn, &c)
}

func runWith(conn *ft12.Conn, c *config) error {
	worker := baos.NewWorker(conn)
	defer worker.Stop()

	u := url.URL{Scheme: "mqtts", Host: c.MqttHost}

	opts := mqtt.ClientOptions{
		Servers:       []*url.URL{&u},
		AutoReconnect: true,
		ClientID:      "baos-mqtt",
		CleanSession:  true,
		KeepAlive:     10,
	}
	_ = opts

	return nil
}

func usage(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "usage: baos-get -device <device> <id>\n")
	flag.CommandLine.SetOutput(writer)
	flag.CommandLine.PrintDefaults()
}

func main() {
	var conf string
	flag.StringVar(&conf, "conf", "", "config file")

	flag.Usage = func() { usage(os.Stdout) }
	flag.Parse()

	if !flag.Parsed() || conf == "" {
		usage(os.Stderr)
		os.Exit(2)
	}

	if err := run(conf); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
