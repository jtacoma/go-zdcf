// Copyright 2013 Joshua Tacoma. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zdcf

import (
	"encoding/json"
	"errors"
	"fmt"

	zpl "github.com/jtacoma/go-zpl"
)

type zdcf0 struct {
	Version float32             `version`
	Context *context1           `context`
	Devices map[string]*device0 `zpl:"*"`
}

type device0 struct {
	Type    string              `type`
	Sockets map[string]*socket1 `zpl:"*"`
}

func unmarshalZdcf0(bytes []byte) (*zdcf0, error) {
	var conf zdcf0
	err_json := json.Unmarshal(bytes, &conf)
	if err_json != nil {
		err_zpl := zpl.Unmarshal(bytes, &conf)
		if err_zpl != nil {
			return nil, fmt.Errorf("failed to parse as JSON (%s) or as ZPL (%s).", err_json, err_zpl)
		}
	}
	if conf.Version < 0 || 1 <= conf.Version {
		return nil, errors.New(fmt.Sprintf(
			"unsupported ZDCF version: %f",
			conf.Version))
	}
	return &conf, nil
}

func (z0 *zdcf0) zdcf1(appName string) *zdcf1 {
	devs := make(map[string]*device1)
	for name, d0 := range z0.Devices {
		devs[name] = &device1{
			Type:    d0.Type,
			Sockets: d0.Sockets,
		}
	}
	return &zdcf1{
		Version: 1.0,
		Apps: map[string]*app1{
			appName: &app1{
				Context: z0.Context,
				Devices: devs,
			},
		},
	}
}
