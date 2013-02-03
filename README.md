# gozdcf

This package provides methods for consuming ZeroMQ Device Configuration Files
(ZDCF, http://rfc.zeromq.org/spec:17) for ØMQ (ZeroMQ, ZMQ,
http://www.zeromq.org/) applications that use gozmq
(http://godoc.org/github.com/alecthomas/gozmq).

## Documentation

See godoc or http://godoc.org/github.com/jtacoma/gozdcf

## Known Issues

* there is no support for setting options on a context.
* when combining configuration sources, any socket options in later sources will replace all options in previous sources.
* multi-valued settings (bind, connect, subscribe) in JSON will only accept arrays.
* configuration validation is mostly absent.

## License

Released under a BSD-style license that can be found in the LICENSE file.
