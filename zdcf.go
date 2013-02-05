// Copyright 2013 Joshua Tacoma. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package provides methods for consuming ZeroMQ Device Configuration Files
// (ZDCF, http://rfc.zeromq.org/spec:17) for ØMQ (ZeroMQ, ZMQ,
// http://www.zeromq.org/) applications that use gozmq
// (http://godoc.org/github.com/alecthomas/gozmq).
//
package gozdcf

import (
	"errors"
	"fmt"
	"sync"

	zmq "github.com/alecthomas/gozmq"
)

func ListenAndServe(appName string, sources ...interface{}) error {
	var (
		wg  sync.WaitGroup
		app *app
		err error
	)
	app, err = newApp(appName, sources...)
	if err != nil {
		return fmt.Errorf("while creating app: %s", err)
	}
	defer app.Close()
	var runners []func()
	app.ForDevices(func(ctx *DeviceContext) {
		var (
			dev func(*DeviceContext)
			ok  bool
		)
		if err == nil {
			if dev, ok = lookupDevice(ctx.Type()); ok {
				runners = append(runners, func() {
					dev(ctx)
					wg.Done()
				})
			} else {
				err = fmt.Errorf("unregistered device type: %s", ctx.Type())
			}
		}
	})
	if err != nil {
		return err
	}
	if len(runners) == 0 {
		return fmt.Errorf("no devices loaded.")
	}
	for _, run := range runners {
		wg.Add(1)
		go run()
	}
	wg.Wait()
	return nil
}

// An app is a ØMQ context with a collection of devices.
type app struct {
	context zmq.Context
	name    string
	devices map[string]*DeviceContext
}

// Create the named app based on the specified configuration.
func newApp(appName string, sources ...interface{}) (a *app, err error) {
	var (
		conf    *zdcf1
		appConf *app1
		ok      bool
	)
	for _, source := range sources {
		var next *zdcf1
		if _, ok := source.(string); ok {
			source = []byte(source.(string))
		}
		switch source.(type) {
		case []byte:
			next, err = unmarshalZdcf1(source.([]byte))
			if err != nil {
				conf0, err0 := unmarshalZdcf0(source.([]byte))
				if err0 != nil {
					return nil, err0
				}
				next = conf0.zdcf1(appName)
				err = nil
			}
		case *zdcf1:
			next = source.(*zdcf1)
		}
		if next == nil {
			return nil, errors.New("unsupported configuration source.")
		}
		if conf == nil {
			conf = next
		} else {
			conf.update(next)
		}
	}
	if context, err := zmq.NewContext(); err != nil {
		return nil, err
	} else {
		a = &app{
			context: context,
			name:    appName,
			devices: map[string]*DeviceContext{},
		}
	}
	appConf, ok = conf.Apps[appName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no such app: %s", appName))
	}
	// TODO: context options (gozmq has no API for this yet)
	for devName, devConf := range appConf.Devices {
		devContext := &DeviceContext{
			app:     a,
			name:    devName,
			sockets: map[string]*socketContext{},
			typ:     devConf.Type,
		}
		for sockName, sockConf := range devConf.Sockets {
			sockContext := newSocketContext(devContext, sockName)
			switch sockConf.Type {
			case "PAIR":
				sockContext.Type = zmq.PAIR
			case "PUB":
				sockContext.Type = zmq.PUB
			case "SUB":
				sockContext.Type = zmq.SUB
			case "REQ":
				sockContext.Type = zmq.REQ
			case "REP":
				sockContext.Type = zmq.REP
			case "DEALER":
				sockContext.Type = zmq.DEALER
			case "ROUTER":
				sockContext.Type = zmq.ROUTER
			case "PULL":
				sockContext.Type = zmq.PULL
			case "PUSH":
				sockContext.Type = zmq.PUSH
			case "XPUB":
				sockContext.Type = zmq.XPUB
			case "XSUB":
				sockContext.Type = zmq.XSUB
			case "XREQ":
				sockContext.Type = zmq.XREQ
			case "XREP":
				sockContext.Type = zmq.XREP
			case "UPSTREAM":
				sockContext.Type = zmq.UPSTREAM
			case "DOWNSTREAM":
				sockContext.Type = zmq.DOWNSTREAM
			}
			// TODO: socket options
			sockContext.Bind = sockConf.Bind       // TODO: copy
			sockContext.Connect = sockConf.Connect // TODO: copy
			devContext.sockets[sockName] = sockContext
		}
		a.devices[devName] = devContext
	}
	return a, nil
}

// ForDevices calls the given function on each device.
func (a *app) ForDevices(do func(*DeviceContext)) {
	for _, devContext := range a.devices {
		do(devContext)
	}
}

// Close the app, including its ØMQ context.
//
// Note that this is constrained by ØMQ's rules for the destruction of its
// contexts, especially that a call to this method will block until all its
// devices' sockets have been closed.
func (a *app) Close() {
	if a != nil && a.context != nil {
		a.context.Close()
	}
}

// A DeviceContext is intended to be all that a ØMQ device needs to do its job.
type DeviceContext struct {
	app     *app
	name    string
	typ     string
	sockets map[string]*socketContext
}

// Type is the name of the device type intended to be instantiated.
//
// This is a string that should be translated to a func (or switch'd to a code
// block) that knows how to create that type of device.
func (d *DeviceContext) Type() string { return d.typ }

// Open creates and binds/connects the named socket.
func (d *DeviceContext) Open(name string) (sock zmq.Socket, err error) {
	var sockContext *socketContext
	var ok bool
	if sockContext, ok = d.sockets[name]; !ok {
		return nil, errors.New("no such socket.")
	}
	return sockContext.Open()
}

// MustOpen creates and binds/connects the named socket or else panics.
func (d *DeviceContext) MustOpen(name string) zmq.Socket {
	sock, err := d.Open(name)
	if err != nil {
		panic(err.Error())
	}
	return sock
}

// A socketContext represents all the information needed to create a socket.
//
// All properties that directly affect the construction, binding, and connecting
// of ØMQ sockets can be set here.  However, a socketContext must be associated
// with a DeviceContext in order to do its job i.e. to create and open a socket.
type socketContext struct {
	device        *DeviceContext
	name          string
	Type          zmq.SocketType
	IntOptions    map[zmq.IntSocketOption]int
	Int64Options  map[zmq.Int64SocketOption]int64
	UInt64Options map[zmq.UInt64SocketOption]uint64
	StringOptions map[zmq.StringSocketOption]string
	Bind          []string
	Connect       []string
}

func newSocketContext(device *DeviceContext, name string) *socketContext {
	if device == nil {
		panic("nil device")
	}
	return &socketContext{
		device:        device,
		name:          name,
		IntOptions:    map[zmq.IntSocketOption]int{},
		Int64Options:  map[zmq.Int64SocketOption]int64{},
		UInt64Options: map[zmq.UInt64SocketOption]uint64{},
		StringOptions: map[zmq.StringSocketOption]string{},
	}
}

// Name returns the name of the socket.
func (s *socketContext) Name() string { return s.name }

// Open a ØMQ socket.
//
// The socket will be affected by all options provided through the socketContext,
// including being bound and/or connected to some addresses: ready to go!
func (s *socketContext) Open() (sock zmq.Socket, err error) {
	var (
		DeviceContext *DeviceContext
		app           *app
	)
	if DeviceContext = s.device; DeviceContext == nil {
		return nil, errors.New("no device context.")
	}
	if app = DeviceContext.app; app == nil {
		return nil, errors.New("device context has no app.")
	}
	if sock, err = app.context.NewSocket(s.Type); err != nil {
		return nil, errors.New(fmt.Sprintf("could not create socket: %s", err.Error()))
	}
	defer func(s zmq.Socket) {
		if err != nil {
			s.Close()
		}
	}(sock)
	for opt, val := range s.IntOptions {
		if err = sock.SetSockOptInt(opt, val); err != nil {
			return nil, errors.New(fmt.Sprintf("could not set option %d = %v : %s",
				opt, val, err.Error()))
		}
	}
	for opt, val := range s.Int64Options {
		if err = sock.SetSockOptInt64(opt, val); err != nil {
			return nil, errors.New(fmt.Sprintf("could not set option %d = %v : %s",
				opt, val, err.Error()))
		}
	}
	for opt, val := range s.UInt64Options {
		if err = sock.SetSockOptUInt64(opt, val); err != nil {
			return nil, errors.New(fmt.Sprintf("could not set option %d = %v : %s",
				opt, val, err.Error()))
		}
	}
	for opt, val := range s.StringOptions {
		if err = sock.SetSockOptString(opt, val); err != nil {
			return nil, errors.New(fmt.Sprintf("could not set option %d = %v : %s",
				opt, val, err.Error()))
		}
	}
	for _, addr := range s.Bind {
		if err = sock.Bind(addr); err != nil {
			return nil, errors.New(fmt.Sprintf("could not bind to %v %s",
				addr, err.Error()))
		}
	}
	for _, addr := range s.Connect {
		if err = sock.Connect(addr); err != nil {
			return nil, errors.New(fmt.Sprintf("could not connect to %v %s",
				addr, err.Error()))
		}
	}
	return
}
