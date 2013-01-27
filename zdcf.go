// This package provides methods for consuming ZDCF in your ØMQ applications.
//
//     app := zdcf.NewApp("myapp", defaults, fileBytes)
//     defer app.Close()
//     echo := zdcf.Device("echo")
//     front := echo.OpenSocket("frontend")
//     defer front.Close()
//     for {
//         msg := front.Recv()
//         front.Send(msg)
//     }
// 
// A more complex example could enumerate an application's devices and use their
// type names to locate appropriate handlers.  See the examples directory for
// more.
//
// For more about ØMQ or ZDCF see http://rfc.zeromq.org/spec:17 or
// http://www.zeromq.org/, respectively.
package gozdcf

import (
	"errors"
	"fmt"

	zmq "github.com/alecthomas/gozmq"
)

// An App is a ZMQ context with a collection of devices.
type App struct {
	context zmq.Context
	name    string
	devices map[string]*DeviceContext
}

// Create the named App based on the specified configuration.
func NewApp(appName string, sources ...interface{}) (app *App, err error) {
	var (
		conf    *Zdcf1
		appConf *App1
		ok      bool
	)
	for _, source := range sources {
		var next *Zdcf1
		switch source.(type) {
		case string:
			next, err = UnmarshalZdcf1([]byte(source.(string)))
			if err != nil {
				return nil, err
			}
		case *Zdcf1:
			next = source.(*Zdcf1)
		}
		if next == nil {
			return nil, errors.New("unsupported configuration source.")
		}
		if conf == nil {
			conf = next
		} else {
			conf.Update(next)
		}
	}
	if context, err := zmq.NewContext(); err != nil {
		return nil, err
	} else {
		app = &App{
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
			app:     app,
			name:    devName,
			sockets: map[string]*SocketContext{},
			typ:     devConf.Type,
		}
		for sockName, sockConf := range devConf.Sockets {
			sockContext := NewSocketContext(devContext, sockName)
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
		app.devices[devName] = devContext
	}
	return app, nil
}

// Device returns the named device or else a second returned value of false.
func (a *App) Device(name string) (devContext *DeviceContext, ok bool) {
	devContext, ok = a.devices[name]
	return
}

func (a *App) ForDevices(do func(*DeviceContext)) {
	for _, devContext := range a.devices {
		do(devContext)
	}
}

// Close the App, including its ZMQ context.
func (a *App) Close() {
	if a != nil && a.context != nil {
		a.context.Close()
	}
}

type DeviceContext struct {
	app     *App
	name    string
	typ     string
	sockets map[string]*SocketContext
}

// Type is the name of the device type intended to be instantiated.
func (d *DeviceContext) Type() string { return d.typ }

// Device returns the named device or else a second returned value of false.
func (d *DeviceContext) Socket(name string) (sockContext *SocketContext, ok bool) {
	sockContext, ok = d.sockets[name]
	return
}

func (d *DeviceContext) OpenSocket(name string) (sock zmq.Socket, err error) {
	var sockContext *SocketContext
	var ok bool
	if sockContext, ok = d.sockets[name]; !ok {
		return nil, errors.New("no such socket.")
	}
	return sockContext.Open()
}

// A SocketContext represents all the information needed to create a socket.
type SocketContext struct {
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

func NewSocketContext(device *DeviceContext, name string) *SocketContext {
	if device == nil {
		panic("nil device")
	}
	return &SocketContext{
		device:        device,
		name:          name,
		IntOptions:    map[zmq.IntSocketOption]int{},
		Int64Options:  map[zmq.Int64SocketOption]int64{},
		UInt64Options: map[zmq.UInt64SocketOption]uint64{},
		StringOptions: map[zmq.StringSocketOption]string{},
	}
}

// Name returns the name of the socket.
func (s *SocketContext) Name() string { return s.name }

// Open a socket based on the socket context.
//
// The socket will be affected by all options provided through the SocketContext,
// including being bound and/or connected to some addresses.
func (s *SocketContext) Open() (sock zmq.Socket, err error) {
	var (
		deviceContext *DeviceContext
		app        *App
	)
	if deviceContext = s.device; deviceContext == nil {
		return nil, errors.New("no device context.")
	}
	if app = deviceContext.app; app == nil {
		return nil, errors.New("device context has no app.")
	}
	if sock, err = app.context.NewSocket(s.Type); err != nil {
		return nil, errors.New(fmt.Sprintf("could not create socket: %s", err.Error()))
	}
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
