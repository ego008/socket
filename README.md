# socket

[![Build Status](https://travis-ci.org/mushroomsir/socket.svg?branch=master)](https://travis-ci.org/mushroomsir/socket)
[![Coverage Status](http://img.shields.io/coveralls/mushroomsir/socket.svg?style=flat-square)](https://coveralls.io/r/mushroomsir/socket)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/mushroomsir/socket/master/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/mushroomsir/socket)

## Installation and Usage
Install the package with:
```go
go get github.com/mushroomsir/socket
```
Import it with:
```go
import "github.com/mushroomsir/socket"
```
and use `socket` as the package name inside the code.
## Features
* Simple API: use it as an easy way to reuse socket(net.conn)
* Socket support Ping() func

## Example

```go 
   // create a new socket pool with an initial capacity of 5 and maximum
   // capacity of 30. The factory will create 5 initial connections and put it
   // into the pool.
   pool, err := socket.NewPool(1, 20, socket.NewFactory(address))
   // now you can get a connection from the pool, if there is no connection
   // available it will create a new one via the factory function.
   socket, err := pool.Get()
   // Ping to detect whether the socket is closed.
   b, err := socket.Ping()

   // close pool any time you want
   pool.Close()
```

