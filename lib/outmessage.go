package lib

import (
	"github.com/funny/binary"
)

type Message []byte

func (msg Message) Send(conn *binary.Writer) error {
	conn.WritePacket(msg, binary.SplitByUint16BE)
	return nil
}

func (msg *Message) Receive(conn *binary.Reader) error {
	*msg = conn.ReadPacket(binary.SplitByUint16BE)
	return nil
}

func MakeMsg(by []byte) Message {
	var msg Message
	msg = by
	return msg
}
