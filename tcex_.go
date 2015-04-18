package main

//import (
//	"encoding/json"
//	"flag"
//	"fmt"
//	"github.com/funny/link"
//	"strconv"
//	"sync/atomic"
//	"testing"
//)

//var (
//	seq int32
//)

//func Test_ex(t *testing.T) {
//	//开启服务端
//	go server_start()
//	//开启生产者
//	go p()

//	//开启3个消费者
//	//go c()

//	//go c()

//	//go c()
//	fmt.Println("test")
//}

////开启服务端
//func server_start() {
//	flag.Parse()

//	link.DefaultConnBufferSize = *buffersize

//	server, err := link.Listen("tcp", "127.0.0.1:10010")
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("1")
//	ULogger.Info("server start:", server.Listener().Addr().String())

//	server.Serve(func(session *link.Session) {
//		ULogger.Info("client", session.Conn().RemoteAddr().String(), "in")

//		session.Process(func(msg *link.InBuffer) error {
//			ULogger.Info("client", session.Conn().RemoteAddr().String(), "say:", string(msg.Data))
//			var m map[string]interface{}

//			var retByte []byte
//			err := json.Unmarshal(msg.Data, &m)
//			if err == nil {
//				ULogger.Errorf("bad request,req is %s\n", string(msg.Data))
//				retByte = nil
//			} else {
//				retByte = Process(session, m)
//			}
//			if retByte == nil {
//				return nil
//			} else {

//				return session.Send(link.Bytes(retByte))
//			}
//		})

//		ULogger.Info("client", session.Conn().RemoteAddr().String(), "close")
//	})
//}

/////生产者
//func p() {
//	fmt.Println("2")
//	ch := make(chan bool, 10)
//	client, err := link.Dial("tcp", "127.0.0.1:10010")
//	if err != nil {
//		panic(err)
//	}
//	go client.Process(func(msg *link.InBuffer) error {
//		println(string(msg.Data))
//		return nil
//	})

//	for {
//		ch <- true
//		q := getId1()
//		ret := map[string]string{"action": "putfile", "seq": q, "fileid": "fid_" + q, "file": "f_" + q}
//		by, _ := json.Marshal(ret)
//		//发送图片
//		client.Send(link.Bytes(by))
//	}

//	client.Close()

//	println("bye")
//}

////消费者
//func c() {

//}
//func getId1() string {
//	return "aaaa"
//	seq := atomic.AddInt32(&seq, 1)
//	return strconv.Itoa(int(seq))
//}
