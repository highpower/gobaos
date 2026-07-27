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
	case knxtypes.BoolCode:
		m[id] = makeHandler(printDPT[knxtypes.DPTBool])
	case knxtypes.ControlCode:
		m[id] = makeHandler(printDPT[knxtypes.DPTControl])
	case knxtypes.StepControlCode:
		m[id] = makeHandler(printDPT[knxtypes.DPTStepControl])
	case knxtypes.CharCode:
		m[id] = makeHandler(printDPT[knxtypes.DPTChar])
	case knxtypes.Uint8Code:
		m[id] = makeHandler(printDPT[knxtypes.DPTUint8])
	case knxtypes.Int8Code:
		m[id] = makeHandler(printDPT[knxtypes.DPTInt8])
	case knxtypes.Uint16Code:
		m[id] = makeHandler(printDPT[knxtypes.DPTUint16])
	case knxtypes.Int16Code:
		m[id] = makeHandler(printDPT[knxtypes.DPTInt16])
	case knxtypes.FloatCode:
		m[id] = makeHandler(printDPT[knxtypes.DPTFloat])
	default:
		return fmt.Errorf("unsupported type %s", code.String())
	}
	return nil
}
