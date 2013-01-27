# gozdcf

This package provides methods for consuming ZDCF in your ØMQ applications.

    app := zdcf.NewApp("myapp", defaults, fileBytes)
    defer app.Close()
    echo := zdcf.Device("echo")
    front := echo.OpenSocket("frontend")
    defer front.Close()
    for {
        msg := front.Recv()
        front.Send(msg)
    }

A more complex example could enumerate an application's devices and use their
type names to locate appropriate handlers.  See the examples directory for
more.

For more about ØMQ or ZDCF see http://rfc.zeromq.org/spec:17 or
http://www.zeromq.org/, respectively.

## Usage

#### type App

```go
type App struct {
}
```

An App is a ZMQ context with a collection of devices.

#### func  NewApp

```go
func NewApp(appName string, sources ...interface{}) (app *App, err error)
```
Create the named App based on the specified configuration.

#### func (*App) Close

```go
func (a *App) Close()
```
Close the App, including its ZMQ context.

#### func (*App) Device

```go
func (a *App) Device(name string) (devContext *DeviceContext, ok bool)
```
Device returns the named device or else a second returned value of false.

#### func (*App) ForDevices

```go
func (a *App) ForDevices(do func(*DeviceContext))
```

#### type DeviceContext

```go
type DeviceContext struct {
}
```


#### func (*DeviceContext) OpenSocket

```go
func (d *DeviceContext) OpenSocket(name string) (sock zmq.Socket, err error)
```

#### func (*DeviceContext) Socket

```go
func (d *DeviceContext) Socket(name string) (sockContext *SocketContext, ok bool)
```
Device returns the named device or else a second returned value of false.

#### func (*DeviceContext) Type

```go
func (d *DeviceContext) Type() string
```
Type is the name of the device type intended to be instantiated.

#### type SocketContext

```go
type SocketContext struct {
	Type          zmq.SocketType
	IntOptions    map[zmq.IntSocketOption]int
	Int64Options  map[zmq.Int64SocketOption]int64
	UInt64Options map[zmq.UInt64SocketOption]uint64
	StringOptions map[zmq.StringSocketOption]string
	Bind          []string
	Connect       []string
}
```

A SocketContext represents all the information needed to create a socket.

#### func  NewSocketContext

```go
func NewSocketContext(device *DeviceContext, name string) *SocketContext
```

#### func (*SocketContext) Name

```go
func (s *SocketContext) Name() string
```
Name returns the name of the socket.

#### func (*SocketContext) Open

```go
func (s *SocketContext) Open() (sock zmq.Socket, err error)
```
Open a socket based on the socket context.

The socket will be affected by all options provided through the SocketContext,
including being bound and/or connected to some addresses.

## Known Issues

* gozmq, while excellent, doesn't yet support setting options on a context.
* when combining configuration sources, any socket options in later sources will completely replace options in previous sources.
* multi-valued settings (bind, connect, subscribe) in JSON will only accept arrays.
* configuration validation is mostly absent.
