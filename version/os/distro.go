// Copyright 2015 Eric Snow
// Licensed under the New BSD License, see LICENSE file for details.

package os

import (
	"strings"

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

// String returns the string representation of the distro. It is
// rendered as the lower-cased distro name.
func (distro Distro) String() string {
	return strings.ToLower(distro.Name)
}

// Distro identifies a linux distribution.
type DistroID uint

// String returns a string representation of the distro.
func (id DistroID) String() string {
	if info, ok := distros[id]; ok {
		return info.String()
	}
	return "unknown"
}

// Info returns information about the distro, if recognized. If not
// recognized then false is returned.
func (id DistroID) Info() (Distro, bool) {
	if distro, ok := distros[id]; ok {
		copied := distro
		return copied, true
	}
	return Distro{}, false
}
