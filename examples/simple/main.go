// Copyright 2013 Joshua Tacoma. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	zdcf "github.com/jtacoma/gozdcf"
)

var defaults = `{
	"version": 1.0,
	"apps": {
		"listener": {
			"context": {
				"iothreads": 1,
				"verbose":   true
			},
			"devices": {
				"main": {
					"type": "zmq_queue",
					"sockets": {
						"frontend": {
							"type": "SUB",
							"options": {
								"hwm":  1000,
								"swap": 25000000
							}
						},
						"backend": {
							"type": "PUB"
						}
					}
				}
			}
		}
	}
}`

var custom = `{
	"version": 1.0001,
	"apps": {
		"listener": {
			"devices": {
				"main": {
					"sockets": {
						"frontend": {
							"option": {
								"subscribe": ["1234 ", "1235 "]
							},
							"bind": ["tcp://127.0.0.1:5555"]
						},
						"backend": {
							"connect": ["tcp://127.0.0.1:5556"]
						}
					}
				}
			}
		}
	}
}`

func main() {
	if listener, err := zdcf.NewApp("listener", defaults, custom); err != nil {
		panic(err)
	} else if main, ok := listener.Device("main"); !ok {
		panic("failed to load device.")
	} else if back, err := main.OpenSocket("backend"); err != nil {
		panic(err)
	} else if front, err := main.OpenSocket("frontend"); err != nil {
		panic(err)
	} else {
		fmt.Println("all ready!", front, back)
	}
}
