package gozdcf

import (
	"encoding/json"
	"errors"
	"fmt"
)

type zdcf1 struct {
	Version float32          `json:"version"`
	Apps    map[string]*app1 `json:"apps"`
}

type app1 struct {
	Context *context1           `json:"context"`
	Devices map[string]*device1 `json:"devices"`
}

type context1 struct {
	IoThreads int  `json:"iothreads"`
	Verbose   bool `json:"verbose"`
}

type device1 struct {
	Type    string              `json:"type"`
	Sockets map[string]*socket1 `json:"sockets"`
}

type socket1 struct {
	Type    string    `json:"type"`
	Options *options1 `json:"option"`
	Bind    []string  `json:"bind"`
	Connect []string  `json:"bind"`
}

type options1 struct {
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

func unmarshalZdcf1(bytes []byte) (*zdcf1, error) {
	var conf zdcf1
	err := json.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, err
	}
	if conf.Version < 1 || 2 <= conf.Version {
		return nil, errors.New(fmt.Sprintf(
			"unsupported ZDCF version: %f",
			conf.Version))
	}
	return &conf, err
}

func (conf *zdcf1) update(other *zdcf1) error {
	if other.Version < 1 || 2 <= other.Version {
		return errors.New(fmt.Sprintf(
			"unsupported ZDCF version: %f",
			other.Version))
	}
	for appName, appConf := range other.Apps {
		if appConf0, already := conf.Apps[appName]; !already {
			conf.Apps[appName] = appConf
		} else {
			// TODO: context? gozmq provides no API for this.
			for devName, devConf := range appConf.Devices {
				if devConf0, already := appConf0.Devices[devName]; !already {
					appConf0.Devices[devName] = devConf
				} else {
					for sockName, sockConf := range devConf.Sockets {
						if sockConf0, already := devConf0.Sockets[sockName]; !already {
							devConf0.Sockets[sockName] = sockConf
						} else {
							if len(sockConf.Type) > 0 {
								sockConf0.Type = sockConf.Type
							}
							if sockConf.Options != nil {
								// TODO: do a proper update here!
								sockConf0.Options = sockConf.Options
							}
							if len(sockConf.Bind) > 0 {
								sockConf0.Bind = sockConf.Bind
							}
							if len(sockConf.Connect) > 0 {
								sockConf0.Connect = sockConf.Connect
							}
						}
					}
				}
			}
		}
	}
	return nil
}
