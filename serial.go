package hserial

import (
	"github.com/albenik/go-serial/v2"
)

type SerialPort struct {
	sport *serial.Port
}

func Open(portName string, mode *Mode) (*SerialPort, error) {
	serialOptions := make([]serial.Option, 0, 6)
	serialOptions = append(serialOptions, serial.WithBaudrate(mode.baudRate))
	serialOptions = append(serialOptions, serial.WithDataBits(mode.dataBits))
	serialOptions = append(serialOptions, serial.WithParity(serial.Parity(mode.parity)))
	serialOptions = append(serialOptions, serial.WithStopBits(serial.StopBits(mode.stopBits)))
	serialOptions = append(serialOptions, serial.WithWriteTimeout(50))
	serialOptions = append(serialOptions, serial.WithReadTimeout(0))

	sport, err := serial.Open(portName, serialOptions...)
	if err != nil {
		return nil, err
	}

	model := new(SerialPort)
	model.sport = sport
	return model, nil
}

func (this *SerialPort) Close() error {
	return this.sport.Close()
}

func (this *SerialPort) ResetBaudrate(baudRate int) error {
	return this.sport.Reconfigure(serial.WithBaudrate(baudRate))
}

func (this *SerialPort) ResetMode(mode Mode) error {
	serialOptions := make([]serial.Option, 0, 6)
	serialOptions = append(serialOptions, serial.WithBaudrate(mode.baudRate))
	serialOptions = append(serialOptions, serial.WithDataBits(mode.dataBits))
	serialOptions = append(serialOptions, serial.WithParity(serial.Parity(mode.parity)))
	serialOptions = append(serialOptions, serial.WithStopBits(serial.StopBits(mode.stopBits)))
	return this.sport.Reconfigure(serialOptions...)
}
