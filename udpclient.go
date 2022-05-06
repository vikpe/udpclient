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
	config Config
}

func New() *Client {
	var config = Config{
		BufferSize:  8192,
		Retries:     3,
		TimeoutInMs: 500,
	}

	return NewWithConfig(config)
}

func NewWithConfig(config Config) *Client {
	return &Client{config: config}
}

func (c Client) GetConfig() Config {
	return c.config
}

func (c Client) SendCommand(address string, command Command) ([]byte, error) {
	response, err := c.SendPacket(address, command.RequestPacket)

	if err != nil {
		return nil, err
	}

	headerLength := len(command.ResponseHeader)
	header := response[0:headerLength]

	isValidResponseHeader := bytes.Equal(header, command.ResponseHeader)
	if !isValidResponseHeader {
		err = errors.New(address + ": Invalid response header.")
		return nil, err
	}

	responseBody := response[headerLength:]

	return responseBody, nil
}

func (c Client) SendPacket(address string, packet []byte) ([]byte, error) {
	conn, err := net.Dial("udp4", address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	responseBuffer := make([]byte, c.config.BufferSize)
	responseLength := 0

	for i := uint8(0); i < c.config.Retries; i++ {
		conn.SetDeadline(getDeadline(c.config.TimeoutInMs))

		_, err = conn.Write(packet)
		if err != nil {
			return nil, err
		}

		conn.SetDeadline(getDeadline(c.config.TimeoutInMs))
		responseLength, err = conn.Read(responseBuffer)
		if err != nil {
			continue
		}

		break
	}

	if err != nil {
		return nil, err
	}

	return responseBuffer[:responseLength], nil
}

func getDeadline(timeoutInMs uint16) time.Time {
	return time.Now().Add(time.Duration(timeoutInMs) * time.Millisecond)
}
