gozdcf
======

This [Go](http://golang.org/) package provides methods for consuming [ZDCF](http://rfc.zeromq.org/spec:17) in your [ØMQ](http://www.zeromq.org/) applications.

    app := zdcf.NewApp("myapp", defaults, fileBytes)
    defer app.Close()
    echo := zdcf.Device("echo")
    front := echo.OpenSocket("frontend")
    defer front.Close()
    for {
        msg := front.Recv()
        front.Send(msg)
    }

A more complex example could enumerate an application's devices and use their type names to locate appropriate handlers.

Known Issues
------------

* [gozmq](https://github.com/alecthomas/gozmq), while excellent, doesn't yet support setting options on a context.
* when combining configuration sources, any socket options in later sources will completely replace options in previous sources.
* multi-valued settings (bind, connect, subscribe) in JSON will only accept arrays.
* configuration validation is mostly absent.

Since this is a young project, there are probably some pretty significant unknown issues.
