// Copyright 2015 Eric Snow
// Licensed under the New BSD License, see LICENSE file for details.

package os

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

// Distro identifies a linux distribution.
type DistroID uint
