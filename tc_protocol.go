package main

import (
	"bytes"
	"errors"
	"github.com/funny/link"
	"io"
	"strconv"
)

var (
	TCMaxPacketSize = 9999
	TCProtocol      = newTcProtocol()
)

// The packet spliting protocol like Erlang's {packet, N}.
// Each packet has a fix length packet header to present packet length.
type tcProtocol struct {
}

func newTcProtocol() *tcProtocol {
	return &tcProtocol{}
}

func (p *tcProtocol) New(v interface{}, _ link.ProtocolSide) (link.ProtocolState, error) {
	return p, nil
}

func (p *tcProtocol) PrepareOutBuffer(buffer *link.OutBuffer, size int) {
	buffer.Prepare(size)
}

func (p *tcProtocol) Write(writer io.Writer, packet *link.OutBuffer) error {
	if len(packet.Data) > TCMaxPacketSize {
		return link.PacketTooLargeError
	}
	lenStr := strconv.Itoa(len(packet.Data))
	writer.Write([]byte(lenStr))
	writer.Write([]byte{'\r'})
	writer.Write([]byte{'\n'})
	if _, err := writer.Write(packet.Data); err != nil {
		return err
	}
	return nil
}

func (p *tcProtocol) Read(reader io.Reader, buffer *link.InBuffer) error {
	// head
	lenBuf := bytes.Buffer{}

	tmpByte := make([]byte, 1)
	writePos := 0
	for {
		///包最大长度不能大于9999
		if writePos > 4 {
			return link.PacketTooLargeError
		}

		if _, err := io.ReadFull(reader, tmpByte); err != nil {
			return err
		}

		//fmt.Print(string(tmpByte))
		if tmpByte[0] == '\r' {
			if _, err := io.ReadFull(reader, tmpByte); tmpByte[0] != '\n' || err != nil {
				return errors.New("bad request")
			}
			len, err := strconv.Atoi(lenBuf.String())
			if err != nil {
				return errors.New("bad request,len is incorrect")
			}
			buffer.Prepare(len)
			if _, err := io.ReadFull(reader, buffer.Data); err != nil {
				return err
			} else {
				return nil
			}
		}
		lenBuf.WriteByte(tmpByte[0])
		writePos++
	}
	return nil
}
