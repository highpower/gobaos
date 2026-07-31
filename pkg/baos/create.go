package baos

import "gobaos/pkg/knxtypes"

func CreateForID(worker *Worker, id DatapointID) (knxtypes.DPT, error) {
	description, err := worker.Description(id)
	if err != nil {
		return nil, err
	}
	return knxtypes.NewDPT(description.DPT)
}
