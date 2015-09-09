// Copyright 2015 Eric Snow
// Licensed under the New BSD License, see LICENSE file for details.

package version_test

import (
	"fmt"

	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/ericsnowcurrently/go-version/version"
)

var (
	_ = gc.Suite(&releaseSuite{})
	_ = gc.Suite(&releaseLevelSuite{})
)

type releaseSuite struct {
	baseSuite
}

func (releaseSuite) TestParseReleaseOkay(c *gc.C) {
	for _, test := range releaseTests {
		test.checkParsing(c, "release")
	}
}

func (releaseSuite) TestParseReleaseLevels(c *gc.C) {
	var i uint
	for _, test := range numberTests {
		if test.err != "" {
			continue
		}
		if test.inexact {
			continue
		}

		num := test.expected.(version.Number)
		vers := test.vers
		for level, name := range version.ReleaseLevelNames {
			if _, ok := version.ReleaseLevelSingletons[level]; ok {
				test = versionTest{
					vers: vers + "-" + name,
					expected: version.Release{
						Number: num,
						Level:  level,
						Serial: 0,
					},
				}
				test.checkParsing(c, "release")
				continue
			}
			for i = 1; i < 4; i += 1 {
				test = versionTest{
					vers: vers + fmt.Sprintf("-%s%d", name, i),
					expected: version.Release{
						Number: num,
						Level:  level,
						Serial: i,
					},
				}
				test.checkParsing(c, "release")
			}
		}
	}
}

func (releaseSuite) TestParseReleaseAbbrev(c *gc.C) {
	var i uint
	for _, test := range numberTests {
		if test.err != "" {
			continue
		}
		if test.inexact {
			continue
		}

		num := test.expected.(version.Number)
		vers := test.vers
		for level, abbrev := range version.ReleaseLevelAbbrevs {
			for i = 1; i < 4; i += 1 {
				test = versionTest{
					vers: vers + fmt.Sprintf("%s%d", abbrev, i),
					expected: version.Release{
						Number: num,
						Level:  level,
						Serial: i,
					},
					inexact: true,
				}
				test.checkParsing(c, "release")
			}
		}
	}
}

func (releaseSuite) TestParseReleaseImplicitFinal(c *gc.C) {
	for _, test := range numberTests {
		if test.err != "" {
			continue
		}
		if test.inexact {
			continue
		}

		var expected version.Release
		expected.Number = test.expected.(version.Number)
		expected.Level = version.ReleaseFinal
		test = versionTest{
			vers:     test.vers,
			expected: expected,
			inexact:  true,
		}
		test.checkParsing(c, "release")
	}
}

func (releaseSuite) TestParseReleaseBadSerial(c *gc.C) {
	for level, name := range version.ReleaseLevelNames {
		c.Logf("checking %q", name)
		serial := 0
		if _, ok := version.ReleaseLevelSingletons[level]; ok {
			serial = 1
		}
		rels := fmt.Sprintf("2.3.1-%s%d", name, serial)
		_, _, err := version.ParseRelease(rels)

		c.Check(err, jc.Satisfies, errors.IsNotValid)
	}

}

func (releaseSuite) TestParseReleaseMajorOnly(c *gc.C) {
	invalid := []string{
		"1-alpha1",
		"2rc3",
		"3",
	}
	for _, rels := range invalid {
		c.Logf("checking %q", rels)
		_, _, err := version.ParseRelease(rels)

		c.Check(err, jc.Satisfies, errors.IsNotValid)
	}
}

func (releaseSuite) TestStringOkay(c *gc.C) {
	rel := newRelease(2, 3, 1, "a1")
	rels := rel.String()

	c.Check(rels, gc.Equals, "2.3.1-alpha1")
}

func (releaseSuite) TestStringZeroValue(c *gc.C) {
	var rel version.Release
	rels := rel.String()

	c.Check(rels, gc.Equals, "0.0.0-dev")
}

func (releaseSuite) TestAbbrev(c *gc.C) {
	rel := newRelease(2, 3, 1, "a1")
	rels := rel.Abbrev()

	c.Check(rels, gc.Equals, "2.3.1a1")
}

func (releaseSuite) TestIsZeroTrue(c *gc.C) {
	var rel version.Release
	isZero := rel.IsZero()

	c.Check(isZero, jc.IsTrue)
}

func (releaseSuite) TestIsZeroFalse(c *gc.C) {
	rel := newRelease(2, 3, 1, "a1")
	isZero := rel.IsZero()

	c.Check(isZero, jc.IsFalse)
}

func (releaseSuite) TestCompare(c *gc.C) {
	var tests []cmpTest
	vers := "3.3.3b2"
	ver := newRelease(3, 3, 3, "b2")

	// more
	for other, ok := ver.Prev(3, 2, 2); ok; other, ok = other.Prev(3, 2, 2) {
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
	for other, ok := ver.Next(3, 2, 2); ok; other, ok = other.Next(3, 2, 2) {
		tests = append(tests, cmpTest{
			vers:     vers,
			others:   other.String(),
			expected: -1,
		})
	}

	for _, test := range tests {
		test.run(c, "release")
	}
	for _, test := range cmpReleasesTests {
		test.run(c, "release")
	}
	for _, test := range cmpNumbersTests {
		test.run(c, "release")
	}
}

func (releaseSuite) TestPrev(c *gc.C) {
	var releases []string
	ok := true
	for ver := newRelease(3, 3, 3, ""); ok; ver, ok = ver.Prev(3, 2, 2) {
		releases = append(releases, ver.String())
	}

	c.Check(releases, jc.DeepEquals, []string{
		"3.3.3-final",
		"3.3.3-candidate2",
		"3.3.3-candidate1",
		"3.3.3-beta2",
		"3.3.3-beta1",
		"3.3.3-alpha3",
		"3.3.3-alpha2",
		"3.3.3-alpha1",
		"3.3.3-dev",
	})
}

func (releaseSuite) TestNext(c *gc.C) {
	var releases []string
	ok := true
	for ver := newRelease(3, 3, 3, "dev"); ok; ver, ok = ver.Next(3, 2, 2) {
		releases = append(releases, ver.String())
	}

	c.Check(releases, jc.DeepEquals, []string{
		"3.3.3-dev",
		"3.3.3-alpha1",
		"3.3.3-alpha2",
		"3.3.3-alpha3",
		"3.3.3-beta1",
		"3.3.3-beta2",
		"3.3.3-candidate1",
		"3.3.3-candidate2",
		"3.3.3-final",
	})
}

func (releaseSuite) TestSerialization(c *gc.C) {
	for format := range marshallers {
		for _, test := range releaseTests {
			if test.err != "" {
				continue
			}
			test.checkSerialization(c, "release", format)
		}
	}
}

func (releaseSuite) TestMarshalJSON(c *gc.C) {
	test := versionTest{
		vers:       "3.2.1a1",
		marshalled: `"3.2.1-alpha1"`,
	}

	test.checkMarshal(c, "release", "json")
}

func (releaseSuite) TestUnMarshalJSON(c *gc.C) {
	test := versionTest{
		vers:       "3.2.1a1",
		marshalled: `"3.2.1-alpha1"`,
	}

	test.checkUnmarshal(c, "release", "json")
}

func (releaseSuite) TestGetYAML(c *gc.C) {
	test := versionTest{
		vers:       "3.2.1a1",
		marshalled: `3.2.1-alpha1` + "\n",
	}

	test.checkMarshal(c, "release", "yaml")
}

func (releaseSuite) TestSetYAML(c *gc.C) {
	test := versionTest{
		vers:       "3.2.1a1",
		marshalled: `3.2.1-alpha1`,
	}

	test.checkUnmarshal(c, "release", "yaml")
}

type releaseLevelSuite struct {
	baseSuite
}

func (releaseLevelSuite) TestReleaseLevelFromName(c *gc.C) {
	for expected, name := range version.ReleaseLevelNames {
		c.Logf("checking %q", name)
		level := version.ReleaseLevelFromName(name)

		c.Check(level, gc.Equals, expected)
	}
}

func (releaseLevelSuite) TestReleaseLevelFromAbbrev(c *gc.C) {
	for expected, abbrev := range version.ReleaseLevelAbbrevs {
		c.Logf("checking %q", abbrev)
		level := version.ReleaseLevelFromAbbrev(abbrev)

		c.Check(level, gc.Equals, expected)
	}
}

func (releaseLevelSuite) TestString(c *gc.C) {
	for level, expected := range version.ReleaseLevelNames {
		c.Logf("checking %q", expected)
		name := level.String()

		c.Check(name, gc.Equals, expected)
	}
}

func (releaseLevelSuite) TestAbbrev(c *gc.C) {
	for level, name := range version.ReleaseLevelNames {
		c.Logf("checking %q", name)
		expected := version.ReleaseLevelAbbrevs[level]
		abbrev := level.Abbrev()

		c.Check(abbrev, gc.Equals, expected)
	}
}

func (releaseLevelSuite) TestIndex(c *gc.C) {
	for level, name := range version.ReleaseLevelNames {
		c.Logf("checking %q", name)
		expected := int(level) - 1
		index := level.Index()

		c.Check(index, gc.Equals, expected)
	}
}

var releaseTests = []versionTest{{
	vers:     "0.0.0",
	expected: newRelease(0, 0, 0, ""),
	inexact:  true,
	//dev:    true,
}, {
	vers:     "1.21.1-alpha1",
	expected: newRelease(1, 21, 1, "a1"),
	//dev:    true,
}, {
	vers:     "1.21.1a1",
	expected: newRelease(1, 21, 1, "a1"),
	inexact:  true,
	//dev:    true,
}, {
	vers:     "1.21-alpha1",
	expected: newRelease(1, 21, 0, "a1"),
	inexact:  true,
	//dev:    true,
}, {
	vers:     "1.21-alpha1.1",
	expected: newRelease(1, 21, 0, "a1"),
	inexact:  true,
	//dev:    true,
}, {
	vers: "1.21.alpha1",
	err:  "release version.* not valid",
}, {
	vers: "1.21.1alpha1",
	err:  "release version.* not valid",
}, {
	vers: "1.21-alpha",
	err:  "release version.* not valid",
}, {
	vers: "1.21-alpha1beta",
	err:  "release version.* not valid",
}, {
	vers: "1.21-alpha-dev",
	err:  "release version.* not valid",
}}

var cmpReleasesTests = []cmpTest{
	{"1.2-alpha1", "1.2.0", -1},
	{"1.2-alpha2", "1.2-alpha1", 1},
	{"1.2-beta1", "1.2-alpha1", 1},
	{"1.2-beta1", "1.2.0", -1},
}
