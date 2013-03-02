// Copyright 2013 Joshua Tacoma. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zdcf

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jtacoma/go-zpl"
)

type zdcf1 struct {
	Version float32          `version`
	Apps    map[string]*app1 `apps`
}

type app1 struct {
	Context *context1           `context`
	Devices map[string]*device1 `devices`
}

type context1 struct {
	IoThreads int  `iothreads`
	Verbose   bool `verbose`
}

type device1 struct {
	Type    string              `type`
	Sockets map[string]*socket1 `sockets`
}

type socket1 struct {
	Type    string    `type`
	Options *options1 `option`
	Bind    []string  `bind`
	Connect []string  `connect`
}

type options1 struct {
	Hwm         int      `hwm`
	Swap        int      `swap`
	Affinity    int      `affinity`
	Identity    string   `identity`
	Subscribe   []string `subscribe`
	Rate        int      `rate`
	RecoveryIvl int      `recovery_ivl`
	McastLoop   bool     `mcast_loop`
	SndBuf      int      `sndbuf`
	RcvBuf      int      `rcvbuf`
}

func unmarshalZdcf1(bytes []byte) (*zdcf1, error) {
	var conf zdcf1
	err_json := json.Unmarshal(bytes, &conf)
	if err_json != nil {
		err_zpl := zpl.Unmarshal(bytes, &conf)
		if err_zpl != nil {
			return nil, fmt.Errorf("failed to parse as JSON (%s) or as ZPL (%s).", err_json, err_zpl)
		}
	}
	if conf.Version < 1 || 2 <= conf.Version {
		return nil, errors.New(fmt.Sprintf(
			"unsupported ZDCF version: %f",
			conf.Version))
	}
	return &conf, nil
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
