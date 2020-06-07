package hserial

type Mode struct {
	baudRate int      // The serial port bitrate (aka Baudrate)
	dataBits int      // Size of the character (must be 5, 6, 7 or 8)
	parity   Parity   // Parity (see Parity type for more info)
	stopBits StopBits // Stop bits (see StopBits type for more info)
}

func NewMode(baudRate int, dataBits int, parity Parity, stopBits StopBits) *Mode {
	model := new(Mode)
	model.baudRate = baudRate
	model.dataBits = dataBits
	model.parity = parity
	model.stopBits = stopBits
	return model
}
