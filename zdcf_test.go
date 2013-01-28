package gozdcf

import (
	"testing"
)

func TestZdcf(t *testing.T) {
	raw := []byte(`{
		"version": 1.0001,
		"apps": {
			"listener": {
				"context": {
					"iothreads": 1,
					"verbose": true
				},
				"devices": {
					"main": {
						"type": "zmq_queue",
						"sockets": {
							"frontend": {
								"type": "SUB",
								"option": {
									"hwm": 1000,
									"swap": 25000000
								},
								"bind": ["tcp://eth0:5555"]
							},
							"backend": {
								"bind": ["tcp://eth0:5556"]
							}
						}
					}
				}
			}
		}
	}`)
	conf, err := unmarshalZdcf1(raw)
	if err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}
	app, err := NewApp("listener", conf)
	if err != nil {
		t.Fatalf("failed to create app: %s", err)
	}
	defer app.Close()
	main, ok := app.Device("main")
	if !ok {
		t.Fatalf("failed to create device: main")
	}
	frontend, err := main.OpenSocket("frontend")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer frontend.Close()
	backendContext, ok := main.Socket("backend")
	if !ok {
		t.Fatalf("failed to find socket context: backend")
	}
	backend, err := backendContext.Open()
	if err != nil {
		t.Fatalf("failed to open socket: backend")
	}
	backend.Close()
}

func ExampleNewApp() {
	// This is a very simplified example that just shows the gist and does
	// not check for errors.
	defaults := `{
		"version": 1.0001,
		"apps": {
			"myapp": {
				"devices": {
					"echo": {
						"sockets": {
							"frontend": {
								"type": "REP",
								"bind": ["tcp://eth0:5555"]
							}
						}
					}
				}
			}
		}
	}`
	app, _ := NewApp("myapp", defaults)
	defer app.Close()
	echo, _ := app.Device("echo")
	front, _ := echo.OpenSocket("frontend")
	defer front.Close()
	for {
		msg, _ := front.Recv(0)
		front.Send(msg, 0)
	}
}
