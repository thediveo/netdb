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
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("ethertypes", func() {

	Context("parsing descriptions", func() {

		It("returns correct descriptions", func() {
			p, err := ParseEtherTypes(strings.NewReader(`
test 9000
RoMON	88BF	mikrotik-rommon mt-rommon		# MikroTik RoMON (unofficial)
`))
			Expect(err).NotTo(HaveOccurred())
			Expect(p).To(ConsistOf(
				MatchFields(IgnoreExtras, Fields{
					"Name":    Equal("test"),
					"Number":  Equal(uint16(0x9000)),
					"Aliases": BeEmpty(),
					"Comment": BeEmpty(),
				}),
				MatchFields(IgnoreExtras, Fields{
					"Name":    Equal("RoMON"),
					"Number":  Equal(uint16(0x88BF)),
					"Aliases": ConsistOf("mikrotik-rommon", "mt-rommon"),
					"Comment": Equal("MikroTik RoMON (unofficial)"),
				}),
			))
		})
		It("ignores comments and empty lines without errors", func() {
			p, err := ParseEtherTypes(strings.NewReader(`
# A comment
		# Another comment

foobar 66
`))
			Expect(err).NotTo(HaveOccurred())
			Expect(p).To(HaveLen(1))
		})
		It("silently skips malformed definitions", func() {
			p, err := ParseEtherTypes(strings.NewReader(`
foobar
			`))
			Expect(err).NotTo(HaveOccurred())
			Expect(p).To(HaveLen(0))
		})
		It("reports invalid definitions", func() {
			p, err := ParseEtherTypes(strings.NewReader(`
foobar 666x
`))
			Expect(err).To(HaveOccurred())
			Expect(p).To(BeNil())
		})
		It("reports scanner errors", func() {
			f, err := os.Open("ethertype_test.go")
			Expect(err).NotTo(HaveOccurred())
			f.Close() // sic! no defer!
			_, err = ParseEtherTypes(f)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("loading", func() {

		It("loads EtherType descriptions from file", func() {
			_, err := LoadEtherTypes("test/non-existing-ethertypes")
			Expect(err).To(HaveOccurred())

			idx, err := LoadEtherTypes("test/ethertypes")
			Expect(err).NotTo(HaveOccurred())
			Expect(idx.Names["test"]).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Name":   Equal("test"),
				"Number": Equal(uint16(0x9000)),
			})))
		})
	})

	Context("indexing", func() {

		It("builds index", func() {
			p, err := ParseEtherTypes(strings.NewReader(`
RoMON	88BF	mikrotik-rommon mt-rommon		# MikroTik RoMON (unofficial)
`))
			Expect(err).NotTo(HaveOccurred())
			idx := NewEtherTypeIndex(p)
			Expect(idx.Names).To(HaveLen(3))
			Expect(idx.Names).To(HaveKey("RoMON"))
			Expect(idx.Names).To(HaveKey("mikrotik-rommon"))
			Expect(idx.Names).To(HaveKey("mt-rommon"))
			Expect(idx.Numbers).To(HaveLen(1))
			Expect(idx.Numbers).To(HaveKey(uint16(0x88BF)))
		})
		It("merges indices", func() {
			p, err := ParseEtherTypes(strings.NewReader(`
RoMON	88BF	mikrotik-rommon mt-rommon		# MikroTik RoMON (unofficial)
`))
			Expect(err).NotTo(HaveOccurred())
			idx := NewEtherTypeIndex(p)

			p, err = ParseEtherTypes(strings.NewReader(`
foobar	66
`))
			Expect(err).NotTo(HaveOccurred())

			idx.MergeIndex(NewEtherTypeIndex(p))
			Expect(idx.Names).To(HaveLen(4))
			Expect(idx.Names).To(HaveKey("RoMON"))
			Expect(idx.Names).To(HaveKey("foobar"))
			Expect(idx.Numbers).To(HaveLen(2))
			Expect(idx.Numbers).To(HaveKey(uint16(0x66)))
		})
	})

	Context("builtins", func() {

		It("looks EtherTypes up by name", func() {
			Expect(EtherTypeByName("IPv4")).NotTo(BeNil())
			Expect(EtherTypeByName("ip")).NotTo(BeNil())
			Expect(EtherTypeByName("IPX")).NotTo(BeNil())
		})
		It("looks EtherTypes up by number", func() {
			Expect(EtherTypeByNumber(2048)).NotTo(BeNil())
			Expect(EtherTypeByNumber(0x800)).NotTo(BeNil())
		})

	})

})
