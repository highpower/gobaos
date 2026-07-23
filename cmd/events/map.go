package main

import (
	"fmt"
	"gobaos/pkg/baos"
	"gobaos/pkg/knxtypes"
	"os"
)

type dataFunc func(id baos.DatapointID, data []byte) error

type dptWrapper[G any] interface {
	*G
	knxtypes.DPT
}

type stringerWrapper[G any] interface {
	*G
	fmt.Stringer
}

func handle[G any, W dptWrapper[G]](id baos.DatapointID, data []byte, f func(id baos.DatapointID, dpt G) error) error {
	var dpt G
	if err := W(&dpt).UnmarshalBinary(data); err != nil {
		return err
	}
	return f(id, dpt)
}

func makeHandler[G any, W dptWrapper[G]](f func(id baos.DatapointID, dpt G) error) dataFunc {
	return func(id baos.DatapointID, data []byte) error {
		return handle[G, W](id, data, f)
	}
}

func printDPT[G any, W stringerWrapper[G]](id baos.DatapointID, value G) error {
	_, _ = fmt.Fprintf(os.Stdout, "datapoint %d: %s\n", id, W(&value).String())
	return nil
}

type datapointMap map[baos.DatapointID]dataFunc

func (m datapointMap) handleIndication(ind *baos.DatapointValueInd) error {
	for _, v := range ind.Values {
		f, ok := m[v.ID]
		if !ok {
			continue
		}
		if err := f(v.ID, v.Data); err != nil {
			return err
		}
	}
	return nil
}

func (m datapointMap) add(id baos.DatapointID, t string) error {
	code, err := knxtypes.ParseTypeCode(t)
	if err != nil {
		return err
	}
	switch code {
	case knxtypes.Bool:
		m[id] = makeHandler(printDPT[knxtypes.DPTBool])
	case knxtypes.Control:
		m[id] = makeHandler(printDPT[knxtypes.DPTControl])
	case knxtypes.StepControl:
		m[id] = makeHandler(printDPT[knxtypes.DPTStepControl])
	case knxtypes.Float:
		m[id] = makeHandler(printDPT[knxtypes.DPTFloat])
	}
	return nil
}
