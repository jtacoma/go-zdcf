package main

import (
	"log"

	zdcf "github.com/jtacoma/gozdcf"
)

var defaults = `{
	"version": 1.0,
	"apps": {
		"handlers": {
			"devices": {
				"server": {
					"type": "echo_service",
					"sockets": {
						"frontend": {
							"type": "REP",
							"bind": ["tcp://127.0.0.1:5555"]
						}
					}
				},
				"client": {
					"type": "echo_once",
					"sockets": {
						"backend": {
							"type":    "REQ",
							"connect": ["tcp://127.0.0.1:5555"]
						}
					}
				}
			}
		}
	}
}`

func EchoService(dev *zdcf.DeviceContext) {
	front, err := dev.OpenSocket("frontend")
	if err != nil {
		log.Println("error: echo service: no frontend socket:", err.Error())
		return
	}
	defer front.Close()
	for {
		msg, err := front.Recv(0)
		if err != nil {
			log.Println("error: echo service:", err.Error())
			return
		}
		front.Send(msg, 0)
	}
}

func EchoClientOnce(dev *zdcf.DeviceContext) {
	back, err := dev.OpenSocket("backend")
	if err != nil {
		log.Println("error: echo client: no backend socket:", err.Error())
		return
	}
	defer back.Close()
	back.Send([]byte("Echo!"), 0)
	msg, err := back.Recv(0)
	if err != nil {
		log.Println("error: echo client:", err.Error())
	}
	log.Println("debug: echo client: message received:", string(msg))
	done <- 1
}

var done = make(chan int)

var handlers = map[string]func(*zdcf.DeviceContext){
	"echo_service": EchoService,
	"echo_once":    EchoClientOnce,
}

func main() {
	app, err := zdcf.NewApp("handlers", defaults)
	if err != nil {
		log.Fatalln(err)
	}
	app.ForDevices(func(dev *zdcf.DeviceContext) {
		if run, ok := handlers[dev.Type()]; ok {
			go run(dev)
		} else {
			log.Println("error: unrecognized device type:", dev.Type())
		}
	})
	<-done
}
