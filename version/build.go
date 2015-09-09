// Copyright 2015 Eric Snow
// Licensed under the New BSD License, see LICENSE file for details.

package version

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/juju/errors"
)

var (
	buildPat = fmt.Sprintf(`(?:(%s)(\.[1-9]\d*)?)`, relPat)
	buildRE  = regexp.MustCompile(`^` + buildPat + `(.*)$`)
)

// Build represents a build of a software release.
type Build struct {
	Release
	// Index uniquely identifies the build of the release.
	Index uint
}

// ParseBuild converts the provided release build version number to a
// Build. The unused portion of the string is also returned.
func ParseBuild(builds string) (Build, string, error) {
	var build Build
	var err error

	parts := buildRE.FindStringSubmatch(builds)
	if len(parts) == 0 {
		return build, "", errors.NotValidf("build version %q", builds)
	}
	remainder := parts[relRE.NumSubexp()]

	rel, _, err := ParseRelease(builds)
	if err != nil {
		return build, "", errors.Trace(err)
	}
	build.Release = rel

	indexStart := relRE.NumSubexp() + 1
	if parts[indexStart] != "" {
		_, err = fmt.Sscanf(parts[indexStart], ".%d", &build.Index)
		err = errors.Trace(err)
	}

	if err != nil {
		return build, "", errors.Wrap(err, errors.NotValidf("build version %q", builds))
	}

	return build, remainder, nil
}

// String returns the string representation of the build.
func (build Build) String() string {
	if build.Index == 0 {
		return build.Release.String()
	}
	return fmt.Sprintf("%s.%d", build.Release, build.Index)
}

// IsZero determines whether or not the Build is the "zero" value.
func (build Build) IsZero() bool {
	return build == Build{}
}

// Validate ensures that the Build is valid. If not then it returns
// errors.NotValid.
func (build Build) Validate() error {
	if err := build.Release.Validate(); err != nil {
		return errors.Trace(err)
	}
	if build.Index == 0 {
		return errors.NotValidf("index (must be non-zero)")
	}
	return nil
}

// Compare returns -1, 0 or 1 depending on whether
// v is less than, equal to or greater than w.
func (build Build) Compare(other Build) int {
	compared := build.Release.Compare(other.Release)
	if compared != 0 {
		return compared
	}
	switch {
	case build.Index < other.Index:
		return -1
	case build.Index > other.Index:
		return 1
	}
	return 0
}

// Prev calculates the previous Build to this one.
// If there is no previous then false is returned.
func (build Build) Prev() (Build, bool) {
	if build.Index <= 1 {
		return build, false
	}
	prev := build
	prev.Index -= 1
	return prev, true
}

// Next calculates the next Build to this one.
// It returns false once the max is hit (if greater than 0).
func (build Build) Next(max int) (Build, bool) {
	if max >= 0 && int(build.Index) >= max {
		return build, false
	}
	next := build
	next.Index += 1
	return next, true
}

// MarshalJSON implements json.Marshaler.
func (build Build) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(build.String())
	if err != nil {
		return data, errors.Trace(err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (build *Build) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return errors.Trace(err)
	}
	// TODO(ericsnow) assert no remainder?
	parsed, _, err := ParseBuild(str)
	if err != nil {
		return errors.Trace(err)
	}
	*build = parsed
	return nil
}

// GetYAML implements goyaml.Getter
func (build Build) GetYAML() (tag string, value interface{}) {
	return "", build.String()
}

// SetYAML implements goyaml.Setter
func (build *Build) SetYAML(tag string, value interface{}) bool {
	str := fmt.Sprintf("%v", value)
	if str == "" {
		return false
	}
	// TODO(ericsnow) assert no remainder?
	parsed, _, err := ParseBuild(str)
	if err != nil {
		return false
	}
	*build = parsed
	return true
}
