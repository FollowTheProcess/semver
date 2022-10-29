// Package semver is a library for parsing and validating semantic versions in Go.
//
// It is fully compliant with the [semver 2.0.0 spec]
//
// [semver 2.0.0 spec]: https://semver.org
package semver

import (
	"fmt"
	"regexp"
	"strconv"
)

// Here primarily to avoid typos.
const (
	major = "major"
	minor = "minor"
	patch = "patch"
	pre   = "pre"
	build = "build"
)

// See https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
// only difference we've made is this one supports an optional 'v' on the front.
var semVerRegex = regexp.MustCompile(fmt.Sprintf(`^v?(?P<%s>0|[1-9]\d*)\.(?P<%s>0|[1-9]\d*)\.(?P<%s>0|[1-9]\d*)(?:-(?P<%s>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<%s>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`, major, minor, patch, pre, build))

// Version encodes a semantic version.
type Version struct {
	Prerelease string
	Build      string
	Major      uint64
	Minor      uint64
	Patch      uint64
}

// String implements the Stringer interface and allows a Version to print itself.
//
//	v := Version{Major: 1, Minor: 2, Patch: 3}
//	fmt.Println(v) // "1.2.3"
func (v Version) String() string {
	base := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Prerelease != "" {
		base += "-" + v.Prerelease
	}
	if v.Build != "" {
		base += "+" + v.Build
	}
	return base
}

// Tag creates a string representation of the version suitable for git tags
// it is identical to the String() method except prepends a 'v' to the result.
//
//	v := Version{Major: 1, Minor: 2, Patch: 3}
//	fmt.Println(v) // "v1.2.3"
func (v Version) Tag() string {
	return "v" + v.String()
}

// New creates and returns a new Version.
func New(major, minor, patch uint64, build, pre string) Version {
	return Version{
		Prerelease: pre,
		Build:      build,
		Major:      major,
		Minor:      minor,
		Patch:      patch,
	}
}

// Parse creates and returns a Version from a semver string.
func Parse(text string) (Version, error) {
	if !semVerRegex.MatchString(text) {
		return Version{}, fmt.Errorf("%q is not a valid semantic version", text)
	}

	groups := make(map[string]string, 5) // 5 elements (fields of Version)
	parts := semVerRegex.FindStringSubmatch(text)
	names := semVerRegex.SubexpNames()

	for i, part := range parts {
		// First element in parts is the substring e.g. "v1.2.4"
		// first element in names is empty string -> skip both
		if i == 0 {
			continue
		}
		groups[names[i]] = part
	}

	// Errors below are ignored because they wouldn't pass the regex check if they
	// weren't parseable numeric digits
	majorInt, _ := strconv.ParseUint(groups[major], 10, 64)
	minorInt, _ := strconv.ParseUint(groups[minor], 10, 64)
	patchInt, _ := strconv.ParseUint(groups[patch], 10, 64)

	v := Version{
		Prerelease: groups[pre],
		Build:      groups[build],
		Major:      majorInt,
		Minor:      minorInt,
		Patch:      patchInt,
	}

	return v, nil
}
