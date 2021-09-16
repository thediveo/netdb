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

// Protocol describes a network communications protocol by its native name and
// official protocol number as appearing within IP headers, with optional alias
// names.
//
// According to
// http://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml the
// Assigned Internet Protocol Numbers are 8 bit (unsigned) numbers.
//
// On purpose, we don't stick with the stuttering POSIX C library names, but
// instead aim for more Go-like type names. After all, Go isn't similar to C,
// except for using letters, signs, and braces.
type Protocol struct {
	Name    string   // Official protocol name.
	Number  uint8    // Protocol number.
	Aliases []string // List of aliases.
}

// ProtocolIndex indexes the known network communication protocols by either
// name (native as well as aliases) and by number.
type ProtocolIndex struct {
	Names   map[string]*Protocol // Index by protocol name, including aliases.
	Numbers map[uint8]*Protocol  // Index by protocol number.
}

// NewProtocolIndex returns a ProtocolsIndex object initialized with the
// specified protocols.
func NewProtocolIndex(protos []Protocol) ProtocolIndex {
	i := ProtocolIndex{
		Names:   map[string]*Protocol{},
		Numbers: map[uint8]*Protocol{},
	}
	i.Merge(protos)
	return i
}

// LoadProtocols returns a ProtocolIndex object initialized from the definitions
// in the named file.
func LoadProtocols(name string) (ProtocolIndex, error) {
	f, err := os.Open(name)
	if err != nil {
		return NewProtocolIndex(nil), err
	}
	defer f.Close()
	protos, err := ParseProtocols(f)
	if err != nil {
		return NewProtocolIndex(nil), err
	}
	return NewProtocolIndex(protos), nil
}

// Merge a list of Protocol descriptions into the current Protocols index,
// potentially overriding existing entries in the index in case of duplicates.
func (i *ProtocolIndex) Merge(protos []Protocol) {
	for idx, proto := range protos {
		// index by name, including aliases
		i.Names[proto.Name] = &protos[idx] // NEVER (re)use &proto! *facepalm*
		for _, alias := range proto.Aliases {
			i.Names[alias] = &protos[idx]
		}
		// index by protocol number
		i.Numbers[proto.Number] = &protos[idx]
	}
}

// MergeIndex merges another ProtocolIndex into the current index, potentially
// overriding existing entries in case of duplicates.
func (i *ProtocolIndex) MergeIndex(pi ProtocolIndex) {
	for name, proto := range pi.Names {
		i.Names[name] = proto
	}
	for number, proto := range pi.Numbers {
		i.Numbers[number] = proto
	}
}

// ParseProtocols parses Internet protocol definitions for the TCP/IP subsystem
// from the given Reader and returns them as a list of Protcol(s).
func ParseProtocols(r io.Reader) ([]Protocol, error) {
	protos := []Protocol{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(strings.SplitN(line, "#", 2)[0]) // There's always an element [0]
		if len(fields) < 2 {
			continue // skip empty lines and also silently ignore malformed lines.
		}

		proto, err := strconv.ParseUint(fields[1], 10, 8)
		if err != nil {
			return nil, err
		}

		protos = append(protos, Protocol{
			Name:    fields[0],
			Number:  uint8(proto), // note that we already checked in ParseUint(..., 8)
			Aliases: fields[2:],
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return protos, nil
}

// ProtocolByName returns the Protocol details for the specified (alias) name,
// or nil if not defined.
func ProtocolByName(name string) *Protocol {
	if Protocols.Numbers == nil {
		Protocols = NewProtocolIndex(BuiltinProtocols)
	}
	return Protocols.Names[name]
}

// ProtocolByNumber returns the Protocol details for the specified protocol
// number, or nil if not defined.
func ProtocolByNumber(number uint8) *Protocol {
	if Protocols.Numbers == nil {
		Protocols = NewProtocolIndex(BuiltinProtocols)
	}
	return Protocols.Numbers[number]
}

// Protocols is the index of protocol names and numbers. If left to the zero
// value then it will be automatically initialized with the builtin definitions
// upon first use of ProtocolByName or ProtocolByNumber.
var Protocols ProtocolIndex
