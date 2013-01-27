package gozdcf

import (
	"testing"
)

func TestunmarshalZdcf1(t *testing.T) {
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
	conf, err := unmarshalZdcf1(raw)
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

func TestZdcf1_update(t *testing.T) {
	var conf = &zdcf1{
		Version: 1.0,
		Apps: map[string]*app1{
			"listener": &app1{
				Context: &context1{
					IoThreads: 1,
					Verbose:   true,
				},
				Devices: map[string]*device1{
					"main": &device1{
						Type: "zmq_queue",
						Sockets: map[string]*socket1{
							"frontend": &socket1{
								Type: "SUB",
								Options: &options1{
									Hwm:       1000,
									Swap:      25000000,
									Subscribe: []string{"4321 "},
								},
								Bind: []string{"tcp://eth0:1111"},
							},
						},
					},
				},
			},
		},
	}
	conf.update(&zdcf1{
		Version: 1.0,
		Apps: map[string]*app1{
			"listener": &app1{
				Context: &context1{
					IoThreads: 1,
					Verbose:   true,
				},
				Devices: map[string]*device1{
					"main": &device1{
						Type: "zmq_queue",
						Sockets: map[string]*socket1{
							"frontend": &socket1{
								Options: &options1{
									Subscribe: []string{
										"1234 ",
										"1235 ",
									},
								},
								Bind: []string{"tcp://eth0:5555"},
							},
							"backend": &socket1{
								Connect: []string{"tcp://eth0:5556"},
							},
						},
					},
				},
			},
			"speaker": &app1{},
		},
	})
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
	_, ok = conf.Apps["speaker"]
	if !ok {
		t.Errorf("apps does not contain %v", "speaker")
	}
}
