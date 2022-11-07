package udpclient_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/udpclient"
	"github.com/vikpe/udphelper"
)

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
		response, err := client.SendPacket(":8001", []byte("ping"))

		assert.Equal(t, []byte{}, response)
		assert.ErrorContains(t, err, "write udp4")
		assert.ErrorContains(t, err, ":8001: i/o timeout")
	})

	t.Run("Read timeout", func(t *testing.T) {
		addr := ":8002"

		go func() {
			udphelper.New(addr).Listen()
		}()
		time.Sleep(10 * time.Millisecond)

		client := udpclient.New()
		client.Config.TimeoutInMs = 20
		response, err := client.SendPacket(":8002", []byte("ping"))

		assert.Equal(t, []byte{}, response)
		assert.ErrorContains(t, err, "read udp4")
		assert.ErrorContains(t, err, ":8002: i/o timeout")
	})

	t.Run("Success", func(t *testing.T) {
		addr := ":8003"

		go func() {
			udphelper.New(addr).Echo()
		}()
		time.Sleep(10 * time.Millisecond)

		client := udpclient.New()
		response, err := client.SendPacket(addr, []byte("ping"))

		assert.Equal(t, "ok:ping", string(response))
		assert.Equal(t, nil, err)
	})
}

func TestClient_SendCommand(t *testing.T) {
	helloWorldCommand := udpclient.Command{
		RequestPacket:  []byte("HELLO WORLD"),
		ResponseHeader: []byte("hello ok:"),
	}

	t.Run("Unknown host", func(t *testing.T) {
		response, err := udpclient.New().SendCommand("foo:666", helloWorldCommand)
		assert.Equal(t, []byte{}, response)
		assert.ErrorContains(t, err, "dial udp4: lookup foo:")
	})

	t.Run("Invalid repsonse header", func(t *testing.T) {
		t.Run("shorter than expected", func(t *testing.T) {
			address := fmt.Sprintf(":%d", 9000)
			udpServer := udphelper.New(address)

			go func() {
				udpServer.Respond([]byte("hello"))
			}()
			time.Sleep(10 * time.Millisecond)

			client := udpclient.New()
			response, err := client.SendCommand(address, helloWorldCommand)
			assert.Equal(t, helloWorldCommand.RequestPacket, udpServer.Requests[0])
			assert.Equal(t, []byte{}, response)
			assert.ErrorContains(t, err, `:9000: Invalid response header, expected "hello ok:"`)
		})

		t.Run("not equal to expected", func(t *testing.T) {
			address := fmt.Sprintf(":%d", 9001)
			udpServer := udphelper.New(address)

			go func() {
				udpServer.Respond([]byte("hello fail:"))
			}()
			time.Sleep(10 * time.Millisecond)

			client := udpclient.New()
			response, err := client.SendCommand(address, helloWorldCommand)
			assert.Equal(t, helloWorldCommand.RequestPacket, udpServer.Requests[0])
			assert.Equal(t, []byte{}, response)
			assert.ErrorContains(t, err, `:9001: Invalid response header, expected "hello ok:"`)
		})
	})

	t.Run("Valid response header", func(t *testing.T) {
		address := fmt.Sprintf(":%d", 9002)
		udpServer := udphelper.New(address)

		go func() {
			udpServer.Respond([]byte("hello ok:HELLO WORLD"))
		}()
		time.Sleep(10 * time.Millisecond)

		client := udpclient.New()
		response, err := client.SendCommand(address, helloWorldCommand)
		assert.Equal(t, []byte("HELLO WORLD"), udpServer.Requests[0])
		assert.Equal(t, []byte("HELLO WORLD"), response)
		assert.Nil(t, err)
	})
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
