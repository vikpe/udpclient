package udpclient_test

import (
	"fmt"

	"github.com/vikpe/udpclient"
)

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
