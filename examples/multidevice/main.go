package main

import (
	"fmt"

	zmq "github.com/alecthomas/gozmq"
	zdcf "github.com/jtacoma/gozdcf"
)

var defaults = &zdcf.Zdcf1{
	Version: 1.0,
	Apps: map[string]*zdcf.App1{
		"multidevice": &zdcf.App1{
			Devices: map[string]*zdcf.Device1{
				"server": &zdcf.Device1{
					Type: "echo_service",
					Sockets: map[string]*zdcf.Socket1{
						"frontend": &zdcf.Socket1{
							Type: "REP",
							Bind: []string{"tcp://127.0.0.1:5555"},
						},
					},
				},
				"client": &zdcf.Device1{
					Type: "echo_once",
					Sockets: map[string]*zdcf.Socket1{
						"backend": &zdcf.Socket1{
							Type:    "REQ",
							Connect: []string{"tcp://127.0.0.1:5555"},
						},
					},
				},
			},
		},
	},
}

func EchoService(front zmq.Socket) {
	defer front.Close()
	for {
		msg, err := front.Recv(0)
		if err != nil {
			panic(err.Error())
		}
		front.Send(msg, 0)
	}
}

func EchoClientOnce(back zmq.Socket) {
	defer back.Close()
	back.Send([]byte("Echo!"), 0)
	msg, err := back.Recv(0)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(msg))
}

func main() {
	listener, err := zdcf.NewApp("multidevice", defaults)
	if err != nil {
		panic(err)
	}
	done := make(chan int)
	listener.ForDevices(func(dev *zdcf.DeviceInfo) {
		var (
			err         error
			front, back zmq.Socket
		)
		switch dev.Type() {
		case "echo_service":
			if front, err = dev.OpenSocket("frontend"); err != nil {
				panic(err)
			}
			go EchoService(front)
		case "echo_once":
			go func() {
				if back, err = dev.OpenSocket("backend"); err != nil {
					panic(err)
				}
				EchoClientOnce(back)
				done <- 1
			}()
		}
	})
	<-done
}
