// Copyright 2015 Eric Snow
// Licensed under the New BSD License, see LICENSE file for details.

package version_test

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	goyaml "gopkg.in/yaml.v1"

	"github.com/ericsnowcurrently/go-version/version"
)

type baseSuite struct {
	testing.IsolationSuite
}

func newNumber(major, minor, micro uint) version.Number {
	return version.Number{
		Major: major,
		Minor: minor,
		Micro: micro,
	}
}

func newRelease(major, minor, micro uint, release string) version.Release {
	var level version.ReleaseLevel
	var serial uint
	switch {
	case release == "dev":
		level = version.ReleaseDevelopment
	case strings.HasPrefix(release, "a"):
		level = version.ReleaseAlpha
		fmt.Sscanf(release, "a%d", &serial)
	case strings.HasPrefix(release, "b"):
		level = version.ReleaseBeta
		fmt.Sscanf(release, "b%d", &serial)
	case strings.HasPrefix(release, "rc"):
		level = version.ReleaseCandidate
		fmt.Sscanf(release, "rc%d", &serial)
	case release == "":
		level = version.ReleaseFinal
	default:
		panic(fmt.Sprintf("unrecognized release %q", release))
	}

	return version.Release{
		Number: newNumber(major, minor, micro),
		Level:  level,
		Serial: serial,
	}
}

func newBuild(major, minor, micro uint, release string, build uint) version.Build {
	return version.Build{
		Release: newRelease(major, minor, micro, release),
		Index:   build,
	}
}

type cmpTest struct {
	vers     string
	others   string
	expected int
}

func (t cmpTest) numbers(c *gc.C) (version.Number, version.Number) {
	ver, _, err := version.ParseNumber(t.vers)
	c.Assert(err, jc.ErrorIsNil)
	other, _, err := version.ParseNumber(t.others)
	c.Assert(err, jc.ErrorIsNil)
	return ver, other
}

func (t cmpTest) compareNumbers(c *gc.C) (int, int) {
	ver, other := t.numbers(c)
	compareVer := ver.Compare(other)
	// Check that reversing the operands has
	// the expected result.
	compareOther := other.Compare(ver)
	return compareVer, compareOther
}

func (t cmpTest) releases(c *gc.C) (version.Release, version.Release) {
	ver, _, err := version.ParseRelease(t.vers)
	c.Assert(err, jc.ErrorIsNil)
	other, _, err := version.ParseRelease(t.others)
	c.Assert(err, jc.ErrorIsNil)
	return ver, other
}

func (t cmpTest) compareReleases(c *gc.C) (int, int) {
	ver, other := t.releases(c)
	compareVer := ver.Compare(other)
	// Check that reversing the operands has
	// the expected result.
	compareOther := other.Compare(ver)
	return compareVer, compareOther
}

func (t cmpTest) builds(c *gc.C) (version.Build, version.Build) {
	ver, _, err := version.ParseBuild(t.vers)
	c.Assert(err, jc.ErrorIsNil)
	other, _, err := version.ParseBuild(t.others)
	c.Assert(err, jc.ErrorIsNil)
	return ver, other
}

func (t cmpTest) compareBuilds(c *gc.C) (int, int) {
	ver, other := t.builds(c)
	compareVer := ver.Compare(other)
	// Check that reversing the operands has
	// the expected result.
	compareOther := other.Compare(ver)
	return compareVer, compareOther
}

func (t cmpTest) run(c *gc.C, kind string) {
	c.Logf("- testing %q <> %q -> %d", t.vers, t.others, t.expected)

	var compareVer, compareOther int
	switch kind {
	case "number":
		compareVer, compareOther = t.compareNumbers(c)
	case "release":
		compareVer, compareOther = t.compareReleases(c)
	case "build":
		compareVer, compareOther = t.compareBuilds(c)
	default:
		c.Logf("unknown kind %q", kind)
		c.FailNow()
	}

	c.Check(compareVer, gc.Equals, t.expected)
	c.Check(compareOther, gc.Equals, -t.expected)
}

type versionTest struct {
	vers       string
	expected   interface{}
	inexact    bool
	marshalled string
	err        string
}

func (t versionTest) zero(c *gc.C, kind string) interface{} {
	switch kind {
	case "number":
		return &version.Number{}
	case "release":
		return &version.Release{}
	case "build":
		return &version.Build{}
	default:
		c.Logf("unknown kind %q", kind)
		c.FailNow()
		return nil
	}
}

func (t versionTest) parsed(c *gc.C, kind string) interface{} {
	switch kind {
	case "number":
		ver, _, err := version.ParseNumber(t.vers)
		c.Assert(err, jc.ErrorIsNil)
		return &ver
	case "release":
		ver, _, err := version.ParseRelease(t.vers)
		c.Assert(err, jc.ErrorIsNil)
		return &ver
	case "build":
		ver, _, err := version.ParseBuild(t.vers)
		c.Assert(err, jc.ErrorIsNil)
		return &ver
	default:
		c.Logf("unknown kind %q", kind)
		c.FailNow()
		return nil
	}
}

func (t versionTest) marshaller(c *gc.C, format string) marshaller {
	marshaller, ok := marshallers[format]
	if !ok {
		c.Logf("unknown format %q", format)
		c.FailNow()
	}
	return marshaller
}

func (t versionTest) checkParsing(c *gc.C, kind string) {
	c.Logf("- testing (%s) parsing of %q", kind, t.vers)

	var ver fmt.Stringer
	var err error
	switch kind {
	case "number":
		ver, _, err = version.ParseNumber(t.vers)
	case "release":
		ver, _, err = version.ParseRelease(t.vers)
	case "build":
		ver, _, err = version.ParseBuild(t.vers)
	default:
		c.Logf("unknown kind %q", kind)
		c.FailNow()
	}

	if t.err == "" {
		if !c.Check(err, jc.ErrorIsNil) {
			return
		}
		c.Check(ver, gc.Equals, t.expected)
		if !t.inexact {
			// Check the round-trip.
			c.Check(ver.String(), gc.Equals, t.vers)
		}
	} else {
		c.Check(err, gc.ErrorMatches, t.err)
	}
}

func (t versionTest) checkSerialization(c *gc.C, kind, format string) {
	c.Logf("- testing %s round-trip for %q (%s)", format, t.vers, kind)
	c.Assert(t.err, gc.Equals, "")

	marshaller := t.marshaller(c, format)
	original := t.parsed(c, kind)
	data, err := marshaller.marshal(original)
	c.Assert(err, jc.ErrorIsNil)

	final := t.zero(c, kind)
	err = marshaller.unmarshal(data, final)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(final, jc.DeepEquals, original)
}

func (t versionTest) checkMarshal(c *gc.C, kind, format string) {
	c.Logf("- testing %s marshalling for %q (%s)", format, t.vers, kind)
	c.Assert(t.err, gc.Equals, "")

	marshaller := t.marshaller(c, format)
	ver := t.parsed(c, kind)
	data, err := marshaller.marshal(ver)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(string(data), gc.Equals, t.marshalled)
	c.Assert(t.err, gc.Equals, "")
}

func (t versionTest) checkUnmarshal(c *gc.C, kind, format string) {
	c.Logf("- testing %s unmarshalling for %q (%s)", format, t.vers, kind)
	c.Assert(t.err, gc.Equals, "")
	if t.expected == nil {
		t.expected = t.parsed(c, kind)
	}

	marshaller := t.marshaller(c, format)
	ver := t.zero(c, kind)
	err := marshaller.unmarshal([]byte(t.marshalled), ver)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(ver, jc.DeepEquals, t.expected)
}

type marshaller struct {
	marshal   func(interface{}) ([]byte, error)
	unmarshal func([]byte, interface{}) error
}

var marshallers = map[string]marshaller{
	"json": marshaller{
		json.Marshal,
		json.Unmarshal,
	},
	"yaml": marshaller{
		goyaml.Marshal,
		goyaml.Unmarshal,
		// TODO(ericsnow) Work around goyaml bug? (#1096149)
		// (SetYAML is not called for non-pointer fields.)
	},
}
