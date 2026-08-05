package baos

import "gobaos/pkg/knxtypes"

// CreateForID creates a new knxtypes.DPT instance for the given datapoint ID by reading its description.
func CreateForID(worker *Worker, id DatapointID) (knxtypes.DPT, error) {
	description, err := worker.Description(id)
	if err != nil {
		return nil, err
	}
	return knxtypes.NewDPT(description.DPT)
}
