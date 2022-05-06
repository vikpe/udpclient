package udpclient

import (
	"bytes"
	"errors"
	"net"
	"time"
)

type UdpClient interface {
	Request(address string, statusPacket []byte, expectedResponseHeader []byte) ([]byte, error)
	GetConfig() Config
}

type Config struct {
	BufferSize  uint16
	Retries     uint8
	TimeoutInMs uint16
}

type defaultClient struct {
	config Config
}

func New() *defaultClient {
	defaultConfig := Config{
		BufferSize:  8192,
		Retries:     3,
		TimeoutInMs: 500,
	}
	return NewWithConfig(defaultConfig)
}

func NewWithConfig(config Config) *defaultClient {
	return &defaultClient{config: config}
}

func (c defaultClient) GetConfig() Config {
	return c.config
}

func (c defaultClient) Request(address string, statusPacket []byte, expectedResponseHeader []byte) ([]byte, error) {
	conn, err := net.Dial("udp4", address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	config := c.GetConfig()
	responseBuffer := make([]byte, config.BufferSize)
	responseLength := 0

	for i := uint8(0); i < config.Retries; i++ {
		conn.SetDeadline(getDeadline(config.TimeoutInMs))

		_, err = conn.Write(statusPacket)
		if err != nil {
			return nil, err
		}

		conn.SetDeadline(getDeadline(config.TimeoutInMs))
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
	headerLength := len(expectedResponseHeader)
	header := response[0:headerLength]

	isValidResponseHeader := bytes.Equal(header, expectedResponseHeader)
	if !isValidResponseHeader {
		err = errors.New(address + ": Invalid response header.")
		return nil, err
	}

	return response[headerLength:], nil
}

func getDeadline(timeoutInMs uint16) time.Time {
	return time.Now().Add(time.Duration(timeoutInMs) * time.Millisecond)
}
