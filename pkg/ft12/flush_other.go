//go:build !linux

package ft12

import "os"

func flush(_ *os.File) error {
	return nil
}
