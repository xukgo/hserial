package hserial

import (
	"fmt"
	"log"
	"time"
)

type CheckDataComplete func([]byte) bool

func (this *SerialPort) Recv(headTimeout, chunkTimeout int, checkFunc CheckDataComplete) ([]byte, error) {
	if this.sport == nil {
		return nil, fmt.Errorf("serial port is nil")
	}

	var err error
	var readLen = 0
	var totalLen = 0
	var buff = make([]byte, 0, 1024)
	var temp = make([]byte, 256) //serialPort.Read(temp)读的是数组类型而不是切片

	err = this.sport.SetReadTimeout(0)
	if err != nil {
		fmt.Println("serialPort SetReadTimeout err", err)
		return nil, err
	}

	//先等待头字节
	startAt := time.Now()
	tryCount := 0
	for {
		headWaitSleep(tryCount)
		readLen, err = this.sport.Read(temp)
		if readLen <= 0 {
			if time.Since(startAt).Milliseconds() < int64(headTimeout) {
				continue
			}

			return nil, nil
		}
		buff = append(buff, temp[:readLen]...)
		break
	}
	fmt.Printf("head duration: %dms\r\n", time.Since(startAt).Milliseconds())
	if checkFunc != nil && checkFunc(buff) {
		return buff, nil
	}

	//开始获取内容
	var lastReadValidAt = time.Now()
	tryCount = 0
	for {
		intervalWaitSleep(tryCount)
		readLen, err = this.sport.Read(temp)
		if err != nil {
			fmt.Println("serialPort read err", err)
			return nil, err
		}

		if readLen < 0 || readLen > 0xffff {
			log.Println("serial read length error,len=", readLen)
			return nil, fmt.Errorf("serial read length error")
		}

		if readLen == 0 {
			if time.Since(lastReadValidAt).Milliseconds() >= int64(chunkTimeout) {
				return buff, nil
			}
			continue
		}

		lastReadValidAt = time.Now()
		totalLen += readLen
		buff = append(buff, temp[:readLen]...)
		if checkFunc != nil && checkFunc(buff) {
			return buff, nil
		}
		//fmt.Printf("now recv total=> %s\r\n", buff)
	}
}

func intervalWaitSleep(tryCount int) {
	//time.Sleep(time.Millisecond * 2)
	//return
	n := 1
	if tryCount == 0 {
		n = 2
	} else if tryCount == 1 {
		n = 4
	} else if tryCount <= 3 {
		n = 8
	} else {
		n = 12
	}
	time.Sleep(time.Millisecond * time.Duration(n))
}

func headWaitSleep(tryCount int) {
	//time.Sleep(time.Millisecond * 2)
	//return
	n := 1
	if tryCount == 0 {
		n = 1
	} else if tryCount == 1 {
		n = 2
	} else if tryCount <= 3 {
		n = 4
	} else {
		n = 8
	}
	time.Sleep(time.Millisecond * time.Duration(n))
}

func getMax(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
