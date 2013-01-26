package zdcf

import (
	"testing"
)

func TestZdcf1(t *testing.T) {
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
								"connect": ["tcp://eth0:5556"]
							}
						}
					}
				}
			}
		}
	}`)
	conf, err := UnmarshalZdcf1(raw)
	if err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}
	if conf == nil {
		t.Fatalf("unmarshal returned two nils.")
	}
	if conf.Version != 1.0001 {
		t.Errorf("version = %f", conf.Version)
	}
	listener, ok := conf.Apps["listener"]
	if !ok {
		t.Errorf("apps does not contain %v", "listener")
	}
	if listener.Context.IoThreads != 1 {
		t.Fatalf("listener.context.iothreads = %v", listener.Context.IoThreads)
	}
	if listener.Context.Verbose != true {
		t.Fatalf("listener.context.verbose = %v", listener.Context.Verbose)
	}
	main, ok := listener.Devices["main"]
	if !ok {
		t.Fatalf("listerner.devices does not contain %v", "main")
	}
	if main.Type != "zmq_queue" {
		t.Fatalf("main.type = %v", main.Type)
	}
	frontend, ok := main.Sockets["frontend"]
	if !ok {
		t.Fatalf("main.sockets does not contain %v", "frontend")
	}
	if frontend.Type != "SUB" {
		t.Fatalf("frontend.type = %v", frontend.Type)
	}
	if frontend.Bind[0] != "tcp://eth0:5555" {
		t.Fatalf("frontend.bind = %v", frontend.Bind)
	}
	backend, ok := main.Sockets["backend"]
	if !ok {
		t.Fatalf("main.sockets does not contain %v", "backend")
	}
	if backend.Connect[0] != "tcp://eth0:5556" {
		t.Fatalf("backend.bind = %v", backend.Bind)
	}
}
