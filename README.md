# gozdcf

This package provides methods for consuming ZDCF in your ØMQ applications.

See also ØMQ (http://rfc.zeromq.org/spec:17), ZDCF (http://www.zeromq.org/),
and gozmq (http://godoc.org/github.com/alecthomas/gozmq).

See [godoc.org](http://godoc.org/github.com/jtacoma/gozdcf) for the familiar pretty docs.

## Example

This is a very simplified example that just shows the gist and does not check for errors.

    defaults := `{
        "version": 1.0001,
        "apps": {
            "myapp": {
                "devices": {
                    "echo": {
                        "sockets": {
                            "frontend": {
                                "type": "REP",
                                "bind": ["tcp://eth0:5555"]
                            }
                        }
                    }
                }
            }
        }
    }`
    app, _ := NewApp("myapp", defaults)
    defer app.Close()
    echo, _ := app.Device("echo")
    front, _ := echo.OpenSocket("frontend")
    defer front.Close()
    for {
        msg, _ := front.Recv(0)
        front.Send(msg, 0)
    }

## Known Issues

* gozmq, while excellent, doesn't yet support setting options on a context.
* when combining configuration sources, any socket options in later sources will completely replace options in previous sources.
* multi-valued settings (bind, connect, subscribe) in JSON will only accept arrays.
* configuration validation is mostly absent.

## License

Released under the MIT license, see LICENSE.txt.
