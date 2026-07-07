package baos

type Service interface {
	Service() byte
	SubService() byte
}

const MainService byte = 0xF0

const (
	SubServiceGetServerItemReq     byte = 0x01
	SubServiceGetServerItemRes     byte = 0x81
	SubServiceSetServerItemReq     byte = 0x02
	SubServiceSetServerItemRes     byte = 0x82
	SubServiceGetDatapointValueReq byte = 0x05
	SubServiceGetDatapointValueRes byte = 0x85
	SubServiceSetDatapointValueReq byte = 0x06
	SubServiceSetDatapointValueRes byte = 0x86
	SubServiceDatapointValueInd    byte = 0xC1
	SubServiceServerItemInd        byte = 0xC2
)