package baos

// MainService is the top-level service identifier for BAOS messages.
const MainService byte = 0xF0

// BAOS sub-service identifiers.
const (
	// SubServiceGetServerItemReq is the request sub-service for reading server items.
	SubServiceGetServerItemReq byte = 0x01
	// SubServiceGetServerItemRes is the response sub-service for reading server items.
	SubServiceGetServerItemRes byte = 0x81
	// SubServiceSetServerItemReq is the request sub-service for writing server items.
	SubServiceSetServerItemReq byte = 0x02
	// SubServiceSetServerItemRes is the response sub-service for writing server items.
	SubServiceSetServerItemRes byte = 0x82
	// SubServiceGetDatapointDescriptionReq is the request sub-service for reading datapoint descriptions.
	SubServiceGetDatapointDescriptionReq byte = 0x03
	// SubServiceGetDatapointDescriptionRes is the response sub-service for reading datapoint descriptions.
	SubServiceGetDatapointDescriptionRes byte = 0x83
	// SubServiceGetDescriptionStringReq is the request sub-service for reading description strings.
	SubServiceGetDescriptionStringReq byte = 0x04
	// SubServiceGetDescriptionStringRes is the response sub-service for reading description strings.
	SubServiceGetDescriptionStringRes byte = 0x84
	// SubServiceGetDatapointValueReq is the request sub-service for reading datapoint values.
	SubServiceGetDatapointValueReq byte = 0x05
	// SubServiceGetDatapointValueRes is the response sub-service for reading datapoint values.
	SubServiceGetDatapointValueRes byte = 0x85
	// SubServiceSetDatapointValueReq is the request sub-service for writing datapoint values.
	SubServiceSetDatapointValueReq byte = 0x06
	// SubServiceSetDatapointValueRes is the response sub-service for writing datapoint values.
	SubServiceSetDatapointValueRes byte = 0x86
	// SubServiceDatapointValueInd is the indication sub-service for datapoint value updates.
	SubServiceDatapointValueInd byte = 0xC1
	// SubServiceServerItemInd is the indication sub-service for server item updates.
	SubServiceServerItemInd byte = 0xC2
)

// Message is the common interface for all BAOS protocol messages.
type Message interface {
	Service() byte
	SubService() byte
}
