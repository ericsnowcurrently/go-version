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
