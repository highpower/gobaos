package ft12

import (
	"syscall"
)

func flock(fd int) error {
	return syscall.Flock(fd, syscall.LOCK_EX | syscall.LOCK_NB)
}
