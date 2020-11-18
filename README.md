# dnacenter
This is a work in progress

## Installation

Install via **go get**:

```shell
go get -u github.com/ApogeeNetworking/dnacenter
```

## Usage
Basic usage can be found below

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    dnac "github.com/ApogeeNetworking/dnacenter"
)

func main() {
    // Used in Development for SelfSigned Certs
    ignoreSSL := true
    dna := dnac.NewClient("host/ip", "user", "pass", ignoreSSL)

    err := dna.Login()
    if err != nil {
        log.Fatalf("%v", err)
    }
    pnpDevice, err := dna.PnP.GetDevice("id")
    if err != nil {
        // Do something about the Error
    }
}
```
