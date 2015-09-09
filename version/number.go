// Copyright 2015 Eric Snow
// Licensed under the New BSD License, see LICENSE file for details.

package version

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/juju/errors"
)

var (
	numPat = strings.Replace(strings.Replace(`
(
    (0 | [1-9]\d*)
    (?:\.(0 | [1-9]\d*))?
    (?:\.(0 | [1-9]\d*))?
)
`, "\n", "", -1), " ", "", -1)
	numRE = regexp.MustCompile(`^` + numPat + `(.*)$`)
)

// Number represents a simple 3-part software/API version.
type Number struct {
	// Major is the version number that changes with a break in
	// compatibility. For APIs it represents the targeted interface version.
	Major uint
	// Minor is the version number that changes when new features are
	// added. For APIs it represents the implementation # of the
	// targeted interface.
	Minor uint
	// Micro is the version number that changes with bug fixes. It's
	// also known as the "patch" level.
	Micro uint
}

// ParseNumber converts a version number string into a Number. This
// conversion works for complete as well as incomplete versions:
//   "2.3.1" -> Number(2, 3, 1)
//   "2.3"   -> Number(2, 3, 0)
//   "2"     -> Number(2, 0, 0)
//
// Unsupported version strings result in errors.NotValid.
func ParseNumber(nums string) (Number, string, error) {
	var num Number
	parts := numRE.FindStringSubmatch(nums)
	if len(parts) == 0 {
		return num, "", errors.NotValidf("version string %q", nums)
	}
	remainder := parts[numRE.NumSubexp()]
	if parts[4] == "" && strings.HasPrefix(remainder, ".") {
		return num, "", errors.NotValidf("version string %q", nums)
	}

	nums = parts[1]
	var err error
	switch strings.Count(nums, ".") {
	case 2:
		_, err = fmt.Sscanf(nums, "%d.%d.%d", &num.Major, &num.Minor, &num.Micro)
		err = errors.Trace(err)
	case 1:
		_, err = fmt.Sscanf(nums, "%d.%d", &num.Major, &num.Minor)
		err = errors.Trace(err)
	case 0:
		_, err = fmt.Sscanf(nums, "%d", &num.Major)
		err = errors.Trace(err)
	}
	if err != nil {
		return num, "", errors.Wrap(err, errors.NotValidf("version string %q", nums))
	}
	return num, remainder, nil
}

// TODO(ericsnow) Add ParseNumberExact?

// String converts the Number to its string representation.
//
// Note that this will not round-trip with ParseNumber if the version
// string was originally incomplete.
func (num Number) String() string {
	return fmt.Sprintf("%d.%d.%d", num.Major, num.Minor, num.Micro)
}

// Feature converts the Number to the string representation of its
// "feature" version string (first 2 numbers only).
func (num Number) Feature() string {
	return fmt.Sprintf("%d.%d", num.Major, num.Minor)
}

// IsZero determines whether or not the Number is the "zero" value.
func (num Number) IsZero() bool {
	return num == Number{}
}

// Validate ensures that the Number is valid. If not then it returns
// errors.NotValid.
func (num Number) Validate() error {
	if num.IsZero() {
		return errors.NotValidf("zero-value Number")
	}
	return nil
}

// Compare returns -1, 0 or 1 depending on whether
// v is less than, equal to or greater than w.
func (num Number) Compare(other Number) int {
	const (
		less    = -1
		equal   = 0
		greater = 1
	)

	switch {
	case num.Major < other.Major:
		return less
	case num.Major > other.Major:
		return greater
	default:
		switch {
		case num.Minor < other.Minor:
			return less
		case num.Minor > other.Minor:
			return greater
		default:
			switch {
			case num.Micro < other.Micro:
				return less
			case num.Micro > other.Micro:
				return greater
			}
		}
	}
	return equal
}

// Prev calculates the previous Number to this one. When it must wrap
// around it wraps to the bound set by the provided max. If there is no
// previous then false is returned.
func (num Number) Prev(max Number) (Number, bool) {
	major, minor, micro := num.Major, num.Minor, num.Micro
	if micro == 0 {
		if minor == 0 {
			if major == 0 {
				return num, false
			}
			major -= 1
			minor = max.Minor
		} else {
			minor -= 1
		}
		micro = max.Micro
	} else {
		micro -= 1
	}
	return Number{major, minor, micro}, true
}

// Next calculates the next Number to this one. When it must wrap
// around it wraps to the bound set by the provided max. If there is no
// next then false is returned.
func (num Number) Next(max Number) (Number, bool) {
	var major, minor, micro uint
	switch {
	case num.Micro < max.Micro:
		major = num.Major
		minor = num.Minor
		micro = num.Micro + 1
	case num.Minor < max.Minor:
		major = num.Major
		minor = num.Minor + 1
	case num.Major < max.Major:
		major = num.Major + 1
	default:
		return Number{}, false
	}
	return Number{major, minor, micro}, true
}

// MarshalJSON implements json.Marshaler.
func (num Number) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(num.String())
	if err != nil {
		return data, errors.Trace(err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (num *Number) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return errors.Trace(err)
	}
	// TODO(ericsnow) assert no remainder?
	parsed, _, err := ParseNumber(str)
	if err != nil {
		return errors.Trace(err)
	}
	*num = parsed
	return nil
}

// GetYAML implements goyaml.Getter
func (num Number) GetYAML() (tag string, value interface{}) {
	return "", num.String()
}

// SetYAML implements goyaml.Setter
func (num *Number) SetYAML(tag string, value interface{}) bool {
	str := fmt.Sprintf("%v", value)
	if str == "" {
		return false
	}
	// TODO(ericsnow) assert no remainder?
	parsed, _, err := ParseNumber(str)
	if err != nil {
		return false
	}
	*num = parsed
	return true
}
