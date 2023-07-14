// Copyright © 2019, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package status

import (
	"net"
	"os"
	"runtime"
	"strings"

	pscpu "github.com/shirou/gopsutil/cpu"
	psmem "github.com/shirou/gopsutil/mem"
	psnet "github.com/shirou/gopsutil/net"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/utils"
)

var logger = utils.MustGetLogger("server", "status.info")

// Iface represents an eth adapter
type Iface struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	VIP  string `json:"vip"`
	Mask string `json:"mask"`
}

// Info struct represents host information
type Info struct {
	Hostname        string                 `json:"hostname"`
	OperatingSystem string                 `json:"operating_system"`
	Architecture    string                 `json:"architecture"`
	Context         map[string]interface{} `json:"context"`
	Ifaces          []Iface                `json:"interfaces"`
	PS              PS                     `json:"ps"`
}

func localAddresses() ([]Iface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		logger.V(1).Info("local addresses: %+v", err)
		return nil, err
	}
	interfaces := make([]Iface, 0)
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			logger.V(1).Info("local addresses: %+v", err)
			continue
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPAddr:
				if v.String() != "" && !v.IP.IsLoopback() {
					iface := Iface{
						Name: i.Name,
						IP:   v.String(),
						Mask: v.IP.DefaultMask().String(),
					}
					interfaces = append(interfaces, iface)
				}
			case *net.IPNet:
				if v.String() != "" && !v.IP.IsLoopback() {
					iface := Iface{
						Name: i.Name,
						IP:   v.String(),
						VIP:  v.IP.String(),
						Mask: v.Mask.String(),
					}
					interfaces = append(interfaces, iface)
				}
			}
		}
	}
	return interfaces, nil
}

// getEnvsWith finds all ENV vars that contain a word
// TODO: use regex?
func getEnvsWith(match string) map[string]string {
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if strings.Contains(pair[0], match) {
			if len(pair[1]) > 0 {
				envs[pair[0]] = strings.TrimSpace(pair[1])
			}
		}
	}
	return envs
}

// GetContext grabs a list of env vars from the env and adds to a receipt
// cnab custom field must be a map[string]interface{} so we return a map
// for convenience
func localContext() map[string]interface{} {
	// A list of environment variables to remove that could be considered sensitive (like passwords, tokens, and api keys)
	filteredEnvVars := []string{"TOKEN", "APIKEY", "PASSWORD", "PASSWD"} // add new sensitive tags here

	matches := []string{"POD", "NODE", "SERVER"}
	envmap := make(map[string]interface{})
	for _, match := range matches {
		envs := getEnvsWith(match)
		if len(envs) > 0 {
			for k, v := range envs {
				// If the key contains ANY of the filtered words, we skip it
				filtered := false
				for _, env := range filteredEnvVars {
					if strings.Contains(strings.ToUpper(k), env) {
						filtered = true
						break
					}
				}
				if filtered {
					envmap[strings.ToLower(k)] = "**********"
				} else {
					envmap[strings.ToLower(k)] = v
				}
			}
		}
	}
	return envmap
}

// PS struct for ps like info
type PS struct {
	CPU interface{} `json:"cpu"`
	Mem interface{} `json:"memory"`
	Net interface{} `json:"network"`
}

func ps() (PS, error) {
	psout := PS{}
	cpu, err := pscpu.Info()
	if err != nil {
		return psout, err
	}
	psout.CPU = cpu
	network, err := psnet.Interfaces()
	if err != nil {
		return psout, err
	}
	psout.Net = network
	memory, err := psmem.VirtualMemory()
	if err != nil {
		return psout, err
	}
	psout.Mem = memory
	return psout, nil
}

// GetInfo returns system and host information
func GetInfo() *Info {
	info := &Info{}
	info.OperatingSystem = runtime.GOOS
	info.Architecture = runtime.GOARCH
	info.Context = localContext()
	hostname, err := os.Hostname()
	if err != nil {
		logger.V(1).Info("unable to determine hostname : %+v", err)
		return info
	}
	info.Hostname = hostname
	ifaces, err := localAddresses()
	if err != nil {
		logger.V(1).Info("unable to determine local interfaces : %+v", err)
		return info
	}
	info.Ifaces = ifaces
	psout, err := ps()
	if err != nil {
		logger.V(1).Info("unable to execute ps : %+v", err)
		return info
	}
	info.PS = psout
	return info
}
