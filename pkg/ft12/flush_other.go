//go:build !linux

package ft12

func flush(fd int) error {
	return nil
}
