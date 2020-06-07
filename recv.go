package hserial

import (
	"fmt"
	"log"
	"time"
)

type CheckDataComplete func([]byte) bool

func (this *SerialPort) Recv(headTimeout, chunkTimeout int, checkFunc CheckDataComplete) ([]byte, ReadDurationPack,error) {
	pack := ReadDurationPack{}
	if this.sport == nil {
		return nil, pack,fmt.Errorf("serial port is nil")
	}

	var err error
	var readLen = 0
	var totalLen = 0
	var buff = make([]byte, 0, 1024)
	var temp = make([]byte, 256) //serialPort.Read(temp)读的是数组类型而不是切片
	var lastReadValidAt time.Time

	err = this.sport.SetReadTimeout(0)
	if err != nil {
		log.Println("serialPort SetReadTimeout err", err)
		return nil, pack,err
	}

	//先等待头字节
	startAt := time.Now()
	tryCount := 0
	for {
		readLen, err = this.sport.Read(temp)
		if readLen <= 0 {
			durMs := time.Since(startAt).Milliseconds()
			if  durMs < int64(headTimeout) {
				headWaitSleep(tryCount)
				continue
			}

			pack.HeadDuration = int(durMs)
			return nil, pack, nil
		}

		lastReadValidAt = time.Now()
		buff = append(buff, temp[:readLen]...)
		break
	}

	startReadContentAt := lastReadValidAt
	pack.HeadDuration = int(lastReadValidAt.Sub(startAt).Milliseconds())
	//fmt.Printf("head duration: %dms\r\n", pack.HeadDuration)

	if checkFunc != nil && checkFunc(buff) {
		return buff,pack, nil
	}

	//开始获取内容
	tryCount = 0
	for {
		intervalWaitSleep(tryCount)
		readLen, err = this.sport.Read(temp)
		if err != nil {
			log.Println("serialPort read err", err)
			pack.HeadDuration = int(time.Since(startReadContentAt).Milliseconds())
			return nil, pack,err
		}

		if readLen < 0 || readLen > 0xffff {
			log.Println("serial read length error,len=", readLen)
			pack.HeadDuration = int(time.Since(startReadContentAt).Milliseconds())
			return nil, pack,fmt.Errorf("serial read length error")
		}

		if readLen == 0 {
			if time.Since(lastReadValidAt).Milliseconds() >= int64(chunkTimeout) {
				pack.HeadDuration = int(time.Since(startReadContentAt).Milliseconds())
				return buff, pack,nil
			}
			continue
		}

		lastReadValidAt = time.Now()
		totalLen += readLen
		buff = append(buff, temp[:readLen]...)
		if checkFunc != nil && checkFunc(buff) {
			pack.HeadDuration = int(time.Since(startReadContentAt).Milliseconds())
			return buff, pack,nil
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
