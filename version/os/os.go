// Copyright 2015 Eric Snow
// Licensed under the New BSD License, see LICENSE file for details.

package os

// These are the names of the operating systems recognized by Go.
const (
	Unknown = ""

	Darwin    = "darwin"
	Dragonfly = "dragonfly"
	FreeBSD   = "freebsd"
	Linux     = "linux"
	Nacl      = "nacl"
	NetBSD    = "netbsd"
	OpenBSD   = "openbsd"
	Solaris   = "solaris"

	Windows = "windows"
)

// unix is the list of unix-like operating systems recognized by Go.
// See http://golang.org/src/path/filepath/path_unix.go.
var unix = map[string]string{
	Darwin:    Darwin,
	Dragonfly: Dragonfly,
	FreeBSD:   FreeBSD,
	Linux:     Linux,
	Nacl:      Nacl,
	NetBSD:    NetBSD,
	OpenBSD:   OpenBSD,
	Solaris:   Solaris,
}

// OSIsUnix determines whether or not the given OS name is one of the
// unix-like operating systems recognized by Go.
func IsUnix(os string) bool {
	_, ok := unix[os]
	return ok
}
