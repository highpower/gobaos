//go:build linux

package ft12

import "fmt"
import "os"
import "golang.org/x/sys/unix"

func flush(file *os.File) error {
	st, err := file.Stat()
	if err != nil {
		
		return fmt.Errorf("stat error: %w", err)
	}
	mode := st.Mode()
	fmt.Printf("file mode is %d\n", mode)
	if mode&os.ModeNamedPipe != 0 {
		return nil
	}
	
	return unix.IoctlSetInt(int(file.Fd()), unix.TCFLSH, unix.TCIFLUSH)
}
