// Copyright 2015 Eric Snow
// Licensed under the New BSD License, see LICENSE file for details.

package os

import (
	"github.com/juju/errors"
)

func init() {
	for id, distro := range distros {
		distro.ID = id
		distros[id] = distro
	}
}

// These are recognized linux distributions.
const (
	DistroUnknown DistroID = iota
	DistroUbuntu
	DistroDebian
	DistroRedHat
	DistroFedora
	DistroCentOS
	DistroArch
	DistroSUSE
)

var distros = map[DistroID]Distro{
	DistroUbuntu: Distro{
		Name: "Ubuntu",
	},
	DistroDebian: Distro{
		Name: "Debian",
	},
	DistroRedHat: Distro{
		Name: "RedHat",
	},
	DistroFedora: Distro{
		Name: "Fedora",
	},
	DistroCentOS: Distro{
		Name: "CentOS",
	},
	DistroArch: Distro{
		Name: "Arch",
	},
	DistroSUSE: Distro{
		Name: "SUSE",
	},
}

// Distro contains information about a linux distribution.
type Distro struct {
	// ID is the unique identifier for the distro.
	ID DistroID
	// Name is the name of the distro.
	Name string
}

// Validate returns an error if the distro is not valid.
func (distro Distro) Validate() error {
	// TODO(ericsnow) Use errors.NotValidf?
	if distro.ID == DistroUnknown {
		return errors.Errorf("distro.ID must be set")
	}
	if distro.Name == "" {
		return errors.Errorf("distro.Name must be set")
	}
	return nil
}

// Distro identifies a linux distribution.
type DistroID uint
