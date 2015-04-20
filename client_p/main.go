package main

import (
	"fmt"
	"github.com/funny/link"
	"github.com/ilahsa/tcexam/lib"
)

// This is an echo client demo work with the echo_server.
// usage:
//     go run echo_client/main.go
func main() {
	link.DefaultProtocol = lib.TCProtocol
	client, err := link.Dial("tcp", "127.0.0.1:10010")
	if err != nil {
		panic(err)
	}
	go client.Process(func(msg *link.InBuffer) error {
		println(string(msg.Data))
		return nil
	})

	for {
		var input string
		if _, err := fmt.Scanf("%s\n", &input); err != nil {
			break
		}
		client.Send(link.String(input))
	}

	client.Close()

	println("bye")
}
