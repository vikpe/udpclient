package udpclient

import (
	"bytes"
	"errors"
	"net"
	"time"
)

type Command struct {
	RequestPacket  []byte
	ResponseHeader []byte
}

type Config struct {
	BufferSize  uint16
	Retries     uint8
	TimeoutInMs uint16
}

type Client struct {
	Config Config
}

func New() *Client {
	var config = Config{
		BufferSize:  8192,
		Retries:     3,
		TimeoutInMs: 800,
	}

	return NewWithConfig(config)
}

func NewWithConfig(config Config) *Client {
	return &Client{Config: config}
}

func (c Client) SendCommand(address string, command Command) ([]byte, error) {
	response, err := c.SendPacket(address, command.RequestPacket)

	if err != nil {
		return []byte{}, err
	}

	headerLength := len(command.ResponseHeader)
	header := response[0:headerLength]

	isValidResponseHeader := bytes.Equal(header, command.ResponseHeader)
	if !isValidResponseHeader {
		err = errors.New(address + ": Invalid response header")
		return []byte{}, err
	}

	responseBody := response[headerLength:]

	return responseBody, nil
}

func (c Client) SendPacket(address string, packet []byte) ([]byte, error) {
	conn, err := net.Dial("udp4", address)
	if err != nil {
		return []byte{}, err
	}
	defer conn.Close()

	response := bytes.NewBuffer(make([]byte, 0))
	frameBuffer := make([]byte, c.Config.BufferSize)
	frameLength := 0
	shouldRetry := true

	for i := uint8(0); i < c.Config.Retries; i++ {
		conn.SetWriteDeadline(getDeadline(c.Config.TimeoutInMs))

		_, err = conn.Write(packet)
		if err != nil {
			return []byte{}, err
		}

		for {
			conn.SetReadDeadline(getDeadline(c.Config.TimeoutInMs))
			frameLength, err = conn.Read(frameBuffer)

			if err != nil { // udp error or end of response
				if response.Len() > 0 {
					err = nil
					shouldRetry = false
				}
				break
			} else { // successfully read frame
				response.Write(frameBuffer[:frameLength])
			}
		}

		if !shouldRetry {
			break
		}
	}

	return response.Bytes(), err
}

func getDeadline(timeoutInMs uint16) time.Time {
	return time.Now().Add(time.Duration(timeoutInMs) * time.Millisecond)
}
