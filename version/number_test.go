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

var _ = gc.Suite(&numberSuite{})

type numberSuite struct {
	baseSuite
}

func (numberSuite) TestParseNumber(c *gc.C) {
	var tests []versionTest

	max := newNumber(2, 2, 2)
	ok := true
	for num := newNumber(0, 0, 0); ok; num, ok = num.Next(max) {
		tests = append(tests, versionTest{
			vers:     fmt.Sprintf("%d.%d.%d", num.Major, num.Minor, num.Micro),
			expected: num,
			//dev: true, // ...if major == 0
		})
	}

	for _, test := range tests {
		test.checkParsing(c, "number")
	}
	for _, test := range numberTests {
		test.checkParsing(c, "number")
	}
}

func (numberSuite) TestStringOkay(c *gc.C) {
	num := version.Number{2, 3, 1}
	nums := num.String()

	c.Check(nums, gc.Equals, "2.3.1")
}

func (numberSuite) TestStringZeroValue(c *gc.C) {
	var num version.Number
	nums := num.String()

	c.Check(nums, gc.Equals, "0.0.0")
}

func (numberSuite) TestStringShort(c *gc.C) {
	num := version.Number{2, 3, 0}
	nums := num.String()

	c.Check(nums, gc.Equals, "2.3.0")
}

func (numberSuite) TestFeature(c *gc.C) {
	num := version.Number{2, 3, 1}
	feature := num.Feature()

	c.Check(feature, gc.Equals, "2.3")
}

func (numberSuite) TestIsZeroTrue(c *gc.C) {
	var num version.Number
	isZero := num.IsZero()

	c.Check(isZero, jc.IsTrue)
}

func (numberSuite) TestIsZeroFalse(c *gc.C) {
	num := version.Number{2, 3, 1}
	isZero := num.IsZero()

	c.Check(isZero, jc.IsFalse)
}

func (numberSuite) TestValidateOkay(c *gc.C) {
	num := newNumber(3, 3, 3)
	err := num.Validate()

	c.Check(err, jc.ErrorIsNil)
}

func (numberSuite) TestValidateZeroValue(c *gc.C) {
	num := newNumber(0, 0, 0)
	err := num.Validate()

	c.Check(err, jc.Satisfies, errors.IsNotValid)
}

func (numberSuite) TestCompare(c *gc.C) {
	var tests []cmpTest
	vers := "3.3.3"
	ver := newNumber(3, 3, 3)

	// more
	max := newNumber(4, 4, 4)
	for other, ok := ver.Prev(max); ok; other, ok = other.Prev(max) {
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
	max = newNumber(5, 5, 5)
	for other, ok := ver.Next(max); ok; other, ok = other.Next(max) {
		tests = append(tests, cmpTest{
			vers:     vers,
			others:   other.String(),
			expected: -1,
		})
	}

	for _, test := range tests {
		test.run(c, "number")
	}
	for _, test := range cmpNumbersTests {
		test.run(c, "number")
	}
}

func (numberSuite) TestPrev(c *gc.C) {
	var got []string
	max := newNumber(0, 4, 4)
	ok := true
	for num := newNumber(2, 3, 1); ok; num, ok = num.Prev(max) {
		got = append(got, num.String())
	}

	c.Check(got, jc.DeepEquals, []string{
		"2.3.1", "2.3.0",
		"2.2.4", "2.2.3", "2.2.2", "2.2.1", "2.2.0",
		"2.1.4", "2.1.3", "2.1.2", "2.1.1", "2.1.0",
		"2.0.4", "2.0.3", "2.0.2", "2.0.1", "2.0.0",

		"1.4.4", "1.4.3", "1.4.2", "1.4.1", "1.4.0",
		"1.3.4", "1.3.3", "1.3.2", "1.3.1", "1.3.0",
		"1.2.4", "1.2.3", "1.2.2", "1.2.1", "1.2.0",
		"1.1.4", "1.1.3", "1.1.2", "1.1.1", "1.1.0",
		"1.0.4", "1.0.3", "1.0.2", "1.0.1", "1.0.0",

		"0.4.4", "0.4.3", "0.4.2", "0.4.1", "0.4.0",
		"0.3.4", "0.3.3", "0.3.2", "0.3.1", "0.3.0",
		"0.2.4", "0.2.3", "0.2.2", "0.2.1", "0.2.0",
		"0.1.4", "0.1.3", "0.1.2", "0.1.1", "0.1.0",
		"0.0.4", "0.0.3", "0.0.2", "0.0.1", "0.0.0",
	})
}

func (numberSuite) TestNext(c *gc.C) {
	var got []string
	max := newNumber(4, 4, 4)
	ok := true
	for num := newNumber(2, 3, 1); ok; num, ok = num.Next(max) {
		got = append(got, num.String())
	}

	c.Check(got, jc.DeepEquals, []string{
		"2.3.1", "2.3.2", "2.3.3", "2.3.4",
		"2.4.0", "2.4.1", "2.4.2", "2.4.3", "2.4.4",

		"3.0.0", "3.0.1", "3.0.2", "3.0.3", "3.0.4",
		"3.1.0", "3.1.1", "3.1.2", "3.1.3", "3.1.4",
		"3.2.0", "3.2.1", "3.2.2", "3.2.3", "3.2.4",
		"3.3.0", "3.3.1", "3.3.2", "3.3.3", "3.3.4",
		"3.4.0", "3.4.1", "3.4.2", "3.4.3", "3.4.4",

		"4.0.0", "4.0.1", "4.0.2", "4.0.3", "4.0.4",
		"4.1.0", "4.1.1", "4.1.2", "4.1.3", "4.1.4",
		"4.2.0", "4.2.1", "4.2.2", "4.2.3", "4.2.4",
		"4.3.0", "4.3.1", "4.3.2", "4.3.3", "4.3.4",
		"4.4.0", "4.4.1", "4.4.2", "4.4.3", "4.4.4",
	})
}

func (numberSuite) TestSerialization(c *gc.C) {
	for format := range marshallers {
		for _, test := range numberTests {
			if test.err != "" {
				continue
			}
			test.checkSerialization(c, "number", format)
		}
	}
}

func (numberSuite) TestMarshalJSON(c *gc.C) {
	test := versionTest{
		vers:       "3.2.1",
		marshalled: `"3.2.1"`,
	}

	test.checkMarshal(c, "number", "json")
}

func (numberSuite) TestUnMarshalJSON(c *gc.C) {
	test := versionTest{
		vers:       "3.2.1",
		marshalled: `"3.2.1"`,
	}

	test.checkUnmarshal(c, "number", "json")
}

func (numberSuite) TestGetYAML(c *gc.C) {
	test := versionTest{
		vers:       "3.2.1",
		marshalled: `3.2.1` + "\n",
	}

	test.checkMarshal(c, "number", "yaml")
}

func (numberSuite) TestSetYAML(c *gc.C) {
	test := versionTest{
		vers:       "3.2.1",
		marshalled: `3.2.1`,
	}

	test.checkUnmarshal(c, "number", "yaml")
}

var cmpNumbersTests = []cmpTest{
	{"1.0.0", "1.0.0", 0},
	{"10.0.0", "9.0.0", 1},
	{"1.0.0", "1.0.1", -1},
	{"1.0.1", "1.0.0", 1},
	{"1.0.0", "1.1.0", -1},
	{"1.1.0", "1.0.0", 1},
	{"1.0.0", "2.0.0", -1},
	{"1.2.1", "1.2.0", 1},
	{"2.0.0", "1.0.0", 1},
}

var numberTests = []versionTest{{
	vers:     "0.0.0",
	expected: newNumber(0, 0, 0),
	//dev:    true,
}, {
	vers:     "0.0.1",
	expected: newNumber(0, 0, 1),
	//dev:    true,
}, {
	vers:     "0.3.0",
	expected: newNumber(0, 3, 0),
	//dev:    true,
}, {
	vers:     "0.3.1",
	expected: newNumber(0, 3, 1),
	//dev:    true,
}, {
	vers:     "2.0.0",
	expected: newNumber(2, 0, 0),
}, {
	vers:     "2.3.1",
	expected: newNumber(2, 3, 1),
}, {
	vers:     "2.3",
	expected: newNumber(2, 3, 0),
	inexact:  true,
}, {
	vers:     "2",
	expected: newNumber(2, 0, 0),
	inexact:  true,
}, {
	vers:     "10.234.3456",
	expected: newNumber(10, 234, 3456),
}, {
	vers:     "10.234.3456.1",
	expected: newNumber(10, 234, 3456),
	inexact:  true,
}, {
	vers:     "1.21-alpha1",
	expected: newNumber(1, 21, 0),
	inexact:  true,
}, {
	vers:     "1.21-alpha1.1",
	expected: newNumber(1, 21, 0),
	inexact:  true,
}, {
	vers: "0.2.",
	err:  `version string .* not valid`,
}, {
	vers: "0.2.spam",
	err:  `version string .* not valid`,
}, {
	vers: "0.2..1",
	err:  `version string .* not valid`,
}}
