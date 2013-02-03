// Copyright 2013 Joshua Tacoma. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gozdcf

import (
	"testing"
)

func TestUnmarshalZdcf0_ZPL(t *testing.T) {
	raw := []byte(`
version = 0.1
context
    iothreads = 1
    verbose = true
main
    type = zmq_queue
    frontend
        type = SUB
        option
            hwm = 1000
            swap = 250000000
        bind = tcp://eth0:5555
    backend
        connect = tcp://eth0:5556`)
	conf, err := unmarshalZdcf0(raw)
	if err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}
	if conf == nil {
		t.Fatalf("unmarshal returned two nils.")
	}
	if conf.Version != 0.1 {
		t.Errorf("version = %f", conf.Version)
	}
	main, ok := conf.Devices["main"]
	if !ok {
		t.Fatalf("conf.devices does not contain %v", "main")
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
