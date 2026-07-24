package baos

// Command represents a datapoint command code.
type Command byte

// Supported datapoint commands.
const (
	NoCommand         Command = 0b0000
	SetNewValue       Command = 0b0001
	SendValueOnBus    Command = 0b0010
	SetAndSendOnBus   Command = 0b0011
	ReadOnBus         Command = 0b0100
	ClearTransmission Command = 0b0101
)
