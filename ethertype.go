// Copyright 2024 Harald Albrecht.
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

// EtherType describes an Ethernet frame type/protocol used on an Ethernet
// network. Each EtherType entry contains a name and the two-octet (uint16)
// identitier. It also optionally contains name aliases as well as a comment.
//
// On purpose, we don't stick with the stuttering POSIX C library names, but
// instead aim for more Go-like type names. After all, Go isn't similar to C,
// except for using letters, signs, and braces.
type EtherType struct {
	Name    string   // Official EtherType name.
	Number  uint16   // EtherType number value.
	Aliases []string // List of aliases.
	Comment string   // Entry comment, if present.
}

// EtherTypeIndex index the known EtherTypes by either name (native as well as
// aliases) and by number.
type EtherTypeIndex struct {
	Names   map[string]*EtherType
	Numbers map[uint16]*EtherType
}

// NewEtherTypeIndex returns an EtherTypeIndex object initialized with the
// specified EtherTypes.
func NewEtherTypeIndex(ethertypes []EtherType) EtherTypeIndex {
	i := EtherTypeIndex{
		Names:   map[string]*EtherType{},
		Numbers: map[uint16]*EtherType{},
	}
	i.Merge(ethertypes)
	return i
}

// LoadEtherTypes returns an EtherTypeIndex object initialized from the
// defintions in the named file.
func LoadEtherTypes(name string) (EtherTypeIndex, error) {
	f, err := os.Open(name)
	if err != nil {
		return NewEtherTypeIndex(nil), err
	}
	defer f.Close()
	ethertypes, err := ParseEtherTypes(f)
	if err != nil {
		return NewEtherTypeIndex(nil), err
	}
	return NewEtherTypeIndex(ethertypes), nil
}

// Merge a list of EtherType descriptions into the current EtherTypes index,
// potentially overriding existing entries in the index in case of duplicates.
func (i *EtherTypeIndex) Merge(ethertypes []EtherType) {
	for idx, ethertype := range ethertypes {
		i.Names[ethertype.Name] = &ethertypes[idx]
		for _, alias := range ethertype.Aliases {
			i.Names[alias] = &ethertypes[idx]
		}
		i.Numbers[ethertype.Number] = &ethertypes[idx]
	}
}

// MergeIndex merges another EtherTypeIndex into the current index, potentially
// overriding existing enties in the case of duplicates.
func (i *EtherTypeIndex) MergeIndex(eti EtherTypeIndex) {
	for name, ethertype := range eti.Names {
		i.Names[name] = ethertype
	}
	for number, ethertype := range eti.Numbers {
		i.Numbers[number] = ethertype
	}
}

// ParseEtherTypes parses EtherType definitions from the given Reader and
// returns them as a list of EtherType objects.
func ParseEtherTypes(r io.Reader) ([]EtherType, error) {
	ethertypes := []EtherType{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			// Skip lines containing only comments
			continue
		}
		components := strings.SplitN(line, "#", 2)
		comment := ""
		if len(components) > 1 {
			comment = strings.TrimSpace(components[1])
		}
		fields := strings.Fields(components[0])
		if len(fields) < 2 {
			// Skip malformed entries
			continue
		}
		number, err := strconv.ParseUint(fields[1], 16, 16)
		if err != nil {
			return nil, err
		}
		ethertypes = append(ethertypes, EtherType{
			Name:    fields[0],
			Number:  uint16(number),
			Aliases: fields[2:],
			Comment: comment,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ethertypes, nil
}

// EtherTypeByName returns the EtherType details for the specified (native or
// aliased) name, or nil if not defined.
func EtherTypeByName(name string) *EtherType {
	if EtherTypes.Numbers == nil {
		EtherTypes = NewEtherTypeIndex(BuiltinEtherTypes)
	}
	return EtherTypes.Names[name]
}

// EtherTypeByNumber returns the EtherType details for the specified EtherType
// number, or nil if not defined.
func EtherTypeByNumber(number uint16) *EtherType {
	if EtherTypes.Numbers == nil {
		EtherTypes = NewEtherTypeIndex(BuiltinEtherTypes)
	}
	return EtherTypes.Numbers[number]
}

// EtherTypes is the index of EtherType names and numbers. If left to the zero
// value, then it will be automatically initialized with the builtin
// definitions upon first use of EtherTypeByName or EtherTypeByNumber.
var EtherTypes EtherTypeIndex
