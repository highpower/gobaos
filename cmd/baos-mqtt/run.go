package main

import (
	"fmt"
	"gobaos/internal/config"
)

func run(c *config.Data) error {
	logger, err := newLogger(c)
	if err != nil {
		return fmt.Errorf("main.run: cannot create logger[%w]", err)
	}
	logger.Infof("main.run: watchdog started with %s", c.String())
	return nil
}
