// Copyright 2015 Eric Snow
// Licensed under the New BSD License, see LICENSE file for details.

package os_test

import (
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/ericsnowcurrently/go-version/version/os"
)

var _ = gc.Suite(&distroSuite{})

type distroSuite struct {
	testing.IsolationSuite
}

func (distroSuite) TestKnownOkay(c *gc.C) {
	known := map[os.DistroID]string{
		os.DistroUbuntu: "Ubuntu",
		os.DistroDebian: "Debian",
		os.DistroRedHat: "RedHat",
		os.DistroFedora: "Fedora",
		os.DistroCentOS: "CentOS",
		os.DistroArch:   "Arch",
		os.DistroSUSE:   "SUSE",
	}
	c.Check(os.Distros, gc.HasLen, len(known))
	for id, name := range known {
		c.Logf("checking %q", name)

		c.Check(os.Distros[id], jc.DeepEquals, os.Distro{
			ID:   id,
			Name: name,
		})
	}
}

func (distroSuite) TestKnownValid(c *gc.C) {
	for _, distro := range os.Distros {
		c.Logf("checking %q", distro)
		err := distro.Validate()

		c.Check(err, jc.ErrorIsNil)
	}
}

func (distroSuite) TestValidateOkay(c *gc.C) {
	distro := os.Distro{
		ID:   os.DistroUbuntu,
		Name: "Ubuntu",
	}
	err := distro.Validate()

	c.Check(err, jc.ErrorIsNil)
}

func (distroSuite) TestValidateMissingID(c *gc.C) {
	distro := os.Distro{
		Name: "Ubuntu",
	}
	err := distro.Validate()

	c.Check(err, gc.ErrorMatches, `.*ID must be set`)
}

func (distroSuite) TestValidateMissingName(c *gc.C) {
	distro := os.Distro{
		ID: os.DistroUbuntu,
	}
	err := distro.Validate()

	c.Check(err, gc.ErrorMatches, `.*Name must be set`)
}

func (distroSuite) TestValidateEmpty(c *gc.C) {
	var distro os.Distro
	err := distro.Validate()

	c.Check(err, gc.NotNil)
}

func (distroSuite) TestStringCapitalized(c *gc.C) {
	var distro os.Distro
	distro.Name = "Spam"
	str := distro.String()

	c.Check(str, gc.Equals, "spam")
}

func (distroSuite) TestStringUpper(c *gc.C) {
	var distro os.Distro
	distro.Name = "SPAM"
	str := distro.String()

	c.Check(str, gc.Equals, "spam")
}

func (distroSuite) TestStringLower(c *gc.C) {
	var distro os.Distro
	distro.Name = "spam"
	str := distro.String()

	c.Check(str, gc.Equals, "spam")
}

var _ = gc.Suite(&distroIDSuite{})

type distroIDSuite struct {
	testing.IsolationSuite
}

func (distroIDSuite) TestStringKnown(c *gc.C) {
	id := os.DistroUbuntu
	str := id.String()

	c.Check(str, gc.Equals, "ubuntu")
}

func (distroIDSuite) TestStringUnknown(c *gc.C) {
	id := os.DistroID(99)
	str := id.String()

	c.Check(str, gc.Equals, "unknown")
}

func (distroIDSuite) TestStringZeroValue(c *gc.C) {
	var id os.DistroID
	str := id.String()

	c.Check(str, gc.Equals, "unknown")
}

func (distroIDSuite) TestInfoKnown(c *gc.C) {
	id := os.DistroUbuntu
	info, ok := id.Info()

	c.Check(ok, jc.IsTrue)
	c.Check(info, gc.Equals, os.Distros[os.DistroUbuntu])
}

func (distroIDSuite) TestInfoUnknown(c *gc.C) {
	id := os.DistroID(99)
	info, ok := id.Info()

	c.Check(ok, jc.IsFalse)
	c.Check(info, gc.Equals, os.Distro{})
}

func (distroIDSuite) TestInfoZeroValue(c *gc.C) {
	var id os.DistroID
	info, ok := id.Info()

	c.Check(ok, jc.IsFalse)
	c.Check(info, gc.Equals, os.Distro{})
}
