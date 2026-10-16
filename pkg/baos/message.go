package baos

// mainService is the top-level service identifier for BAOS messages.
const mainService byte = 0xF0

// BAOS sub-service identifiers.
const (
	// subServiceGetServerItemReq is the request sub-service for reading server items.
	subServiceGetServerItemReq byte = 0x01
	// subServiceGetServerItemRes is the response sub-service for reading server items.
	subServiceGetServerItemRes byte = 0x81
	// subServiceSetServerItemReq is the request sub-service for writing server items.
	subServiceSetServerItemReq byte = 0x02
	// subServiceSetServerItemRes is the response sub-service for writing server items.
	subServiceSetServerItemRes byte = 0x82
	// subServiceGetDatapointDescriptionReq is the request sub-service for reading datapoint descriptions.
	subServiceGetDatapointDescriptionReq byte = 0x03
	// subServiceGetDatapointDescriptionRes is the response sub-service for reading datapoint descriptions.
	subServiceGetDatapointDescriptionRes byte = 0x83
	// subServiceGetDescriptionStringReq is the request sub-service for reading description strings.
	subServiceGetDescriptionStringReq byte = 0x04
	// subServiceGetDescriptionStringRes is the response sub-service for reading description strings.
	subServiceGetDescriptionStringRes byte = 0x84
	// subServiceGetDatapointValueReq is the request sub-service for reading datapoint values.
	subServiceGetDatapointValueReq byte = 0x05
	// subServiceGetDatapointValueRes is the response sub-service for reading datapoint values.
	subServiceGetDatapointValueRes byte = 0x85
	// subServiceSetDatapointValueReq is the request sub-service for writing datapoint values.
	subServiceSetDatapointValueReq byte = 0x06
	// subServiceSetDatapointValueRes is the response sub-service for writing datapoint values.
	subServiceSetDatapointValueRes byte = 0x86
	// subServiceDatapointValueInd is the indication sub-service for datapoint value updates.
	subServiceDatapointValueInd byte = 0xC1
	// subServiceServerItemInd is the indication sub-service for server item updates.
	subServiceServerItemInd byte = 0xC2
)
