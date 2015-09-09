// Copyright 2015 Eric Snow
// Licensed under the New BSD License, see LICENSE file for details.

package version

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/juju/errors"
)

// TODO(ericsnow) Is this a good enough "unknown" value?
const unknown = "unknown"

var (
	binFieldPat  = `[0-9A-Za-z]`
	binSuffixPat = fmt.Sprintf(`(?:-(%[1]s+)-(%[1]s+))`, binFieldPat)
	binPat       = fmt.Sprintf(`(%s)(%s)`, buildPat, binSuffixPat)
	binRE        = regexp.MustCompile(`^` + binPat + `(.*)$`)
)

// Binary represents the version of a built binary for some software.
type Binary struct {
	Build
	// Series is the targeted OS series, as identified by the operating
	// system (e.g. trusty).
	Series string
	// Arch is the targeted host architecture (e.g. amd64).
	Arch string
}

// ParseBinary converts the provided binary version string to a Binary.
// The unused portion of the string is also returned.
func ParseBinary(bins string) (Binary, string, error) {
	var bin Binary

	parts := binRE.FindStringSubmatch(bins)
	if len(parts) == 0 {
		return bin, "", errors.NotValidf("binary version %q", bins)
	}
	remainder := parts[binRE.NumSubexp()]

	suffixStart := buildRE.NumSubexp() + 1
	series, arch := parts[suffixStart+1], parts[suffixStart+2]
	// TODO(ericsnow) Leave them as "unknown"?
	if series == unknown {
		series = ""
	}
	if arch == unknown {
		arch = ""
	}

	build, _, err := ParseBuild(bins)
	if err != nil {
		return bin, "", errors.Trace(err)
	}
	bin.Build = build
	bin.Series = series
	bin.Arch = arch
	return bin, remainder, nil
}

// String returns the string representation of this Binary.
func (bin Binary) String() string {
	// TODO(ericsnow) Omit the series/arch if not set? Fail?
	series, arch := bin.Series, bin.Arch
	if series == "" {
		series = unknown
	}
	if arch == "" {
		arch = unknown
	}
	return fmt.Sprintf("%s-%s-%s", bin.Build, series, arch)
}

// IsZero determines whether or not the Binary is the "zero" value.
func (bin Binary) IsZero() bool {
	return bin == Binary{}
}

// Validate() ensures the Binary is valid. If not then it returns
// errors.NotValid.
func (bin Binary) Validate() error {
	if err := bin.Build.Validate(); err != nil {
		return errors.Trace(err)
	}

	if bin.Series == "" || bin.Series == unknown {
		return errors.NotValidf("binary series missing")
	}
	fieldPat := fmt.Sprintf(`^%s+$`, binFieldPat)
	if matched, _ := regexp.MatchString(fieldPat, bin.Series); !matched {
		return errors.NotValidf("unrecognized binary series %q", bin.Series)
	}

	if bin.Arch == "" || bin.Arch == unknown {
		return errors.NotValidf("binary arch missing")
	}
	if matched, _ := regexp.MatchString(fieldPat, bin.Arch); !matched {
		return errors.NotValidf("unrecognized binary arch %q", bin.Arch)
	}

	return nil
}

// Compare returns -1, 0 or 1 depending on whether
// v is less than, equal to or greater than w.
func (bin Binary) Compare(other Binary) int {
	return bin.Build.Compare(other.Build)
}

// Prev calculates the previous Binary to this one.
// If there is no previous then false is returned.
func (bin Binary) Prev() (Binary, bool) {
	build, ok := bin.Build.Prev()
	if !ok {
		return bin, ok
	}
	prev := Binary{
		Build:  build,
		Series: bin.Series,
		Arch:   bin.Arch,
	}
	return prev, true
}

// Next calculates the next Binary to this one.
// It returns false once the max is hit (if greater than 0).
func (bin Binary) Next(max int) (Binary, bool) {
	build, ok := bin.Build.Next(max)
	if !ok {
		return bin, ok
	}
	next := Binary{
		Build:  build,
		Series: bin.Series,
		Arch:   bin.Arch,
	}
	return next, true
}

// MarshalJSON implements json.Marshaler.
func (bin Binary) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(bin.String())
	if err != nil {
		return data, errors.Trace(err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (bin *Binary) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return errors.Trace(err)
	}
	// TODO(ericsnow) assert no remainder?
	parsed, _, err := ParseBinary(str)
	if err != nil {
		return errors.Trace(err)
	}
	*bin = parsed
	return nil
}

// GetYAML implements goyaml.Getter
func (bin Binary) GetYAML() (tag string, value interface{}) {
	return "", bin.String()
}

// SetYAML implements goyaml.Setter
func (bin *Binary) SetYAML(tag string, value interface{}) bool {
	str := fmt.Sprintf("%v", value)
	if str == "" {
		return false
	}
	// TODO(ericsnow) assert no remainder?
	parsed, _, err := ParseBinary(str)
	if err != nil {
		return false
	}
	*bin = parsed
	return true
}
