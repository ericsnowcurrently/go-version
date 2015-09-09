// Copyright 2015 Eric Snow
// Licensed under the New BSD License, see LICENSE file for details.

package version_test

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/ericsnowcurrently/go-version/version"
)

var _ = gc.Suite(&buildSuite{})

type buildSuite struct {
	baseSuite
}

func (buildSuite) TestParseBuild(c *gc.C) {
	for _, test := range buildTests {
		test.checkParsing(c, "build")
	}
	for _, test := range releaseTests {
		if test.err != "" {
			continue
		}
		if test.inexact {
			continue
		}

		var expected version.Build
		expected.Release = test.expected.(version.Release)
		test := versionTest{
			vers:     test.vers,
			expected: expected,
			inexact:  true,
		}
		test.checkParsing(c, "build")
	}
	for _, test := range numberTests {
		if test.err != "" {
			continue
		}
		if test.inexact {
			continue
		}

		var expected version.Build
		expected.Number = test.expected.(version.Number)
		expected.Level = version.ReleaseFinal
		test := versionTest{
			vers:     test.vers,
			expected: expected,
			inexact:  true,
		}
		test.checkParsing(c, "build")
	}
}

func (buildSuite) TestStringOkay(c *gc.C) {
	build := newBuild(2, 3, 1, "a1", 2)
	builds := build.String()

	c.Check(builds, gc.Equals, "2.3.1-alpha1.2")
}

func (buildSuite) TestStringZeroValue(c *gc.C) {
	var build version.Build
	builds := build.String()

	c.Check(builds, gc.Equals, "0.0.0-dev")
}

func (buildSuite) TestIsZeroTrue(c *gc.C) {
	var build version.Build
	isZero := build.IsZero()

	c.Check(isZero, jc.IsTrue)
}

func (buildSuite) TestIsZeroFalse(c *gc.C) {
	build := newBuild(2, 3, 1, "a1", 2)
	isZero := build.IsZero()

	c.Check(isZero, jc.IsFalse)
}

func (buildSuite) TestCompare(c *gc.C) {
	var tests []cmpTest
	vers := "3.3.3b2.4"
	ver := newBuild(3, 3, 3, "b2", 4)

	// more
	for other, ok := ver.Prev(); ok; other, ok = other.Prev() {
		tests = append(tests, cmpTest{
			vers:     vers,
			others:   other.String(),
			expected: 1,
		})
	}
	// same
	tests = append(tests, cmpTest{
		vers:     vers,
		others:   vers,
		expected: 0,
	})
	// less
	for other, ok := ver.Next(7); ok; other, ok = other.Next(7) {
		tests = append(tests, cmpTest{
			vers:     vers,
			others:   other.String(),
			expected: -1,
		})
	}

	for _, test := range tests {
		test.run(c, "build")
	}
	for _, test := range cmpBuildsTests {
		test.run(c, "build")
	}
	for _, test := range cmpReleasesTests {
		test.run(c, "build")
	}
	for _, test := range cmpNumbersTests {
		test.run(c, "build")
	}
}

func (buildSuite) TestPrev(c *gc.C) {
	var builds []string
	ok := true
	for ver := newBuild(3, 3, 3, "a2", 5); ok; ver, ok = ver.Prev() {
		builds = append(builds, ver.String())
	}

	c.Check(builds, jc.DeepEquals, []string{
		"3.3.3-alpha2.5",
		"3.3.3-alpha2.4",
		"3.3.3-alpha2.3",
		"3.3.3-alpha2.2",
		"3.3.3-alpha2.1",
	})
}

func (buildSuite) TestNext(c *gc.C) {
	var releases []string
	ok := true
	for ver := newBuild(3, 3, 3, "a2", 1); ok; ver, ok = ver.Next(5) {
		releases = append(releases, ver.String())
	}

	c.Check(releases, jc.DeepEquals, []string{
		"3.3.3-alpha2.1",
		"3.3.3-alpha2.2",
		"3.3.3-alpha2.3",
		"3.3.3-alpha2.4",
		"3.3.3-alpha2.5",
	})
}

func (buildSuite) TestSerialization(c *gc.C) {
	for format := range marshallers {
		for _, test := range releaseTests {
			if test.err != "" {
				continue
			}
			test.checkSerialization(c, "build", format)
		}
	}
}

func (buildSuite) TestMarshalJSON(c *gc.C) {
	test := versionTest{
		vers:       "3.2.1a1",
		marshalled: `"3.2.1-alpha1"`,
	}

	test.checkMarshal(c, "build", "json")
}

func (buildSuite) TestUnMarshalJSON(c *gc.C) {
	test := versionTest{
		vers:       "3.2.1a1",
		marshalled: `"3.2.1-alpha1"`,
	}

	test.checkUnmarshal(c, "build", "json")
}

func (buildSuite) TestGetYAML(c *gc.C) {
	test := versionTest{
		vers:       "3.2.1a1",
		marshalled: `3.2.1-alpha1` + "\n",
	}

	test.checkMarshal(c, "build", "yaml")
}

func (buildSuite) TestSetYAML(c *gc.C) {
	test := versionTest{
		vers:       "3.2.1a1",
		marshalled: `3.2.1-alpha1`,
	}

	test.checkUnmarshal(c, "build", "yaml")
}

var buildTests = []versionTest{{
	vers:     "0.0.0.0",
	expected: newBuild(0, 0, 0, "", 0),
	inexact:  true,
	//dev:    true,
}, {
	vers:     "2.3.1.5",
	expected: newBuild(2, 3, 1, "", 5),
	inexact:  true,
	//dev:    true,
}, {
	vers:     "10.234.3456.64",
	expected: newBuild(10, 234, 3456, "", 64),
	inexact:  true,
	//dev:    true,
}, {
	vers:     "1.21-alpha1.1",
	expected: newBuild(1, 21, 0, "a1", 1),
	inexact:  true,
	//dev:    true,
}, {
	vers:     "1.21.1a1.1",
	expected: newBuild(1, 21, 1, "a1", 1),
	inexact:  true,
	//dev:    true,
}}

var cmpBuildsTests = []cmpTest{
	{"1.2-alpha2.1", "1.2-alpha2", 1},
	{"1.2-alpha2.2", "1.2-alpha2.1", 1},
	{"1.2-beta1", "1.2-alpha2.1", 1},
	{"2.0.0.0", "2.0.0", 0},
	{"2.0.0.0", "2.0.0.0", 0},
	{"2.0.0.1", "2.0.0.0", 1},
	{"2.0.1.10", "2.0.0.0", 1},
	// TODO(ericsnow) Support these?
	//{"2.0-_0", "2.0-00", 1},
	//{"2.0-_0", "2.0.0", -1},
	//{"2.0-_0", "2.0-alpha1.0", -1},
	//{"2.0-_0", "1.999.0", 1},
}
