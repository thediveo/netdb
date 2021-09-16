// Copyright 2021 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package netdb

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

// Service describes a network service by its official service name, port number
// and network protocol, with optional servicse alias names.
//
// On purpose, we don't stick with the stuttering POSIX C library names, but
// instead aim for more Go-like type names. After all, Go isn't similar to C,
// except for using letters, signs, and braces.
type Service struct {
	Name         string    // Official service name.
	Port         int       // Transport port number.
	ProtocolName string    // Name of protocol to use.
	Protocol     *Protocol // Protocol details, if known.
	Aliases      []string  // List of service name aliases.
}

// ServiceIndex indexes the known network services by either (alias) name or by
// transport port number.
type ServiceIndex struct {
	Names map[ServiceProtocol]*Service // Index by service name and protocol name.
	Ports map[ServicePort]*Service     // Index by port number.
}

// ServiceProtocol represents a Service index key.
type ServiceProtocol struct {
	Name     string // Service name.
	Protocol string // Protocol name; might be zero.
}

// ServicePort represents a Service index key.
type ServicePort struct {
	Port     int    // Transport port number.
	Protocol string // Protocol name; might be zero.
}

// NewServiceIndex returns a Services index object initialized with the
// specified services.
func NewServiceIndex(services []Service) ServiceIndex {
	i := ServiceIndex{
		Names: map[ServiceProtocol]*Service{},
		Ports: map[ServicePort]*Service{},
	}
	i.Merge(services)
	return i
}

// LoadServices returns a ServiceIndex object initialized from the
// definitions in the named file.
func LoadServices(name string, protos ProtocolIndex) (ServiceIndex, error) {
	f, err := os.Open(name)
	if err != nil {
		return NewServiceIndex(nil), err
	}
	defer f.Close()
	services, err := ParseServices(f, protos)
	if err != nil {
		return NewServiceIndex(nil), err
	}
	return NewServiceIndex(services), nil
}

// Merge a list of service descriptions into the current Services index,
// potentially overriding existing entries in the index in case of duplicates.
func (i *ServiceIndex) Merge(services []Service) {
	for idx, service := range services {
		// only register first transport-agnostic instance of a service.
		namekey := ServiceProtocol{Name: service.Name}
		if _, ok := i.Names[namekey]; !ok {
			i.Names[namekey] = &services[idx] // NEVER (re)use &service! *facepalm*
		}
		i.Names[ServiceProtocol{Name: service.Name, Protocol: service.ProtocolName}] = &services[idx]
		for _, alias := range service.Aliases {
			namekey := ServiceProtocol{Name: alias}
			if _, ok := i.Names[namekey]; !ok {
				i.Names[namekey] = &services[idx]
			}
			i.Names[ServiceProtocol{Name: alias, Protocol: service.ProtocolName}] = &services[idx]
		}
		// only register first transport-agnostic instance of a service.
		portkey := ServicePort{Port: service.Port}
		if _, ok := i.Ports[portkey]; !ok {
			i.Ports[portkey] = &services[idx]
		}
		i.Ports[ServicePort{Port: service.Port, Protocol: service.ProtocolName}] = &services[idx]
	}
}

// MergeIndex merges another ServiceIndex into the current index, potentially
// overriding existing entries in case of duplicates.
func (i *ServiceIndex) MergeIndex(si ServiceIndex) {
	for key, service := range si.Names {
		i.Names[key] = service
	}
	for key, service := range si.Ports {
		i.Ports[key] = service
	}
}

// ByName returns the named Service for the given protocol, or nil if not found.
// If the protocol is the zero value ("") then the "first" Service matching the
// name is returned, where "first" refers to the order in which the services
// were originally described in a list of services, such as /etc/services.
func (i *ServiceIndex) ByName(name string, protocol string) *Service {
	return i.Names[ServiceProtocol{Name: name, Protocol: protocol}]
}

// ByPort returns the service for the given port and protocol, or nil if not
// found. If the protocol is the zero value ("") then the "first" Service
// matching the name is returned, where "first" refers to the order in which the
// services were originally described in a list of services, such as
// /etc/services.
func (i *ServiceIndex) ByPort(port int, protocol string) *Service {
	return i.Ports[ServicePort{Port: port, Protocol: protocol}]
}

// ParseServices parses network service definitions from the given Reader and
// returns them as a list of Service(s).
func ParseServices(r io.Reader, p ProtocolIndex) ([]Service, error) {
	services := []Service{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(strings.SplitN(line, "#", 2)[0]) // There's always an element [0]
		if len(fields) < 2 {
			continue // skip empty lines and also silently ignore malformed lines.
		}

		portprotocol := strings.Split(fields[1], "/")
		if len(portprotocol) != 2 {
			continue
		}

		port, err := strconv.ParseUint(portprotocol[0], 10, 16)
		if err != nil {
			continue
		}

		proto, ok := p.Names[portprotocol[1]]
		if !ok {
			continue
		}

		services = append(services, Service{
			Name:         fields[0],
			Port:         int(port),
			ProtocolName: portprotocol[1],
			Protocol:     proto,
			Aliases:      fields[2:],
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

// ServiceByName returns the Service details for the specified (alias) name and
// (optional) protocol name, or nil if not defined.
func ServiceByName(name string, protocol string) *Service {
	if Services.Names == nil {
		Services = NewServiceIndex(BuiltinServices)
	}
	return Services.ByName(name, protocol)
}

// ServiceByPort returns the Service details for the specified port number and
// (optional) protocol name, or nil if not defined.
func ServiceByPort(port int, protocol string) *Service {
	if Services.Names == nil {
		Services = NewServiceIndex(BuiltinServices)
	}
	return Services.ByPort(port, protocol)
}

// Services is the index of service names and protocols. If left to the zero
// value then it will be automatically initialized with the builtin definitions
// upon first use of ServiceByName or ServiceByPort.
var Services ServiceIndex
