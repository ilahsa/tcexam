package main

import (
	"encoding/json"
	"fmt"
	"github.com/funny/link"
	"github.com/ilahsa/tcexam/lib"
	//"strconv"
	"sync"
)

// This is an echo client demo work with the echo_server.
// usage:
//     go run echo_client/main.go
func main() {

	link.DefaultProtocol = lib.TCProtocol
	client, err := link.Dial("tcp", "10.1.9.27:10010")
	if err != nil {
		panic(err)
	}
	stopWait := sync.WaitGroup{}
	acount := 0
	go client.Process(func(msg *link.InBuffer) error {
		fmt.Println("接收", string(msg.Data))
		//处理接收到的文件r
		dat := map[string]string{}
		json.Unmarshal(msg.Data, &dat)
		action := dat["action"]
		ret := dat["result"]
		if action == "res_cstart" && ret == "1" {
			stopWait.Done()
			fmt.Println("登陆成功")
			return nil
		}
		if action == "res_getfile" {
			acount = acount + 1
			fmt.Println("收到答案的总数", acount)
		}
		fid := dat["id"]
		seq := dat["seq"]
		dat1 := map[string]string{
			"action": "answer", "seq": seq, "id": fid,
		}
		by1, _ := json.Marshal(dat1)
		_ = by1
		//client.Send(link.Bytes(by1))

		return nil
	})

	//登陆
	dat := map[string]string{
		"action": "cstart", "seq": "00001", "userid": "u_001", "password": "123456",
	}
	by, _ := json.Marshal(dat)
	client.Send(link.Bytes(by))
	fmt.Println("发送登陆", string(by))
	stopWait.Add(1)
	for {
		//等待登陆成功
		stopWait.Wait()

		var input string
		if _, err := fmt.Scanf("%s\n", &input); err == nil {
			//fmt.Println(input + "ssss")
		} else {
			break
		}

		dat := map[string]string{
			//"action": "putfile", "file": "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", "seq": id,
			"action": "conrol", "event": input, "seq": "_tmp_0001",
		}
		by, _ := json.Marshal(dat)
		client.Send(link.Bytes(by))

		client.Send(link.String(input))
	}

	client.Close()

	println("bye")
}
