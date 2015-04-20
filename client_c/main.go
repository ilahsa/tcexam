package main

import (
	"encoding/json"
	"fmt"
	"github.com/funny/link"
	"github.com/ilahsa/tcexam/lib"
	"strconv"
)

// This is an echo client demo work with the echo_server.
// usage:
//     go run echo_client/main.go
func main() {
	link.DefaultProtocol = lib.TCProtocol
	client, err := link.Dial("tcp", "10.1.9.103:10010")
	if err != nil {
		panic(err)
	}
	acount := 0
	go client.Process(func(msg *link.InBuffer) error {
		println(string(msg.Data))
		//处理接收到的文件
		dat := map[string]string{}
		json.Unmarshal(msg.Data, &dat)
		action := dat["action"]
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
		client.Send(link.Bytes(by1))

		return nil
	})

	//登陆
	dat := map[string]string{
		"action": "cstart", "seq": "00001", "userid": "u_001", "password": "123456",
	}
	by, _ := json.Marshal(dat)
	client.Send(link.Bytes(by))
	for {

		//发10个获取
		for i := 0; i < 10; i++ {
			id := strconv.Itoa(i)
			dat := map[string]string{
				"action": "getfile", "seq": id,
			}
			by, _ := json.Marshal(dat)
			client.Send(link.Bytes(by))
		}
		fmt.Println("发送10个获取")
		var input string
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			break
		}
		client.Send(link.String(input))
	}

	client.Close()

	println("bye")
}