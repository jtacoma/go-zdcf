package main

import (
	"fmt"

	zdcf "github.com/jtacoma/gozdcf"
)

var defaults = &zdcf.Zdcf1{
	Version: 1.0,
	Apps: map[string]*zdcf.App1{
		"handlers": &zdcf.App1{
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

func EchoService(dev *zdcf.DeviceInfo) {
	if front, err := dev.OpenSocket("frontend"); err != nil {
		panic(err)
	} else {
		defer front.Close()
		for {
			msg, err := front.Recv(0)
			if err != nil {
				panic(err.Error())
			}
			front.Send(msg, 0)
		}
	}
}

func EchoClientOnce(dev *zdcf.DeviceInfo) {
	if back, err := dev.OpenSocket("backend"); err != nil {
		panic(err)
	} else {
		defer back.Close()
		back.Send([]byte("Echo!"), 0)
		msg, err := back.Recv(0)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(string(msg))
		done <- 1
	}
}

var done = make(chan int)

var handlers = map[string]func(*zdcf.DeviceInfo){
	"echo_service": EchoService,
	"echo_once":    EchoClientOnce,
}

func main() {
	listener, err := zdcf.NewApp("handlers", defaults)
	if err != nil {
		panic(err)
	}
	listener.ForDevices(func(dev *zdcf.DeviceInfo) {
		if run, ok := handlers[dev.Type()]; ok {
			go run(dev)
		}
	})
	<-done
}
