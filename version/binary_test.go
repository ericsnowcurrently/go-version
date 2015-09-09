// Copyright 2015 Eric Snow
// Licensed under the New BSD License, see LICENSE file for details.

package version_test

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/ericsnowcurrently/go-version/version"
)

var _ = gc.Suite(&binarySuite{})

type binarySuite struct {
	baseSuite
}

func expectedBinary(c *gc.C, test versionTest) version.Binary {
	var expected version.Binary

	switch orig := test.expected.(type) {
	case version.Binary:
		expected = orig
	// TODO(ericsnow) Move the rest to expectedBuild, etc.
	case version.Build:
		expected.Build = orig
	case version.Release:
		expected.Release = orig
	case version.Number:
		expected.Number = orig
		expected.Level = version.ReleaseFinal
	default:
		c.Logf("unknown version type %T", orig)
		c.FailNow()
	}

	return expected
}

func (binarySuite) TestParseBinary(c *gc.C) {
	var tests []versionTest

	addTest := func(test versionTest) {
		if test.err != "" {
			return
		}
		if test.inexact {
			return
		}
		expected := expectedBinary(c, test)

		expected.Series = "trusty"
		expected.Arch = "amd64"
		tests = append(tests, versionTest{
			vers:     test.vers + "-trusty-amd64",
			expected: expected,
			inexact:  true,
		})
	}

	for _, test := range binaryTests {
		tests = append(tests, test)
	}
	for _, test := range buildTests {
		addTest(test)
	}
	for _, test := range releaseTests {
		addTest(test)
	}
	for _, test := range numberTests {
		addTest(test)
	}

	for _, test := range tests {
		test.checkParsing(c, "binary")
	}
}

func (binarySuite) TestStringOkay(c *gc.C) {
	bin := newBinary(2, 3, 1, "a1", 2, "trusty", "amd64")
	bins := bin.String()

	c.Check(bins, gc.Equals, "2.3.1-alpha1.2-trusty-amd64")
}

func (binarySuite) TestStringZeroValue(c *gc.C) {
	var bin version.Binary
	bins := bin.String()

	c.Check(bins, gc.Equals, "0.0.0-dev-unknown-unknown")
}

func (binarySuite) TestIsZeroTrue(c *gc.C) {
	var bin version.Binary
	isZero := bin.IsZero()

	c.Check(isZero, jc.IsTrue)
}

func (binarySuite) TestIsZeroFalse(c *gc.C) {
	bin := newBinary(2, 3, 1, "a1", 2, "trusty", "amd64")
	isZero := bin.IsZero()

	c.Check(isZero, jc.IsFalse)
}

func (binarySuite) TestCompare(c *gc.C) {
	var tests []cmpTest
	vers := "3.3.3b2.4"
	ver := newBinary(3, 3, 3, "b2", 4, "trusty", "amd64")

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

var binaryTests = []versionTest{{
	vers:     "1.2.3-trusty-amd64",
	expected: newBinary(1, 2, 3, "", 0, "trusty", "amd64"),
	inexact:  true,
}, {
	vers:     "1.2.3.4-trusty-amd64",
	expected: newBinary(1, 2, 3, "", 4, "trusty", "amd64"),
	inexact:  true,
	//dev: true,
}, {
	vers:     "1.2-alpha3-trusty-amd64",
	expected: newBinary(1, 2, 0, "a3", 0, "trusty", "amd64"),
	inexact:  true,
}, {
	vers:     "1.2.3-alpha3-trusty-amd64",
	expected: newBinary(1, 2, 3, "a3", 0, "trusty", "amd64"),
}, {
	vers:     "1.2.3-alpha3.4-trusty-amd64",
	expected: newBinary(1, 2, 3, "a3", 4, "trusty", "amd64"),
	//dev: true,
}, {
	vers: "1.2.3",
	err:  "binary version .* not valid",
}, {
	vers: "1.2-beta1",
	err:  "binary version .* not valid",
}, {
	vers: "1.2.3--amd64",
	err:  "binary version .* not valid",
}, {
	vers: "1.2.3-trusty-",
	err:  "binary version .* not valid",
}}
