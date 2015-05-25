package lib

import (
	"errors"
	//"fmt"
	"github.com/funny/link"
	"io"
	"strconv"
)

var Protocol = protocol{}

//var _ link.Codec = protocol{}

func New(delim byte) link.Protocol {
	return protocol{}
}

type protocol struct {
}

func (protocol protocol) NewCodec() link.Codec {
	return protocol
}

func (codec protocol) makeBuffer(buf *link.Buffer, msg link.Message) error {
	// prepend packet buffer
	size := 1024

	if sizeable, ok := msg.(link.Sizeable); ok {
		size = sizeable.BufferSize()
	}

	buf.Reset(0, size)

	// write pakcet content
	if err := msg.WriteBuffer(buf); err != nil {
		return err
	}
	by1 := buf.Data

	lenTmp := strconv.Itoa(len(by1))
	//	fmt.Println(lenTmp + "____")
	buf.Reset(0, len(by1)+len(lenTmp)+2)
	tmpByte := make([]byte, 0, len(by1)+len(lenTmp)+2)
	tmpByte = append(tmpByte, []byte(lenTmp)...)
	tmpByte = append(tmpByte, '\r', '\n')
	tmpByte = append(tmpByte, by1...)

	buf.WriteBytes(tmpByte)
	return nil
}

func (codec protocol) MakeBroadcast(buf *link.Buffer, msg link.Message) error {
	return errors.New("not implement")
}

func (codec protocol) SendBroadcast(conn *link.Conn, buf *link.Buffer) error {
	return errors.New("not implement")
}

func (codec protocol) SendMessage(conn *link.Conn, buf *link.Buffer, msg link.Message) error {
	err := codec.makeBuffer(buf, msg)
	if err != nil {
		return err
	}
	_, err = conn.Write(buf.Data)
	//	fmt.Println(222222222222222)
	//	fmt.Println(string(buf.Data))
	//	fmt.Println(333333333333333)
	return err
}

func (codec protocol) ProcessRequest(conn *link.Conn, buf *link.Buffer, handler link.RequestHandler) error {
	data, _, err := conn.Reader.ReadLine()
	if err != nil {
		return err
	}

	size, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}

	buf.Reset(size, size)
	if _, err := io.ReadFull(conn, buf.Data); err != nil {
		//		fmt.Println("error 1")
		return err
	}
	//	fmt.Println(string(buf.Data))
	return handler(buf)
}
