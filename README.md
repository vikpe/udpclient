# UDP Client [![Go Reference](https://pkg.go.dev/badge/github.com/vikpe/udpclient.svg)](https://pkg.go.dev/github.com/vikpe/udpclient)
> UDP client for Go


## Example
Send a `status 23` command to a QuakeWorld server.

```go
package main

import (
	"github.com/vikpe/udpclient"
)

func main() {
	client := udpclient.New()
	statusPacket := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}
	expectedResponseHeader := []byte{0xff, 0xff, 0xff, 0xff, 'n', '\\'}

	response, err := client.Request(
		"qw.foppa.dk:27502",
		statusPacket,
		expectedResponseHeader,
	)
}
```
