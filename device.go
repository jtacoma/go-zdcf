// Copyright 2013 Joshua Tacoma. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zdcf

import (
	"fmt"
	"regexp"

	zmq "github.com/alecthomas/gozmq"
)

func DeviceFunc(deviceTypePattern string, device func(*DeviceContext)) error {
	if pattern, err := regexp.Compile(deviceTypePattern); err != nil {
		return err
	} else {
		registry = append(registry, registration{pattern, device})
	}
	return nil
}

func builtinDevice(dev *DeviceContext) {
	var (
		typ         zmq.DeviceType
		back, front zmq.Socket
		err         error
	)
	switch dev.Type() {
	case "zmq_forwarder":
		typ = zmq.FORWARDER
	case "zmq_streamer":
		typ = zmq.STREAMER
	case "zmq_queue":
		typ = zmq.QUEUE
	default:
		panic(fmt.Sprintf("device has unknown type: %s.", dev.Type()))
	}
	back = dev.MustOpen("backend")
	front = dev.MustOpen("frontend")
	err = zmq.Device(typ, front, back)
	//err = zmq.Proxy( front, back, capture)
	if err != nil {
	}
}

type registration struct {
	pattern *regexp.Regexp
	device  func(*DeviceContext)
}

var registry = []registration{
	{regexp.MustCompile(`zmq_[a-z0-9_]*`), builtinDevice},
}

func lookupDevice(typeName string) (func(*DeviceContext), bool) {
	for i := len(registry) - 1; i >= 0; i-- {
		if registry[i].pattern.MatchString(typeName) {
			return registry[i].device, true
		}
	}
	return nil, false
}
