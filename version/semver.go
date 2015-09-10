// Copyright 2015 Eric Snow
// Licensed under the New BSD License, see LICENSE file for details.

package version

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/juju/errors"
)

var (
	identPat = strings.Replace(strings.Replace(`
(
    [0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*
)
`, "\n", "", -1), " ", "", -1)

	semVerPat = fmt.Sprintf(`(%s(?:-%s)?(:+%s)?)`, numPat, identPat, identPat)
	semverRE  = regexp.MustCompile(`^` + semVerPat + `$`)
)

/*
var semverRE = regexp.MustCompile(strings.Replace(strings.Replace(`^
`^
(
    (0 | [1-9]\d*)
    \.(0 | [1-9]\d*)
    \.(0 | [1-9]\d*)
)
(?:-(
    [0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*
))?
(?:\+(
    [0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*
))?
$`, "\n", "", -1), " ", "", -1))
*/

// See http://semver.org/.
type SemVer struct {
	Number     Number
	PreRelease []string
	Build      []string
}

func ParseSemVer(vers string) (SemVer, error) {
	var ver SemVer
	parts := semverRE.FindStringSubmatch(vers)
	if len(parts) == 0 {
		return ver, errors.NotValidf("SemVer (2.0) %q", vers)
	}
	num, _, err := ParseNumber(parts[1])
	if err != nil {
		return ver, errors.Trace(err)
	}
	ver.Number = num
	ver.PreRelease = strings.Split(parts[5], ".")
	ver.Build = strings.Split(parts[6], ".")
	return ver, nil
}

func (ver SemVer) Major() uint {
	return ver.Number.Major
}

func (ver SemVer) Minor() uint {
	return ver.Number.Minor
}

func (ver SemVer) Patch() uint {
	return ver.Number.Micro
}

func (ver SemVer) String() string {
	str, _ := ver.toString()
	return str
}

func (ver SemVer) Validate() error {
	_, err := ver.toString()
	return errors.Trace(err)
}

func (ver SemVer) toString() (string, error) {
	vers := ver.String()
	if len(ver.PreRelease) > 0 {
		vers += "-" + strings.Join(ver.PreRelease, ".")
		if _, err := ParseSemVer(vers); err != nil {
			return "", errors.NotValidf("pre-release %v", ver.PreRelease)
		}
	}
	if len(ver.Build) > 0 {
		vers += "+" + strings.Join(ver.Build, ".")
		if _, err := ParseSemVer(vers); err != nil {
			return "", errors.NotValidf("build %v", ver.Build)
		}
	}
	return vers, nil
}
