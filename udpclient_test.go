package udpclient_test

import (
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/udpclient"
)

func captureUdpPacket(addr string) {
	conn, _ := net.ListenPacket("udp", addr)
	buffer := make([]byte, 1024)
	messageLength, dst, _ := conn.ReadFrom(buffer)
	message := buffer[:messageLength]
	response := "OK: " + string(message)
	conn.WriteTo([]byte(response), dst)
}

func TestClient_SendPacket(t *testing.T) {
	addr := ":8000"

	go func() {
		captureUdpPacket(addr)
	}()

	time.Sleep(10 * time.Millisecond)

	client := udpclient.New()
	response, err := client.SendPacket(addr, []byte("PING"))

	assert.Equal(t, "OK: PING", string(response))
	assert.Equal(t, nil, err)
}

func TestClient_SendCommand(t *testing.T) {
	testCases := []struct {
		port                 int
		command              udpclient.Command
		expectedResponseBody []byte
		expectedError        error
	}{
		{
			9000,
			udpclient.Command{
				RequestPacket:  []byte("HELLO WORLD"),
				ResponseHeader: []byte("OK: "),
			},
			[]byte("HELLO WORLD"),
			nil,
		},
		{
			9001,
			udpclient.Command{
				RequestPacket:  []byte("HELLO WORLD"),
				ResponseHeader: []byte("NOT OK: "),
			},
			nil,
			errors.New(":9001: Invalid response header"),
		},
	}

	for _, tc := range testCases {
		udpAddress := fmt.Sprintf(":%d", tc.port)

		go func() {
			captureUdpPacket(udpAddress)
		}()

		time.Sleep(10 * time.Millisecond)

		client := udpclient.New()
		response, err := client.SendCommand(udpAddress, tc.command)

		assert.Equal(t, tc.expectedResponseBody, response)
		assert.Equal(t, tc.expectedError, err)
	}
}

func ExampleNewWithConfig() {
	config := udpclient.Config{
		BufferSize:  256,
		Retries:     0,
		TimeoutInMs: 800,
	}
	client := udpclient.NewWithConfig(config)

	fmt.Print(client.GetConfig())
}

func ExampleClient_SendPacket() {
	client := udpclient.New()
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}

	response, err := client.SendPacket(
		"qw.foppa.dk:27502",
		statusPacket,
	)

	fmt.Println(response, err)
}
