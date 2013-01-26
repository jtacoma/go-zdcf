package main

import (
	"fmt"

	zdcf "github.com/jtacoma/gozdcf"
)

var conf = &zdcf.Zdcf1{
	Version: 1.0,
	Apps: map[string]*zdcf.App1{
		"listener": &zdcf.App1{
			Context: &zdcf.Context1{
				IoThreads: 1,
				Verbose:   true,
			},
			Devices: map[string]*zdcf.Device1{
				"main": &zdcf.Device1{
					Type: "zmq_queue",
					Sockets: map[string]*zdcf.Socket1{
						"frontend": &zdcf.Socket1{
							Type: "SUB",
							Options: &zdcf.Options1{
								Hwm:  1000,
								Swap: 25000000,
								Subscribe: []string{
									"1234 ",
									"1235 ",
								},
							},
							Bind: []string{"tcp://eth0:5555"},
						},
						"backend": &zdcf.Socket1{
							Bind: []string{"tcp://eth0:5556"},
						},
					},
				},
			},
		},
	},
}

func main() {
	if listener, err := zdcf.NewApp("listener", conf); err != nil {
		panic(err)
	} else if main, ok := listener.Device("main"); !ok {
		panic("failed to load device.")
	} else if back, err := main.OpenSocket("backend"); err != nil {
		panic(err)
	} else if front, err := main.OpenSocket("frontend"); err != nil {
		panic(err)
	} else {
		fmt.Println("all ready!", front, back)
	}
}
