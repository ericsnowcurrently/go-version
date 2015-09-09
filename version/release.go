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
	relPat = fmt.Sprintf(`(?:%s%s?)`, numPat, releaseLevelPat)
	relRE  = regexp.MustCompile(`^` + relPat + `(.*)$`)
)

// These are the recognized release levels.
const (
	ReleaseDevelopment ReleaseLevel = iota
	ReleaseAlpha
	ReleaseBeta
	ReleaseCandidate
	ReleaseFinal
)

// Release represents a 3-part version plus release info.
type Release struct {
	Number
	// Level is the release level (e.g. "alpha")
	Level ReleaseLevel
	// Serial is the increment within the release level.
	//
	// Note that Serial does not apply to "development" and "final"
	// releases.
	Serial uint
}

// ParseRelease converts a release version string into a Release.
// Supported release names are "dev" (development), "alpha", "beta",
// "candidate", and "final":
//   "2.3.1-dev" -> Release(2, 3, 1, ReleaseDevelopment, 0)
//   "2.3.1-alpha1" -> Release(2, 3, 1, ReleaseAlpha, 1)
//   "2.3.1-beta1" -> Release(2, 3, 1, ReleaseBeta, 1)
//   "2.3.1-candidate1" -> Release(2, 3, 1, ReleaseCandidate, 1)
//   "2.3.1-final" -> Release(2, 3, 1, ReleaseFinal, 0)
//
// Abbreviated release names are supported for alpha, beta, and candidate:
//   "2.3.1a1" -> Release(2, 3, 1, ReleaseAlpha, 1)
//   "2.3.1b1" -> Release(2, 3, 1, ReleaseBeta, 1)
//   "2.3.1rc1" -> Release(2, 3, 1, ReleaseCandidate, 1)
//
// If the release part of the string is missing then it defaults to "final":
//   "2.3.1" -> Release(2, 3, 1, ReleaseFinal, 0)
//
// The conversion fails for incomplete-number release strings:
//   "2.3a1"     -> errors.NotValid
//   "2.3-beta2" -> errors.NotValid
//
// Unsupported version strings result in errors.NotValid.
func ParseRelease(rels string) (Release, string, error) {
	var rel Release

	parts := relRE.FindStringSubmatch(rels)
	if len(parts) == 0 {
		return rel, "", errors.NotValidf("release version %q", rels)
	}
	if parts[3] == "" { // minor
		return rel, "", errors.NotValidf("release version %q", rels)
	}

	num, _, err := ParseNumber(rels)
	if err != nil {
		return rel, "", errors.Wrap(err, errors.NotValidf("release version %q", rels))
	}
	rel.Number = num

	remainder := parts[relRE.NumSubexp()]
	if remainder != "" {
		remainder2 := remainder
		if remainder2[0] == '-' {
			remainder2 = remainder2[1:]
		}
		for _, name := range releaseLevelNames {
			if strings.HasPrefix(remainder2, name) {
				return rel, "", errors.NotValidf("release version %q", rels)
			}
		}
	}

	relStart := numRE.NumSubexp()
	switch {
	case parts[relStart] == "":
		rel.Level = ReleaseFinal
	case parts[relStart+1] != "":
		rel.Level = ReleaseLevelFromName(parts[relStart+1])
		oops := -1
		fmt.Sscanf(remainder, "%d", &oops)
		if oops >= 0 {
			return rel, "", errors.NotValidf("release version %q", rels)
		}
	case parts[relStart+2] != "":
		rel.Level = ReleaseLevelFromName(parts[relStart+3])
		fmt.Sscanf(parts[relStart+4], "%d", &rel.Serial)
	case parts[relStart+5] != "":
		rel.Level = ReleaseLevelFromAbbrev(parts[relStart+6])
		fmt.Sscanf(parts[relStart+7], "%d", &rel.Serial)
	}

	return rel, remainder, nil
}

// String returns the string representation of the release version.
func (rel Release) String() string {
	switch rel.Level {
	case ReleaseDevelopment, ReleaseFinal:
		return fmt.Sprintf("%s-%s", rel.Number, rel.Level)
	default:
		return fmt.Sprintf("%s-%s%d", rel.Number, rel.Level, rel.Serial)
	}
}

// Abbrev returns the abbreviated string representation of the release
// version, of the full string if the level does not have an
// abbreviation.
func (rel Release) Abbrev() string {
	switch rel.Level {
	case ReleaseDevelopment, ReleaseFinal:
		return rel.Number.String()
	default:
		return fmt.Sprintf("%s%s%d", rel.Number, rel.Level.Abbrev(), rel.Serial)
	}
}

// IsZero determines whether or not the release version is the zero value.
func (rel Release) IsZero() bool {
	return rel == Release{}
}

// Compare returns -1, 0 or 1 depending on whether
// v is less than, equal to or greater than w.
func (rel Release) Compare(other Release) int {
	compared := rel.Number.Compare(other.Number)
	if compared != 0 {
		return compared
	}
	switch {
	case rel.Level < other.Level:
		return -1
	case rel.Level > other.Level:
		return 1
	case rel.Serial < other.Serial:
		return -1
	case rel.Serial > other.Serial:
		return 1
	}
	return 0
}

// Prev calculates the previous release level. When the level must wrap
// around it wraps to the bounds set by the provided max alpha, beta,
// and release. If there is no previous then false is returned.
func (rel Release) Prev(aMax, bMax, rcMax uint) (Release, bool) {
	prev := Release{Number: rel.Number}
	if rel.Serial <= 1 {
		switch rel.Level {
		case ReleaseFinal:
			prev.Level = ReleaseCandidate
			prev.Serial = rcMax
		case ReleaseCandidate:
			prev.Level = ReleaseBeta
			prev.Serial = bMax
		case ReleaseBeta:
			prev.Level = ReleaseAlpha
			prev.Serial = aMax
		case ReleaseAlpha:
			prev.Level = ReleaseDevelopment
			prev.Serial = 0
		case ReleaseDevelopment:
			return rel, false
		}
	} else {
		prev.Level = rel.Level
		prev.Serial = rel.Serial - 1
	}
	return prev, true
}

// Next calculates the previous release level. When the level must wrap
// around it wraps to the bounds set by the provided max alpha, beta,
// and release. If there is no next then false is returned.
func (rel Release) Next(aMax, bMax, rcMax uint) (Release, bool) {
	next := Release{
		Number: rel.Number,
		Level:  rel.Level,
		Serial: 1,
	}
	switch rel.Level {
	case ReleaseDevelopment:
		next.Level = ReleaseAlpha
	case ReleaseAlpha:
		if rel.Serial < aMax {
			next.Serial = rel.Serial + 1
		} else {
			next.Level = ReleaseBeta
		}
	case ReleaseBeta:
		if rel.Serial < bMax {
			next.Serial = rel.Serial + 1
		} else {
			next.Level = ReleaseCandidate
		}
	case ReleaseCandidate:
		if rel.Serial < rcMax {
			next.Serial = rel.Serial + 1
		} else {
			next.Level = ReleaseFinal
			next.Serial = 0
		}
	case ReleaseFinal:
		return rel, false
	}
	return next, true
}

// MarshalJSON implements json.Marshaler.
func (rel Release) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(rel.String())
	if err != nil {
		return data, errors.Trace(err)
	}
	return data, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (rel *Release) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return errors.Trace(err)
	}
	parsed, _, err := ParseRelease(str)
	if err != nil {
		return errors.Trace(err)
	}
	*rel = parsed
	return nil
}

// GetYAML implements goyaml.Getter
func (rel Release) GetYAML() (tag string, value interface{}) {
	return "", rel.String()
}

// SetYAML implements goyaml.Setter
func (rel *Release) SetYAML(tag string, value interface{}) bool {
	str := fmt.Sprintf("%v", value)
	if str == "" {
		return false
	}
	parsed, _, err := ParseRelease(str)
	if err != nil {
		return false
	}
	*rel = parsed
	return true
}

var (
	releaseLevelPat = strings.Replace(strings.Replace(`
(
    -(dev|final)
    |
    -((alpha|beta|candidate)([1-9]\d*))
    |
    ((a|b|rc)([1-9]\d*))
)
`, "\n", "", -1), " ", "", -1)
	releaseLevelRE = regexp.MustCompile(`^` + releaseLevelPat + `$`)
)

var (
	releaseLevelNames = map[ReleaseLevel]string{
		ReleaseDevelopment: "dev",
		ReleaseAlpha:       "alpha",
		ReleaseBeta:        "beta",
		ReleaseCandidate:   "candidate",
		ReleaseFinal:       "final",
	}

	// These must remain alphabetically ordered.
	releaseLevelAbbrevs = map[ReleaseLevel]string{
		ReleaseAlpha:     "a",
		ReleaseBeta:      "b",
		ReleaseCandidate: "rc",
	}

	releaseLevelSingletons = map[ReleaseLevel]bool{
		ReleaseDevelopment: true,
		ReleaseFinal:       true,
	}
)

// ReleaseLevel represents one of the recognized release levels.
type ReleaseLevel int

// ReleaseLevelFromName converts the provided string into the
// corresponding release level. It defaults to "final".
func ReleaseLevelFromName(name string) ReleaseLevel {
	for level, known := range releaseLevelNames {
		if name == known {
			return level
		}
	}
	return ReleaseFinal
}

// ReleaseLevelFromAbbrev converts the provided string into the
// corresponding release level. It defaults to "final".
func ReleaseLevelFromAbbrev(abbrev string) ReleaseLevel {
	for level, known := range releaseLevelAbbrevs {
		if abbrev == known {
			return level
		}
	}
	return ReleaseFinal
}

// String returns the string representation of the release level.
func (rl ReleaseLevel) String() string {
	return releaseLevelNames[rl]
}

// Abbrev returns the abbreviated string representation of the
// release level.
func (rl ReleaseLevel) Abbrev() string {
	return releaseLevelAbbrevs[rl]
}

// Index returns an integer that gives the absolute order of the
// release level. This can be used when the release level must
// be represented as an integer.
func (rl ReleaseLevel) Index() int {
	switch rl {
	case ReleaseAlpha:
		return 0
	case ReleaseBeta:
		return 1
	case ReleaseCandidate:
		return 2
	case ReleaseFinal:
		return 3
	default:
		return -1
	}
}
