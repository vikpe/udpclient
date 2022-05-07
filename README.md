# UDP Client [![Go Reference](https://pkg.go.dev/badge/github.com/vikpe/udpclient.svg)](https://pkg.go.dev/github.com/vikpe/udpclient) [![Test](https://github.com/vikpe/udpclient/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/vikpe/udpclient/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/vikpe/udpclient/branch/main/graph/badge.svg)](https://codecov.io/gh/vikpe/udpclient)

> UDP client for Go

## Example

Send a `status 23` packet to a QuakeWorld server.

```go
package main

import (
	"fmt"

	"github.com/vikpe/udpclient"
)

func main() {
	client := udpclient.New()
	packet := []byte{0xff, 0xff, 0xff, 0xff, 's', 't', 'a', 't', 'u', 's', ' ', '2', '3', 0x0a}
	address := "qw.foppa.dk:27502"
	response, err := client.SendPacket(address, packet)

	fmt.Println(response, err)
}
```
