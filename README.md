# dnacenter
This is a work in progress

## Installation

Install via **go get**:

```shell
go get -u github.com/drkchiloll/dnacenter
```

## Usage
Basic usage can be found below

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    dnac "github.com/drkchiloll/dnacenter"
)

func main() {
    dna := dnac.NewClient("host/ip", "user", "pass", true)

    err := dna.Login()
    if err != nil {
        log.Fatalf("%v", err)
    }
    devices, err := dna.GetNetDevice()
    if err != nil {
        log.Fatalf("%v", err)
    }
    fmt.Println(device)
}
```