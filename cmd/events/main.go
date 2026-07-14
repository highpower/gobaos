package main

import (
	"fmt"
	"gobaos/pkg/ft12"
	"os"
)

func main() {
	conn, err := ft12.NewConn("/dev/ttyAMA0")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "connect error: %s\n", err.Error())
		os.Exit(1)
	}
	defer func() { _ = conn.Close() }()

	for {
		if _, err := conn.Read(nil); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "read error: %s\n", err.Error())
			os.Exit(1)
		}
		l := conn.Available()
		buffer := make([]byte, l)
		_, err := conn.Read(buffer)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "buffer read error: %s\n", err.Error())
		}
	}
}
