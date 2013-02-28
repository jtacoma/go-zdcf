// Copyright 2013 Joshua Tacoma. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zdcf

import (
	"fmt"
	"testing"
	"time"
)

func TestZdcf(t *testing.T) {
	conf := `{
		"version": 1.0001,
		"apps": {
			"listener": {
				"context": {
					"iothreads": 1,
					"verbose": true
				},
				"devices": {
					"main": {
						"type": "zmq_streamer",
						"sockets": {
							"frontend": {
								"type": "PULL",
								"option": {
									"hwm": 1000,
									"swap": 25000000
								},
								"bind": ["tcp://127.0.0.1:5555"]
							},
							"backend": {
								"type": "PUSH",
								"bind": ["tcp://127.0.0.1:5556"]
							}
						}
					},
					"sender": {
						"type": "test_send",
						"sockets": {
							"out": {
								"type": "PUSH",
								"connect": ["tcp://127.0.0.1:5555"]
							}
						}
					},
					"receiver": {
						"type": "test_recv",
						"sockets": {
							"in": {
								"type": "PULL",
								"connect": ["tcp://127.0.0.1:5556"]
							}
						}
					}
				}
			}
		}
	}`
	var (
		received_err     = make(chan error)
		received_message = make(chan string)
	)
	DeviceFunc("test_send", func(ctx *DeviceContext) {
		out := ctx.MustOpen("out")
		defer out.Close()
		out.Send([]byte("PASS"), 0)
	})
	DeviceFunc("test_recv", func(ctx *DeviceContext) {
		in := ctx.MustOpen("in")
		defer in.Close()
		msg, err := in.Recv(0)
		if err != nil {
			received_err <- err
		} else {
			received_message <- string(msg)
		}
	})
	go func() {
		err := ListenAndServe("listener", conf)
		if err != nil {
			t.Fatalf("failed to start: %s", err)
		}
	}()
	select {
	case err := <-received_err:
		t.Fatalf("received error: %s", err)
	case <-received_message:
		// :-)
	case <-time.After(1 * time.Second):
		t.Fatalf("timed out :-(")
	}
}

func Example() {
	defaults := `
version = 0.1
echo1
    type = echo_once
    cliff
        type = REP
        bind = tcp://127.0.0.1:5557

yodel42
    type = yodel_once
    vocals
        type = REQ
        connect = tcp://127.0.0.1:5557
`
	DeviceFunc("echo_once", func(ctx *DeviceContext) {
		cliff := ctx.MustOpen("cliff")
		defer cliff.Close()
		msg, err := cliff.Recv(0)
		if err != nil {
			panic(err)
		}
		cliff.Send(msg, 0)
	})
	DeviceFunc("yodel_once", func(ctx *DeviceContext) {
		vocals := ctx.MustOpen("vocals")
		defer vocals.Close()
		vocals.Send([]byte("YOOOODEL!"), 0)
		msg, _ := vocals.Recv(0)
		fmt.Println(string(msg))
	})
	err := ListenAndServe("myapp", defaults)
	if err != nil {
		fmt.Println("error:", err.Error())
	}
	// Output: YOOOODEL!
}
