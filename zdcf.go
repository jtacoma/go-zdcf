package zdcf

import (
	"errors"
	"fmt"

	zmq "github.com/alecthomas/gozmq"
)

// An App is a ZMQ context with a collection of devices.
type App struct {
	context zmq.Context
	name    string
	devices map[string]*DeviceInfo
}

// Create the named App based on the specified configuration.
func NewApp(appName string, sources ...interface{}) (app *App, err error) {
	var (
		conf    *Zdcf1
		appConf *App1
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
			devices: map[string]*DeviceInfo{},
		}
	}
	appConf = conf.Apps[appName]
	// TODO: context options (gozmq has no API for this yet)
	for devName, devConf := range appConf.Devices {
		devInfo := &DeviceInfo{
			app:     app,
			name:    devName,
			sockets: map[string]*SocketInfo{},
			typ:     devConf.Type,
		}
		for sockName, sockConf := range devConf.Sockets {
			sockInfo := NewSocketInfo(devInfo, sockName)
			switch sockConf.Type {
			case "PAIR":
				sockInfo.Type = zmq.PAIR
			case "PUB":
				sockInfo.Type = zmq.PUB
			case "SUB":
				sockInfo.Type = zmq.SUB
			case "REQ":
				sockInfo.Type = zmq.REQ
			case "REP":
				sockInfo.Type = zmq.REP
			case "DEALER":
				sockInfo.Type = zmq.DEALER
			case "ROUTER":
				sockInfo.Type = zmq.ROUTER
			case "PULL":
				sockInfo.Type = zmq.PULL
			case "PUSH":
				sockInfo.Type = zmq.PUSH
			case "XPUB":
				sockInfo.Type = zmq.XPUB
			case "XSUB":
				sockInfo.Type = zmq.XSUB
			case "XREQ":
				sockInfo.Type = zmq.XREQ
			case "XREP":
				sockInfo.Type = zmq.XREP
			case "UPSTREAM":
				sockInfo.Type = zmq.UPSTREAM
			case "DOWNSTREAM":
				sockInfo.Type = zmq.DOWNSTREAM
			}
			// TODO: socket options
			sockInfo.Bind = sockConf.Bind       // TODO: copy
			sockInfo.Connect = sockConf.Connect // TODO: copy
			devInfo.sockets[sockName] = sockInfo
		}
		app.devices[devName] = devInfo
	}
	return app, nil
}

// Device returns the named device or else a second returned value of false.
func (a *App) Device(name string) (devInfo *DeviceInfo, ok bool) {
	devInfo, ok = a.devices[name]
	return
}

// Close the App, including its ZMQ context.
func (a *App) Close() {
	if a != nil && a.context != nil {
		a.context.Close()
	}
}

type DeviceInfo struct {
	app     *App
	name    string
	typ     string
	sockets map[string]*SocketInfo
}

// Type is the name of the device type intended to be instantiated.
func (d *DeviceInfo) Type() string { return d.typ }

// Device returns the named device or else a second returned value of false.
func (d *DeviceInfo) Socket(name string) (sockInfo *SocketInfo, ok bool) {
	sockInfo, ok = d.sockets[name]
	return
}

func (d *DeviceInfo) OpenSocket(name string) (sock zmq.Socket, err error) {
	var sockInfo *SocketInfo
	var ok bool
	if sockInfo, ok = d.sockets[name]; !ok {
		return nil, errors.New("no such socket.")
	}
	return sockInfo.Open()
}

// A SocketInfo represents all the information needed to create a socket.
type SocketInfo struct {
	device        *DeviceInfo
	name          string
	Type          zmq.SocketType
	IntOptions    map[zmq.IntSocketOption]int
	Int64Options  map[zmq.Int64SocketOption]int64
	UInt64Options map[zmq.UInt64SocketOption]uint64
	StringOptions map[zmq.StringSocketOption]string
	Bind          []string
	Connect       []string
}

func NewSocketInfo(device *DeviceInfo, name string) *SocketInfo {
	if device == nil {
		panic("nil device")
	}
	return &SocketInfo{
		device:        device,
		name:          name,
		IntOptions:    map[zmq.IntSocketOption]int{},
		Int64Options:  map[zmq.Int64SocketOption]int64{},
		UInt64Options: map[zmq.UInt64SocketOption]uint64{},
		StringOptions: map[zmq.StringSocketOption]string{},
	}
}

func (s *SocketInfo) Name() string { return s.name }

// Open a socket based on the socket info.
//
// The socket will be affected by all options provided through the SocketInfo,
// including being bound and/or connected to some addresses.
func (s *SocketInfo) Open() (sock zmq.Socket, err error) {
	var (
		deviceInfo *DeviceInfo
		app        *App
	)
	if deviceInfo = s.device; deviceInfo == nil {
		return nil, errors.New("no device info.")
	}
	if app = deviceInfo.app; app == nil {
		return nil, errors.New("device info has no app.")
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
