package main

import (
	"fmt"
	"thermo-watchdog/internal/config"
	"thermo-watchdog/internal/core"
)

func run(c *config.Data) error {
	logger, err := newLogger(c)
	if err != nil {
		return fmt.Errorf("main.run: cannot create logger[%w]", err)
	}
	watchdog, err := core.NewWatchdog(c, c, logger)
	if err != nil {
		return fmt.Errorf("main.run: cannot create watchdog[%w]", err)
	}
	logger.Infof("main.run: watchdog started with %s", c.String())
	return watchdog.Run()
}
