# latitudesh-go
[![GoDoc](https://godoc.org/github.com/latitudesh/latitudesh-go?status.svg)](https://godoc.org/github.com/latitudesh/latitudesh-go)

latitude-shgo is a Go client library for accessing the Latitude.sh API.

You can view the API docs here: https://docs.latitude.sh/reference


## Install
```sh
go get github.com/latitudesh/latitudesh-go@vX.Y.Z
```

where X.Y.Z is the [version](https://github.com/latitudesh/latitudesh-go/releases) you need.

## Usage

```go
package main

import (
    latitude "github.com/latitudesh/latitudesh-go"
)

func main() {
    client := latitude.NewClientWithAuth("Latitude.sh", apiToken, nil)
}
```

