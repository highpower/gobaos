//go:build !linux

package ft12

func flush(_ int) error {
	return nil
}
