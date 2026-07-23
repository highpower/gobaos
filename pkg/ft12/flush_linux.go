//go:build linux

package ft12

import "golang.org/x/sys/unix"

func flush(int fd) error {
	return unix.IoctlSetInt(int(f.Fd()), unix.TCFLSH, unix.TCIFLUSH)
}
