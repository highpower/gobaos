//go:build linux

package ft12

import "golang.org/x/sys/unix"

func flush(fd int) error {
	return unix.IoctlSetInt(fd, unix.TCFLSH, unix.TCIFLUSH)
}
