package udpclient

import (
	"bytes"
	"errors"
	"net"
	"time"
)

type Config struct {
	BufferSize  uint16
	Retries     uint8
	TimeoutInMs uint16
}

type Client struct {
	Config Config
}

func New() *Client {
	defaultConfig := Config{
		BufferSize:  8192,
		Retries:     3,
		TimeoutInMs: 500,
	}
	return NewWithConfig(defaultConfig)
}

func NewWithConfig(config Config) *Client {
	return &Client{Config: config}
}

func (client Client) Request(address string, statusPacket []byte, expectedResponseHeader []byte) ([]byte, error) {
	conn, err := net.Dial("udp4", address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	responseBuffer := make([]byte, client.Config.BufferSize)
	responseLength := 0

	for i := uint8(0); i < client.Config.Retries; i++ {
		conn.SetDeadline(client.getDeadline())

		_, err = conn.Write(statusPacket)
		if err != nil {
			return nil, err
		}

		conn.SetDeadline(client.getDeadline())
		responseLength, err = conn.Read(responseBuffer)
		if err != nil {
			continue
		}

		break
	}

	if err != nil {
		return nil, err
	}

	response := responseBuffer[:responseLength]

	isValidResponseHeader := bytes.Equal(response[:len(expectedResponseHeader)], expectedResponseHeader)
	if !isValidResponseHeader {
		err = errors.New(address + ": Invalid response header.")
		return nil, err
	}

	return response, nil
}

func (client Client) getDeadline() time.Time {
	return time.Now().Add(time.Duration(client.Config.TimeoutInMs) * time.Millisecond)
}
