package hserial

import (
	"fmt"
	"time"
)

func (this *SerialPort) Send(buff []byte) (int, error) {
	if this.sport == nil {
		return 0, fmt.Errorf("serial port is nil")
	}

	n, err := this.sport.Write(buff)
	return n, err
}

func (this *SerialPort) SendBatch(buff []byte, onceMax int, inteval int) (int, error) {
	if this.sport == nil {
		return 0, fmt.Errorf("serial port is nil")
	}

	if onceMax >= len(buff) {
		n, err := this.sport.Write(buff)
		return n, err
	}

	inSendLen := len(buff)
	totalSendLen := 0
	data := make([]byte, 0, onceMax)

	for {
		data, copyLen := copyBytes(data, buff, onceMax)
		buff = buff[copyLen:]

		n, err := this.sport.Write(data)
		totalSendLen += n
		if err != nil {
			return totalSendLen, err
		}
		if totalSendLen >= inSendLen {
			return totalSendLen, nil
		}

		data = data[0:]
		time.Sleep(time.Millisecond * time.Duration(inteval))
	}
}

func copyBytes(slice []byte, src []byte, maxCount int) ([]byte, int) {
	n := 0
	if maxCount > len(src) {
		n = len(src)
	} else {
		n = maxCount
	}
	slice = append(slice, src[:n]...)
	return slice, n
}
