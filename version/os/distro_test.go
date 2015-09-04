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
