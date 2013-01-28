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

See also ØMQ (http://rfc.zeromq.org/spec:17), ZDCF (http://www.zeromq.org/),
and gozmq (http://godoc.org/github.com/alecthomas/gozmq).

See [godoc.org](http://godoc.org/github.com/jtacoma/gozdcf) for the familiar pretty docs.

## Known Issues

* gozmq, while excellent, doesn't yet support setting options on a context.
* when combining configuration sources, any socket options in later sources will completely replace options in previous sources.
* multi-valued settings (bind, connect, subscribe) in JSON will only accept arrays.
* configuration validation is mostly absent.

## License

Released under the MIT license, see LICENSE.txt.
