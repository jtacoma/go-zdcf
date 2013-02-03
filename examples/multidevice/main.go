// Copyright 2013 Joshua Tacoma. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	zmq "github.com/alecthomas/gozmq"
	zdcf "github.com/jtacoma/gozdcf"
)

var defaults = `{
	"version": 1.0,
	"apps": {
		"multidevice": {
			"devices": {
				"server": {
					"type": "echo_service",
					"sockets": {
						"frontend": {
							"type": "REP",
							"bind": ["tcp://127.0.0.1:5555"]
						}
					}
				},
				"client": {
					"type": "echo_once",
					"sockets": {
						"backend": {
							"type":    "REQ",
							"connect": ["tcp://127.0.0.1:5555"]
						}
					}
				}
			}
		}
	}
}`

func EchoService(front zmq.Socket) {
	defer front.Close()
	for {
		msg, err := front.Recv(0)
		if err != nil {
			panic(err.Error())
		}
		front.Send(msg, 0)
	}
}

func EchoClientOnce(back zmq.Socket) {
	defer back.Close()
	back.Send([]byte("Echo!"), 0)
	msg, err := back.Recv(0)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(msg))
}

func main() {
	listener, err := zdcf.NewApp("multidevice", defaults)
	if err != nil {
		panic(err)
	}
	done := make(chan int)
	listener.ForDevices(func(dev *zdcf.DeviceContext) {
		var (
			err         error
			front, back zmq.Socket
		)
		switch dev.Type() {
		case "echo_service":
			if front, err = dev.OpenSocket("frontend"); err != nil {
				panic(err)
			}
			go EchoService(front)
		case "echo_once":
			go func() {
				if back, err = dev.OpenSocket("backend"); err != nil {
					panic(err)
				}
				EchoClientOnce(back)
				done <- 1
			}()
		}
	})
	<-done
}
