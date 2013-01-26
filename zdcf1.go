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
			"invalid zdcf1 version: %f",
			zdcf1.Version))
	}
	return &zdcf1, err
}
