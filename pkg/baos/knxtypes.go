package baos

type (
	Control byte
	Dimming byte
)

func NewControl(manual bool, on bool) Control {
	return 0
}

func NewDimming(on bool, stepCode byte) Dimming {
	return 0
}
