package udpclient_test

import "github.com/vikpe/udpclient"

func ExampleNewWithConfig() {
	config := udpclient.Config{
		BufferSize:  256,
		Retries:     0,
		TimeoutInMs: 800,
	}
	udpClient := udpclient.NewWithConfig(config)
}

func ExampleClient_Request() {
	udpClient := udpclient.New()
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}
	expectedHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'}

	response, err := udpClient.Request(
		"qw.foppa.dk:27502",
		statusPacket,
		expectedHeader,
	)
}
