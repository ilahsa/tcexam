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
	client, err := link.Dial("tcp", "85.195.110.142:10010")
	if err != nil {
		panic(err)
	}
	acount := 0
	go client.Process(func(msg *link.InBuffer) error {
		fmt.Println("接收 ", string(msg.Data))
		//处理接收到的文件
		dat := map[string]string{}
		json.Unmarshal(msg.Data, &dat)
		action := dat["action"]
		if action == "res_putfile" {
			id := dat["id"]

			acount = acount + 1
			fmt.Println("收到答案的总数", acount)

			dat1 := map[string]string{
				"action": "reportanswer", "seq": "_seq_0001", "id": id, "result": "0",
			}
			by, _ := json.Marshal(dat1)
			client.Send(link.Bytes(by))
			fmt.Println("发送reportanswer", string(by))
		}

		return nil
	})

	for {
		//发10个文件
		for i := 0; i < 10; i++ {
			id := strconv.Itoa(i)
			dat := map[string]string{
				"action": "putfile", "file": "file_" + id, "seq": id,
			}
			by, _ := json.Marshal(dat)
			client.Send(link.Bytes(by))
			fmt.Println("发送", string(by))
		}
		fmt.Println("发送10个文件")
		var input string
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			break
		}
		client.Send(link.String(input))
	}

	client.Close()

	println("bye")
}
