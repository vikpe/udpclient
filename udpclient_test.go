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

func udpListenAndEcho(addr string) {
	conn, _ := net.ListenPacket("udp", addr)
	buffer := make([]byte, 1024)
	messageLength, dst, _ := conn.ReadFrom(buffer)
	message := buffer[:messageLength]
	response := "OK: " + string(message)
	conn.WriteTo([]byte(response), dst)
}

func udpListen(addr string) {
	conn, _ := net.ListenPacket("udp", addr)
	conn.ReadFrom(make([]byte, 1024))
}

func TestClient_SendPacket(t *testing.T) {
	t.Run("Unknown host", func(t *testing.T) {
		client := udpclient.New()
		response, err := client.SendPacket("foo:666", []byte{1, 2, 3})
		assert.Equal(t, []byte{}, response)
		assert.ErrorContains(t, err, "dial udp4: lookup foo:")
	})

	t.Run("Write timeout", func(t *testing.T) {
		client := udpclient.New()
		client.Config.TimeoutInMs = 0
		response, err := client.SendPacket(":8001", []byte("PING"))

		assert.Equal(t, []byte{}, response)
		assert.ErrorContains(t, err, "write udp4")
		assert.ErrorContains(t, err, ":8001: i/o timeout")
	})

	t.Run("Read timeout", func(t *testing.T) {
		addr := ":8002"

		go func() {
			udpListen(addr)
		}()

		time.Sleep(10 * time.Millisecond)

		client := udpclient.New()
		client.Config.TimeoutInMs = 20
		response, err := client.SendPacket(":8002", []byte("PING"))

		assert.Equal(t, []byte{}, response)
		assert.ErrorContains(t, err, "read udp4")
		assert.ErrorContains(t, err, ":8002: i/o timeout")
	})

	t.Run("Success", func(t *testing.T) {
		addr := ":8003"

		go func() {
			udpListenAndEcho(addr)
		}()

		time.Sleep(10 * time.Millisecond)

		client := udpclient.New()
		response, err := client.SendPacket(addr, []byte("PING"))

		assert.Equal(t, "OK: PING", string(response))
		assert.Equal(t, nil, err)
	})
}

func TestClient_SendCommand(t *testing.T) {
	t.Run("Unknown host", func(t *testing.T) {
		client := udpclient.New()
		command := udpclient.Command{
			RequestPacket:  []byte("HELLO WORLD"),
			ResponseHeader: []byte("OK: "),
		}
		response, err := client.SendCommand("foo:666", command)
		assert.Equal(t, []byte{}, response)
		assert.ErrorContains(t, err, "dial udp4: lookup foo:")
	})

	testCases := []struct {
		testName             string
		port                 int
		command              udpclient.Command
		expectedResponseBody []byte
		expectedError        error
	}{
		{
			"Valid response header",
			9000,
			udpclient.Command{
				RequestPacket:  []byte("HELLO WORLD"),
				ResponseHeader: []byte("OK: "),
			},
			[]byte("HELLO WORLD"),
			nil,
		},
		{
			"Invalid response header",
			9001,
			udpclient.Command{
				RequestPacket:  []byte("HELLO WORLD"),
				ResponseHeader: []byte("NOT OK: "),
			},
			[]byte{},
			errors.New(":9001: Invalid response header"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			udpAddress := fmt.Sprintf(":%d", tc.port)

			go func() {
				udpListenAndEcho(udpAddress)
			}()

			time.Sleep(10 * time.Millisecond)

			client := udpclient.New()
			response, err := client.SendCommand(udpAddress, tc.command)
			assert.Equal(t, tc.expectedResponseBody, response)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func ExampleNewWithConfig() {
	config := udpclient.Config{
		BufferSize:  256,
		Retries:     0,
		TimeoutInMs: 800,
	}
	client := udpclient.NewWithConfig(config)

	fmt.Print(client.Config)
}

func ExampleClient_SendPacket() {
	client := udpclient.New()
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}

	response, err := client.SendPacket("qw.foppa.dk:27502", statusPacket)
	fmt.Println(response, err)
}

func ExampleClient_SendCommand() {
	client := udpclient.New()
	command := udpclient.Command{
		RequestPacket:  []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a},
		ResponseHeader: []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'},
	}

	response, err := client.SendCommand("qw.foppa.dk:27502", command)
	fmt.Println(response, err)
}
