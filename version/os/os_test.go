// Copyright 2015 Eric Snow
// Licensed under the New BSD License, see LICENSE file for details.

package os_test

import (
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/ericsnowcurrently/go-version/version/os"
)

var _ = gc.Suite(&osSuite{})

type osSuite struct {
	testing.IsolationSuite
}

func (osSuite) TestIsUnixKnown(c *gc.C) {
	for _, known := range os.Unix {
		c.Logf("checking %q", known)
		isUnix := os.IsUnix(known)

		c.Check(isUnix, jc.IsTrue)
	}
}

func (osSuite) TestIsUnixWindows(c *gc.C) {
	isUnix := os.IsUnix("windows")

	c.Check(isUnix, jc.IsFalse)
}

func (osSuite) TestIsUnixZeroValue(c *gc.C) {
	isUnix := os.IsUnix("")

	c.Check(isUnix, jc.IsFalse)
}

func (osSuite) TestIsUnixUnknown(c *gc.C) {
	isUnix := os.IsUnix("<unknown OS>")

	c.Check(isUnix, jc.IsFalse)
}
