package baos

import (
	"gobaos/pkg/ft12"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runDevice(wg *sync.WaitGroup, file *os.File) {
	defer wg.Done()
	defer func() { _ = file.Close() }()
	_, _ = file.Write([]byte{0x4c})
}

func openNamedPipe(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR, 0o644|os.ModeNamedPipe)
}

func TestFailingDevice(t *testing.T) {

	const name = "testdata/failing"
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0o644|os.ModeNamedPipe)
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)

	go runDevice(&wg, file)
	conn, err := ft12.NewConn(name, ft12.WithOpenFileFunc(openNamedPipe))
	require.NoError(t, err)

	worker := NewWorker(conn)
	defer worker.Stop()

	assert.Error(t, worker.Run())
	wg.Wait()
}
