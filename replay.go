package main

import (
	"net"
	"flag"
	"log"
	"io/ioutil"
	"time"
	"sync/atomic"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
)

var (
	pkgNum int
	pkgNumPerPeriod int
	remoteAddr string
	pkgFile string
	pkgToSend []byte
	seqOff int
	seq uint32
	period int
)

func init() {
	flag.IntVar(&pkgNum, "c", 10, "pkg count to send")//发包总数
	flag.StringVar(&remoteAddr, "r", "183.60.48.140:8000", "remote addr to send pkg to")//发包目的地址
	flag.StringVar(&pkgFile, "f", "req.bin", "pkg file to send")//发包内容，建议使用抓包工具直接将包抓下来
	flag.IntVar(&pkgNumPerPeriod, "p", 10000, "pkg cnt per period")//每个周期发多少个包
	flag.IntVar(&seqOff, "o", 2000, "seq offset")// TODO seq
	flag.IntVar(&period, "t", 100, "period, in ms")//周期间隔时间
	flag.Parse()

	var err error
	if pkgToSend, err = ioutil.ReadFile(pkgFile); err != nil {
		log.Fatal(err)
	}

	seq = uint32(seqOff)
}

func nextSeq() uint16 {
    var newSeq uint32
    if seq >= (1<<16 - 1) {
        newSeq = 0
        atomic.StoreUint32(&seq, 0)
    } else {
        newSeq = atomic.AddUint32(&seq, 1)//原子操作
    }

    return uint16(newSeq)
}

func changeSeq(buf []byte) {
    if len(buf) < 12 {
        panic("buf is too small")
    }

    var bk [4]byte
    tmp := bk[0:4]

    binary.BigEndian.PutUint16(tmp, nextSeq())//大端 ＝＝ 网络序
    for i := 0; i < 2; i++ {
        buf[1+2+2+i] = tmp[i]
    }
}

func main() {
	ra, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("remote addr: %+v\n", ra)

	conn, err := net.DialUDP("udp", nil, ra)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	sendCnt := 0
	allCnt := 0

	sendPkg := func () {
		//changeSeq(pkgToSend) // 如果你需要改变序列号则重写该函数
		conn.Write(pkgToSend)
		sendCnt += 1
		allCnt += 1
	}

	var failSec uint64

	begin :=time.Now()
	start:= time.Now()
	for allCnt < pkgNum {
		elapsed := time.Since(start) / (time.Millisecond * time.Duration(period))
		if 0 == elapsed {
			if sendCnt < pkgNumPerPeriod {
				sendPkg()
			}
		} else {
			if sendCnt < pkgNumPerPeriod {
				failSec += 1
			}
			sendCnt = 0
			start = time.Now()
			sendPkg()
		}
	}

	report := func () {
		fmt.Printf("%d pkg, takes %v\n", pkgNum, time.Since(begin))
		fmt.Printf("fail sec: %d\n", failSec)
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c) // 收到信号则提前退出
		<-c

		fmt.Println("signaled")

		report()
		os.Exit(-1)
	}()

	report()
}
