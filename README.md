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
