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
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("services", func() {

	var protos ProtocolIndex

	BeforeEach(func() {
		p, err := ParseProtocols(strings.NewReader(`
foobar	12
baz		234
`))
		Expect(err).NotTo(HaveOccurred())
		protos = NewProtocolIndex(p)
	})

	Context("parsing descriptions", func() {

		It("returns correct descriptions", func() {
			s, err := ParseServices(strings.NewReader(`
crash 666/foobar burn
crash 666/baz burn
`), protos)
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(ConsistOf(
				MatchFields(IgnoreExtras, Fields{
					"Name":         Equal("crash"),
					"Port":         Equal(666),
					"ProtocolName": Equal("foobar"),
					"Protocol": PointTo(MatchFields(IgnoreExtras, Fields{
						"Name":   Equal("foobar"),
						"Number": Equal(uint8(12)),
					})),
					"Aliases": ConsistOf("burn"),
				}),
				MatchFields(IgnoreExtras, Fields{
					"Name":         Equal("crash"),
					"Port":         Equal(666),
					"ProtocolName": Equal("baz"),
					"Protocol": PointTo(MatchFields(IgnoreExtras, Fields{
						"Name":   Equal("baz"),
						"Number": Equal(uint8(234)),
					})),
					"Aliases": ConsistOf("burn"),
				}),
			))
		})

		It("ignores comments and empty lines without errors", func() {
			s, err := ParseServices(strings.NewReader(`
# A comment
			# Another comment
crash 666/foobar burn # And this one.
`), protos)
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(HaveLen(1))
		})

		It("silently skips malformed definitions (including non-defined protocols)", func() {
			s, err := ParseServices(strings.NewReader(`
crash and burn
crash and/and burn
crash 666/and burn
`), protos)
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(HaveLen(0))
		})

		It("reports scanner errors", func() {
			f, err := os.Open("service_test.go")
			Expect(err).NotTo(HaveOccurred())
			f.Close() // sic! no defer!
			_, err = ParseServices(f, ProtocolIndex{})
			Expect(err).To(HaveOccurred())
		})

	})

	Context("loading", func() {

		It("loads service descriptions from file", func() {
			_, err := LoadServices("test/non-existing-services", ProtocolIndex{})
			Expect(err).To(HaveOccurred())

			idx, err := LoadServices("test/services", protos)
			Expect(err).NotTo(HaveOccurred())
			Expect(idx.ByPort(666, "foobar")).NotTo(BeNil())
		})

	})

	Context("indexing", func() {

		It("builds index", func() {
			s, err := ParseServices(strings.NewReader(`
crash 666/foobar burn
crash 666/baz burn
`), protos)
			Expect(err).NotTo(HaveOccurred())
			idx := NewServiceIndex(s)

			Expect(idx.Names).To(HaveLen(6))
			Expect(idx.Names).To(HaveKey(ServiceProtocol{Name: "crash", Protocol: ""}))
			Expect(idx.Names).To(HaveKey(ServiceProtocol{Name: "crash", Protocol: "foobar"}))
			Expect(idx.Names).To(HaveKey(ServiceProtocol{Name: "crash", Protocol: "baz"}))
			Expect(idx.Names).To(HaveKey(ServiceProtocol{Name: "burn", Protocol: ""}))
			Expect(idx.Names).To(HaveKey(ServiceProtocol{Name: "burn", Protocol: "foobar"}))
			Expect(idx.Names).To(HaveKey(ServiceProtocol{Name: "burn", Protocol: "baz"}))

			Expect(idx.Ports).To(HaveLen(3))
			Expect(idx.Ports).To(HaveKey(ServicePort{Port: 666, Protocol: ""}))
			Expect(idx.Ports).To(HaveKey(ServicePort{Port: 666, Protocol: "foobar"}))
			Expect(idx.Ports).To(HaveKey(ServicePort{Port: 666, Protocol: "baz"}))

			Expect(idx.ByName("frotz", "")).To(BeNil())
			Expect(idx.ByName("burn", "")).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Name":    Equal("crash"),
				"Port":    Equal(666),
				"Aliases": ConsistOf("burn"),
			})))

			Expect(idx.ByPort(12345, "")).To(BeNil())
			Expect(idx.ByPort(666, "")).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Name":    Equal("crash"),
				"Port":    Equal(666),
				"Aliases": ConsistOf("burn"),
			})))
		})

		It("merges indices", func() {
			s, err := ParseServices(strings.NewReader(`
crash 666/foobar
`), protos)
			Expect(err).NotTo(HaveOccurred())
			idx := NewServiceIndex(s)

			s, err = ParseServices(strings.NewReader(`
crash 666/baz
`), protos)
			Expect(err).NotTo(HaveOccurred())
			idx.MergeIndex(NewServiceIndex(s))

			Expect(idx.Names).To(HaveLen(3)) // sic! incl. zero protocol name
			Expect(idx.Ports).To(HaveLen(3)) // dto.
		})

	})

	Context("builtins", func() {

		BeforeEach(func() {
			Services = ServiceIndex{}
		})

		It("looks services up by name", func() {
			Expect(ServiceByName("domain", "udp")).NotTo(BeNil())
		})

		It("looks services up by port", func() {
			Expect(ServiceByPort(53, "tcp")).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Name": Equal("domain"),
			})))
		})

	})

})
