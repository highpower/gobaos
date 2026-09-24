package main

import (
	"context"
	"flag"
	"fmt"
	"gobaos/internal/config"
	"gobaos/pkg/baos"
	"gobaos/pkg/ft12"
	"io"
	"net/url"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

func run(name string) error {
	c, err := config.Parse(name)
	if err != nil {
		return err
	}
	conn, err := ft12.NewConn(c.Device, ft12.WithTimeout(50*time.Millisecond))
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	logger, err := newLogger(c.Logger.Level, c.Logger.Output)
	if err != nil {
		return err
	}

	tc, err := newTLSConfig(c.Mqtt.CaCert)
	if err != nil {
		return err
	}

	opts := mqtt.ClientOptions{
		Servers:        []*url.URL{{Scheme: "tcps", Host: c.Mqtt.Host}},
		AutoReconnect:  true,
		ClientID:       "baos-mqtt",
		CleanSession:   true,
		ConnectTimeout: time.Second,
		KeepAlive:      10,
		TLSConfig:      tc,
		Username:       c.Mqtt.User,
		Password:       c.Mqtt.Password,
	}
	client := mqtt.NewClient(&opts)
	return runWith(conn, client, logger)
}

func runWorker(w *baos.Worker, wg *sync.WaitGroup, log *zap.SugaredLogger) {
	defer wg.Done()
	if err := w.Run(); err != nil {
		log.Errorf("main.runWorker: %s", err.Error())
	}
}

func runWith(conn *ft12.Conn, client mqtt.Client, log *zap.SugaredLogger) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	worker := baos.NewWorker(conn)
	defer worker.Stop()

	wg.Add(1)
	go runWorker(worker, &wg, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := worker.AllDescriptions(ctx)
	if err != nil {
		log.Errorf("main.runWith: descriptions error %s", err.Error())
		return err
	}

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		err := token.Error()
		log.Errorf("main.runWith: broker error %s", err.Error())
		return token.Error()
	}
	log.Infof("main.runWith: connected to mqtt")
	return nil
}

func usage(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "usage: baos-mqtt -config <config file>\n")
	flag.CommandLine.SetOutput(writer)
	flag.CommandLine.PrintDefaults()
}

func main() {
	var conf string
	flag.StringVar(&conf, "config", "", "config file")

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
