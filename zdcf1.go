package zdcf

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Zdcf1 struct {
	Version float32          `json:"version"`
	Apps    map[string]*App1 `json:"apps"`
}

type App1 struct {
	Context *Context1           `json:"context"`
	Devices map[string]*Device1 `json:"devices"`
}

type Context1 struct {
	IoThreads int  `json:"iothreads"`
	Verbose   bool `json:"verbose"`
}

type Device1 struct {
	Type    string              `json:"type"`
	Sockets map[string]*Socket1 `json:"sockets"`
}

type Socket1 struct {
	Type    string    `json:"type"`
	Options *Options1 `json:"option"`
	Bind    []string  `json:"bind"`
	Connect []string  `json:"bind"`
}

type Options1 struct {
	Hwm         int      `json:"hwm"`
	Swap        int      `json:"swap"`
	Affinity    int      `json:"affinity"`
	Identity    string   `json:"identity"`
	Subscribe   []string `json:"subscribe"`
	Rate        int      `json:"rate"`
	RecoveryIvl int      `json:"recovery_ivl"`
	McastLoop   bool     `json:"mcast_loop"`
	SndBuf      int      `json:"sndbuf"`
	RcvBuf      int      `json:"rcvbuf"`
}

func UnmarshalZdcf1(bytes []byte) (*Zdcf1, error) {
	var zdcf1 Zdcf1
	err := json.Unmarshal(bytes, &zdcf1)
	if err != nil {
		return nil, err
	}
	if zdcf1.Version < 1 || 2 <= zdcf1.Version {
		return nil, errors.New(fmt.Sprintf(
			"unsupported zdcf1 version: %f",
			zdcf1.Version))
	}
	return &zdcf1, err
}

func (conf *Zdcf1) Update(other *Zdcf1) error {
	if other.Version < 1 || 2 <= other.Version {
		return errors.New(fmt.Sprintf(
			"unsupported zdcf1 version: %f",
			other.Version))
	}
	for appName, appConf1 := range other.Apps {
		if appConf0, already := conf.Apps[appName]; !already {
			conf.Apps[appName] = appConf1
		} else {
			// TODO: context? gozmq provides no API for this.
			for devName, devConf1 := range appConf1.Devices {
				if devConf0, already := appConf0.Devices[devName]; !already {
					appConf0.Devices[devName] = devConf1
				} else {
					for sockName, sockConf1 := range devConf1.Sockets {
						if sockConf0, already := devConf0.Sockets[sockName]; !already {
							devConf0.Sockets[sockName] = sockConf1
						} else {
							if len(sockConf1.Type) > 0 {
								sockConf0.Type = sockConf1.Type
							}
							if sockConf1.Options != nil {
								// TODO: do a proper update here!
								sockConf0.Options = sockConf1.Options
							}
							if len(sockConf1.Bind) > 0 {
								sockConf0.Bind = sockConf1.Bind
							}
							if len(sockConf1.Connect) > 0 {
								sockConf0.Connect = sockConf1.Connect
							}
						}
					}
				}
			}
		}
	}
	return nil
}
