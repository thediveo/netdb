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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("protocols", func() {

	Context("parsing descriptions", func() {

		It("returns correct descriptions", func() {
			p, err := ParseProtocols(strings.NewReader(`
foobar 66
ratzfatz	123 schwuppdiwupp siebenmeilenstiefler # aber nicht rumpelpumpel
`))
			Expect(err).NotTo(HaveOccurred())
			Expect(p).To(ConsistOf(
				MatchFields(IgnoreExtras, Fields{
					"Name":    Equal("foobar"),
					"Number":  Equal(uint8(66)),
					"Aliases": BeEmpty(),
				}),
				MatchFields(IgnoreExtras, Fields{
					"Name":    Equal("ratzfatz"),
					"Number":  Equal(uint8(123)),
					"Aliases": ConsistOf("schwuppdiwupp", "siebenmeilenstiefler"),
				}),
			))
		})

		It("ignores comments and empty lines without errors", func() {
			p, err := ParseProtocols(strings.NewReader(`
# A comment
		# Another comment

foobar 66
`))
			Expect(err).NotTo(HaveOccurred())
			Expect(p).To(HaveLen(1))
		})

		It("silently skips malformed definitions", func() {
			p, err := ParseProtocols(strings.NewReader(`
foobar
`))
			Expect(err).NotTo(HaveOccurred())
			Expect(p).To(HaveLen(0))
		})

		It("reports invalid protocol definitions", func() {
			p, err := ParseProtocols(strings.NewReader(`
foobar 666
`))
			Expect(err).To(HaveOccurred())
			Expect(p).To(BeNil())
		})

		It("reports scanner errors", func() {
			f, err := os.Open("protocol_test.go")
			Expect(err).NotTo(HaveOccurred())
			f.Close() // sic! no defer!
			_, err = ParseProtocols(f)
			Expect(err).To(HaveOccurred())
		})

	})

	Context("loading", func() {

		It("loads protocol descriptions from file", func() {
			_, err := LoadProtocols("test/non-existing-protocols")
			Expect(err).To(HaveOccurred())

			idx, err := LoadProtocols("test/protocols")
			Expect(err).NotTo(HaveOccurred())
			Expect(idx.Names["ratzfatz"]).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Name":    Equal("ratzfatz"),
				"Number":  Equal(uint8(123)),
				"Aliases": ConsistOf("schwuppdiwupp", "siebenmeilenstiefler"),
			})))
		})

	})

	Context("indexing", func() {

		It("builds index", func() {
			p, err := ParseProtocols(strings.NewReader(`
ratzfatz	123 schwuppdiwupp siebenmeilenstiefler
`))
			Expect(err).NotTo(HaveOccurred())
			idx := NewProtocolIndex(p)
			Expect(idx.Names).To(HaveLen(3))
			Expect(idx.Names).To(HaveKey("ratzfatz"))
			Expect(idx.Names).To(HaveKey("schwuppdiwupp"))
			Expect(idx.Names).To(HaveKey("siebenmeilenstiefler"))
			Expect(idx.Numbers).To(HaveLen(1))
			Expect(idx.Numbers).To(HaveKey(uint8(123)))
		})

		It("merges indices", func() {
			p, err := ParseProtocols(strings.NewReader(`
ratzfatz	123 schwuppdiwupp siebenmeilenstiefler
`))
			Expect(err).NotTo(HaveOccurred())
			idx := NewProtocolIndex(p)

			p, err = ParseProtocols(strings.NewReader(`
foobar	66
`))
			Expect(err).NotTo(HaveOccurred())

			idx.MergeIndex(NewProtocolIndex(p))
			Expect(idx.Names).To(HaveLen(4))
			Expect(idx.Names).To(HaveKey("ratzfatz"))
			Expect(idx.Names).To(HaveKey("foobar"))
			Expect(idx.Numbers).To(HaveLen(2))
			Expect(idx.Numbers).To(HaveKey(uint8(66)))
		})

	})

	Context("builtins", func() {

		BeforeEach(func() {
			Protocols = ProtocolIndex{}
		})

		It("looks protocols up by name", func() {
			Expect(ProtocolByName("tcp")).NotTo(BeNil())
			Expect(ProtocolByName("udp")).NotTo(BeNil())
		})

		It("looks protocols up by number", func() {
			Expect(ProtocolByNumber(6)).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Name": Equal("tcp"),
			})))
			Expect(ProtocolByNumber(17)).To(PointTo(MatchFields(IgnoreExtras, Fields{
				"Name": Equal("udp"),
			})))
		})

	})

})
