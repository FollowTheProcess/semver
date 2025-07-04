// Package semver is a library for parsing, validating and bumping semantic versions in Go.
//
// It is fully compliant with the [semver 2.0.0 spec]
//
// [semver 2.0.0 spec]: https://semver.org
package semver // import "go.followtheprocess.codes/semver"

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

	groups = 5 // major, minor, patch, pre, build
)

// See https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
// only difference we've made is this one supports an optional 'v' on the front.
var semVerRegex = regexp.MustCompile(
	`^v?(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<pre>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<build>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`,
)

// Version encodes a semantic version.
type Version struct {
	Prerelease string // Optional pre-release e.g. "rc1"
	Build      string // Optional build metadata e.g. "build.123"
	Major      uint   // Major version
	Minor      uint   // Minor version
	Patch      uint   // Patch version
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
// Numeric parts are unsigned integers so that e.g -1 becomes a compile time error
//
//	v := New(1, 7, 6, "", "")
//	Version{Major: 1, Minor: 7, Patch: 6, Prerelease: "", Build: ""}
func New(major, minor, patch uint, build, pre string) Version {
	return Version{
		Prerelease: pre,
		Build:      build,
		Major:      major,
		Minor:      minor,
		Patch:      patch,
	}
}

// Parse creates and returns a Version from a semver string.
//
// If the string is not a valid semantic version, an error will be returned
//
//	v, _ := Parse("v1.8.9")
//	Version{Major: 1, Minor: 8, Patch: 9, Prerelease: "", Build: ""}
func Parse(text string) (Version, error) {
	if !IsValid(text) {
		return Version{}, fmt.Errorf("%q is not a valid semantic version", text)
	}

	groups := make(map[string]string, groups) // 5 elements (fields of Version)
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
	majorInt, _ := strconv.ParseUint(groups[major], 10, 64) //nolint: errcheck // Must be valid digits to pass regex
	minorInt, _ := strconv.ParseUint(groups[minor], 10, 64) //nolint: errcheck // Must be valid digits to pass regex
	patchInt, _ := strconv.ParseUint(groups[patch], 10, 64) //nolint: errcheck // Must be valid digits to pass regex

	v := Version{
		Prerelease: groups[pre],
		Build:      groups[build],
		Major:      uint(majorInt),
		Minor:      uint(minorInt),
		Patch:      uint(patchInt),
	}

	return v, nil
}

// BumpMajor returns a new Version with it's major version bumped.
func BumpMajor(current Version) Version {
	// Everything else set to zero value
	return Version{
		Major: current.Major + 1,
	}
}

// BumpMinor returns a new Version with it's minor version bumped.
func BumpMinor(current Version) Version {
	// Keep major, bump minor, everything else -> zero value
	return Version{
		Major: current.Major,
		Minor: current.Minor + 1,
	}
}

// BumpPatch returns a new Version with it's patch version bumped.
func BumpPatch(current Version) Version {
	// Keep major and minor, bump patch, everything else -> zero value
	return Version{
		Major: current.Major,
		Minor: current.Minor,
		Patch: current.Patch + 1,
	}
}

// IsValid returns whether or not a string is a valid semantic version.
func IsValid(text string) bool {
	return semVerRegex.MatchString(text)
}
