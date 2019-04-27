AvatarMe
========

[![GoDoc](https://godoc.org/github.com/SaraTrawnik/avatarme?status.svg)](https://godoc.org/github.com/SaraTrawnik/avatarme)
[![GoReportCard](https://goreportcard.com/badge/SaraTrawnik/avatarme)](https://goreportcard.com/report/SaraTrawnik/avatarme)

Library for github style Identicon generation.

## Installation

Get
```
go get -u github.com/SaraTrawnik/avatarme
```
Test
```
go test -cover github.com/SaraTrawnik/avatarme
```

## Usage

```go
package main

import (
  "github.com/SaraTrawnik/avatarme"
  "fmt"
)

func main () {
  ident := avatarme.New([]byte("test value"), "filename")
  fmt.Println(ident.Base64())
  ident.Draw() // saves resulting png image named "filename" to current directory
}
```
