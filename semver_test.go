package semver_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/FollowTheProcess/semver"
)

// Map of valid semver strings to their expected semver.Version, no error should
// be returned parsing these.
var valid = map[string]semver.Version{
	"0.0.4":                 {Major: 0, Minor: 0, Patch: 4},
	"1.2.3":                 {Major: 1, Minor: 2, Patch: 3},
	"10.20.30":              {Major: 10, Minor: 20, Patch: 30},
	"1.1.2-prerelease+meta": {Major: 1, Minor: 1, Patch: 2, Prerelease: "prerelease", Build: "meta"},
	"1.1.2+meta":            {Major: 1, Minor: 1, Patch: 2, Build: "meta"},
	"1.1.2+meta-valid":      {Major: 1, Minor: 1, Patch: 2, Build: "meta-valid"},
	"1.0.0-alpha":           {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha"},
	"1.0.0-beta":            {Major: 1, Minor: 0, Patch: 0, Prerelease: "beta"},
	"1.0.0-alpha.beta":      {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.beta"},
	"1.0.0-alpha.beta.1":    {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.beta.1"},
	"1.0.0-alpha.1":         {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.1"},
	"1.0.0-alpha0.valid":    {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha0.valid"},
	"1.0.0-alpha.0valid":    {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.0valid"},
	"1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay": {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha-a.b-c-somethinglong", Build: "build.1-aef.1-its-okay"},
	"1.0.0-rc.1+build.1":                   {Major: 1, Minor: 0, Patch: 0, Prerelease: "rc.1", Build: "build.1"},
	"2.0.0-rc.1+build.123":                 {Major: 2, Minor: 0, Patch: 0, Prerelease: "rc.1", Build: "build.123"},
	"1.2.3-beta":                           {Major: 1, Minor: 2, Patch: 3, Prerelease: "beta"},
	"10.2.3-DEV-SNAPSHOT":                  {Major: 10, Minor: 2, Patch: 3, Prerelease: "DEV-SNAPSHOT"},
	"1.2.3-SNAPSHOT-123":                   {Major: 1, Minor: 2, Patch: 3, Prerelease: "SNAPSHOT-123"},
	"1.0.0":                                {Major: 1, Minor: 0, Patch: 0},
	"2.0.0":                                {Major: 2, Minor: 0, Patch: 0},
	"1.1.7":                                {Major: 1, Minor: 1, Patch: 7},
	"2.0.0+build.1848":                     {Major: 2, Minor: 0, Patch: 0, Build: "build.1848"},
	"2.0.1-alpha.1227":                     {Major: 2, Minor: 0, Patch: 1, Prerelease: "alpha.1227"},
	"1.0.0-alpha+beta":                     {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha", Build: "beta"},
	"1.2.3----RC-SNAPSHOT.12.9.1--.12+788": {Major: 1, Minor: 2, Patch: 3, Prerelease: "---RC-SNAPSHOT.12.9.1--.12", Build: "788"},
	"1.2.3----R-S.12.9.1--.12+meta":        {Major: 1, Minor: 2, Patch: 3, Prerelease: "---R-S.12.9.1--.12", Build: "meta"},
	"1.2.3----RC-SNAPSHOT.12.9.1--.12":     {Major: 1, Minor: 2, Patch: 3, Prerelease: "---RC-SNAPSHOT.12.9.1--.12"},
	"1.0.0+0.build.1-rc.10000aaa-kk-0.1":   {Major: 1, Minor: 0, Patch: 0, Build: "0.build.1-rc.10000aaa-kk-0.1"},
	"99999999999999999.99999999999999999.99999999999999999": {Major: 99999999999999999, Minor: 99999999999999999, Patch: 99999999999999999},
	"1.0.0-0A.is.legal":      {Major: 1, Minor: 0, Patch: 0, Prerelease: "0A.is.legal"},
	"v0.0.4":                 {Major: 0, Minor: 0, Patch: 4},
	"v1.2.3":                 {Major: 1, Minor: 2, Patch: 3},
	"v10.20.30":              {Major: 10, Minor: 20, Patch: 30},
	"v1.1.2-prerelease+meta": {Major: 1, Minor: 1, Patch: 2, Prerelease: "prerelease", Build: "meta"},
	"v1.1.2+meta":            {Major: 1, Minor: 1, Patch: 2, Build: "meta"},
	"v1.1.2+meta-valid":      {Major: 1, Minor: 1, Patch: 2, Build: "meta-valid"},
	"v1.0.0-alpha":           {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha"},
	"v1.0.0-beta":            {Major: 1, Minor: 0, Patch: 0, Prerelease: "beta"},
	"v1.0.0-alpha.beta":      {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.beta"},
	"v1.0.0-alpha.beta.1":    {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.beta.1"},
	"v1.0.0-alpha.1":         {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.1"},
	"v1.0.0-alpha0.valid":    {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha0.valid"},
	"v1.0.0-alpha.0valid":    {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha.0valid"},
	"v1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay": {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha-a.b-c-somethinglong", Build: "build.1-aef.1-its-okay"},
	"v1.0.0-rc.1+build.1":                   {Major: 1, Minor: 0, Patch: 0, Prerelease: "rc.1", Build: "build.1"},
	"v2.0.0-rc.1+build.123":                 {Major: 2, Minor: 0, Patch: 0, Prerelease: "rc.1", Build: "build.123"},
	"v1.2.3-beta":                           {Major: 1, Minor: 2, Patch: 3, Prerelease: "beta"},
	"v10.2.3-DEV-SNAPSHOT":                  {Major: 10, Minor: 2, Patch: 3, Prerelease: "DEV-SNAPSHOT"},
	"v1.2.3-SNAPSHOT-123":                   {Major: 1, Minor: 2, Patch: 3, Prerelease: "SNAPSHOT-123"},
	"v1.0.0":                                {Major: 1, Minor: 0, Patch: 0},
	"v2.0.0":                                {Major: 2, Minor: 0, Patch: 0},
	"v1.1.7":                                {Major: 1, Minor: 1, Patch: 7},
	"v2.0.0+build.1848":                     {Major: 2, Minor: 0, Patch: 0, Build: "build.1848"},
	"v2.0.1-alpha.1227":                     {Major: 2, Minor: 0, Patch: 1, Prerelease: "alpha.1227"},
	"v1.0.0-alpha+beta":                     {Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha", Build: "beta"},
	"v1.2.3----RC-SNAPSHOT.12.9.1--.12+788": {Major: 1, Minor: 2, Patch: 3, Prerelease: "---RC-SNAPSHOT.12.9.1--.12", Build: "788"},
	"v1.2.3----R-S.12.9.1--.12+meta":        {Major: 1, Minor: 2, Patch: 3, Prerelease: "---R-S.12.9.1--.12", Build: "meta"},
	"v1.2.3----RC-SNAPSHOT.12.9.1--.12":     {Major: 1, Minor: 2, Patch: 3, Prerelease: "---RC-SNAPSHOT.12.9.1--.12"},
	"v1.0.0+0.build.1-rc.10000aaa-kk-0.1":   {Major: 1, Minor: 0, Patch: 0, Build: "0.build.1-rc.10000aaa-kk-0.1"},
	"v99999999999999999.99999999999999999.99999999999999999": {Major: 99999999999999999, Minor: 99999999999999999, Patch: 99999999999999999},
	"v1.0.0-0A.is.legal": {Major: 1, Minor: 0, Patch: 0, Prerelease: "0A.is.legal"},
}

// Invalid semver strings, should all return an error and a zero semver.Version.
var invalid = [...]string{
	"1",
	"1.2",
	"1.2.3-0123",
	"1.2.3-0123.0123",
	"1.1.2+.123",
	"+invalid",
	"-invalid",
	"-invalid+invalid",
	"-invalid.01",
	"alpha",
	"alpha.beta",
	"alpha.beta.1",
	"alpha.1",
	"alpha+beta",
	"alpha_beta",
	"alpha.",
	"alpha..",
	"beta",
	"1.0.0-alpha_beta",
	"-alpha.",
	"1.0.0-alpha..",
	"1.0.0-alpha..1",
	"1.0.0-alpha...1",
	"1.0.0-alpha....1",
	"1.0.0-alpha.....1",
	"1.0.0-alpha......1",
	"1.0.0-alpha.......1",
	"01.1.1",
	"1.01.1",
	"1.1.01",
	"1.2",
	"1.2.3.DEV",
	"1.2-SNAPSHOT",
	"1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788",
	"1.2-RC-SNAPSHOT",
	"-1.0.3-gamma+b7718",
	"+justmeta",
	"9.8.7+meta+meta",
	"9.8.7-whatever+meta+meta",
	"99999999999999999.99999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12",
	"v1",
	"v1.2",
	"v1.2.3-0123",
	"v1.2.3-0123.0123",
	"v1.1.2+.123",
	"v+invalid",
	"v-invalid",
	"v-invalid+invalid",
	"v-invalid.01",
	"valpha",
	"valpha.beta",
	"valpha.beta.1",
	"valpha.1",
	"valpha+beta",
	"valpha_beta",
	"valpha.",
	"valpha..",
	"vbeta",
	"v1.0.0-alpha_beta",
	"v-alpha.",
	"v1.0.0-alpha..",
	"v1.0.0-alpha..1",
	"v1.0.0-alpha...1",
	"v1.0.0-alpha....1",
	"v1.0.0-alpha.....1",
	"v1.0.0-alpha......1",
	"v1.0.0-alpha.......1",
	"v01.1.1",
	"v1.01.1",
	"v1.1.01",
	"v1.2",
	"v1.2.3.DEV",
	"v1.2-SNAPSHOT",
	"v1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788",
	"v1.2-RC-SNAPSHOT",
	"v-1.0.3-gamma+b7718",
	"v+justmeta",
	"v9.8.7+meta+meta",
	"v9.8.7-whatever+meta+meta",
	"v99999999999999999.99999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12",
}

func TestParse(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		for str, want := range valid {
			got, err := semver.Parse(str)
			if err != nil {
				t.Fatalf("Parsing a valid semver string (%q) resulted in an error: %v", str, err)
			}

			if got != want {
				t.Errorf("\nGot:\t%#v\nWanted:\t%#v\n", got, want)
			}
		}
	})

	t.Run("invalid", func(t *testing.T) {
		for _, str := range invalid {
			got, err := semver.Parse(str)
			if err == nil {
				t.Fatalf("Parsing an invalid semver string (%q) did not return an error: %v", str, err)
			}

			want := semver.Version{}
			if got != want {
				t.Errorf("\nGot:\t%#v\nWanted:\t%#v\n", got, want)
			}
		}
	})
}

func TestVersionString(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		version semver.Version
	}{
		{
			name:    "empty",
			version: semver.Version{},
			want:    "0.0.0",
		},
		{
			name:    "just version",
			version: semver.Version{Major: 1, Minor: 6, Patch: 12},
			want:    "1.6.12",
		},
		{
			name:    "prerelease",
			version: semver.Version{Major: 1, Minor: 6, Patch: 12, Prerelease: "rc.1"},
			want:    "1.6.12-rc.1",
		},
		{
			name:    "prerelease and build",
			version: semver.Version{Major: 1, Minor: 6, Patch: 12, Prerelease: "rc.1", Build: "build.123"},
			want:    "1.6.12-rc.1+build.123",
		},
		{
			name:    "complex",
			version: semver.Version{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha-a.b-c-somethinglong", Build: "build.1-aef.1-its-okay"},
			want:    "1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.version.String(); got != tt.want {
				t.Errorf("got %q, wanted %q", got, tt.want)
			}
		})
	}
}

func TestVersionTagString(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		version semver.Version
	}{
		{
			name:    "empty",
			version: semver.Version{},
			want:    "v0.0.0",
		},
		{
			name:    "just version",
			version: semver.Version{Major: 1, Minor: 6, Patch: 12},
			want:    "v1.6.12",
		},
		{
			name:    "prerelease",
			version: semver.Version{Major: 1, Minor: 6, Patch: 12, Prerelease: "rc.1"},
			want:    "v1.6.12-rc.1",
		},
		{
			name:    "prerelease and build",
			version: semver.Version{Major: 1, Minor: 6, Patch: 12, Prerelease: "rc.1", Build: "build.123"},
			want:    "v1.6.12-rc.1+build.123",
		},
		{
			name:    "complex",
			version: semver.Version{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha-a.b-c-somethinglong", Build: "build.1-aef.1-its-okay"},
			want:    "v1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.version.Tag(); got != tt.want {
				t.Errorf("got %q, wanted %q", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	want := semver.Version{
		Prerelease: "rc.1",
		Build:      "build.123",
		Major:      4,
		Minor:      16,
		Patch:      3,
	}
	got := semver.New(4, 16, 3, "build.123", "rc.1")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, wanted %#v", got, want)
	}
}

func TestBumpMajor(t *testing.T) {
	tests := []struct {
		name    string
		current semver.Version
		want    semver.Version
	}{
		{
			name:    "zeros",
			current: semver.Version{},
			want:    semver.Version{Major: 1},
		},
		{
			name:    "minor",
			current: semver.Version{Minor: 1},
			want:    semver.Version{Major: 1},
		},
		{
			name:    "patch",
			current: semver.Version{Patch: 1},
			want:    semver.Version{Major: 1},
		},
		{
			name:    "everything",
			current: semver.Version{Major: 0, Minor: 32, Patch: 6, Prerelease: "rc.1", Build: "build.123"},
			want:    semver.Version{Major: 1},
		},
		{
			name:    "big numbers",
			current: semver.Version{Major: 123, Minor: 32, Patch: 6, Prerelease: "rc.1", Build: "build.123"},
			want:    semver.Version{Major: 124},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := semver.BumpMajor(tt.current); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, wanted %#v", got, tt.want)
			}
		})
	}
}

func TestBumpMinor(t *testing.T) {
	tests := []struct {
		name    string
		current semver.Version
		want    semver.Version
	}{
		{
			name:    "zeros",
			current: semver.Version{},
			want:    semver.Version{Minor: 1},
		},
		{
			name:    "minor",
			current: semver.Version{Minor: 1},
			want:    semver.Version{Minor: 2},
		},
		{
			name:    "patch",
			current: semver.Version{Patch: 1},
			want:    semver.Version{Minor: 1},
		},
		{
			name:    "everything",
			current: semver.Version{Major: 0, Minor: 32, Patch: 6, Prerelease: "rc.1", Build: "build.123"},
			want:    semver.Version{Minor: 33},
		},
		{
			name:    "big numbers",
			current: semver.Version{Major: 123, Minor: 32, Patch: 6, Prerelease: "rc.1", Build: "build.123"},
			want:    semver.Version{Major: 123, Minor: 33},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := semver.BumpMinor(tt.current); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, wanted %#v", got, tt.want)
			}
		})
	}
}

func TestBumpPatch(t *testing.T) {
	tests := []struct {
		name    string
		current semver.Version
		want    semver.Version
	}{
		{
			name:    "zeros",
			current: semver.Version{},
			want:    semver.Version{Patch: 1},
		},
		{
			name:    "minor",
			current: semver.Version{Minor: 1},
			want:    semver.Version{Minor: 1, Patch: 1},
		},
		{
			name:    "patch",
			current: semver.Version{Patch: 1},
			want:    semver.Version{Patch: 2},
		},
		{
			name:    "everything",
			current: semver.Version{Major: 0, Minor: 32, Patch: 6, Prerelease: "rc.1", Build: "build.123"},
			want:    semver.Version{Minor: 32, Patch: 7},
		},
		{
			name:    "big numbers",
			current: semver.Version{Major: 123, Minor: 32, Patch: 6, Prerelease: "rc.1", Build: "build.123"},
			want:    semver.Version{Major: 123, Minor: 32, Patch: 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := semver.BumpPatch(tt.current); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, wanted %#v", got, tt.want)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{
			name: "simple",
			text: "1.2.4",
			want: true,
		},
		{
			name: "simple with v",
			text: "v1.2.4",
			want: true,
		},
		{
			name: "prerelease",
			text: "v2.3.7-rc.1",
			want: true,
		},
		{
			name: "prerelease and build",
			text: "v8.1.0-rc.1+build.123",
			want: true,
		},
		{
			name: "beta",
			text: "1.2.3-beta",
			want: true,
		},
		{
			name: "obviously wrong",
			text: "moby dick",
			want: false,
		},
		{
			name: "invalid",
			text: "1",
			want: false,
		},
		{
			name: "prerelease digits",
			text: "1.2.3-0123",
			want: false,
		},
		{
			name: "extra parts",
			text: "1.2.3.4",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := semver.IsValid(tt.text)

			if got != tt.want {
				t.Errorf("got %#v, wanted %#v", got, tt.want)
			}
		})
	}
}

// FuzzVersionParse fuzzes our parse method with random input and makes sure it never panics.
func FuzzVersionParse(f *testing.F) {
	// Combine all the valid and invalid examples as the corpus
	var all []string
	for str := range valid {
		all = append(all, str)
	}
	for _, str := range invalid {
		all = append(all, str)
	}

	for _, example := range all {
		f.Add(example)
	}

	f.Fuzz(func(t *testing.T, s string) {
		semver.Parse(s) //nolint: errcheck
	})
}

func BenchmarkVersionParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := semver.Parse("v12.4.3-rc1+build.123")
		if err != nil {
			b.Fatalf("Parse returned an error: %v", err)
		}
	}
}

func BenchmarkVersionString(b *testing.B) {
	v := semver.Version{
		Prerelease: "rc1",
		Build:      "build.123",
		Major:      3,
		Minor:      4,
		Patch:      12,
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.String()
	}
}

func BenchmarkVersionTag(b *testing.B) {
	v := semver.Version{
		Prerelease: "rc1",
		Build:      "build.123",
		Major:      3,
		Minor:      4,
		Patch:      12,
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = v.Tag()
	}
}

func ExampleParse() {
	v, err := semver.Parse("v1.19.0")
	if err != nil {
		fmt.Println("Uh oh! That's not a valid semantic version")
	}
	fmt.Println(v)
	// Output: 1.19.0
}

func ExampleBumpMajor() {
	current, err := semver.Parse("3.12.0")
	if err != nil {
		fmt.Println("could not parse")
	}
	next := semver.BumpMajor(current)
	fmt.Println(next)
	// Output: 4.0.0
}

func ExampleBumpMinor() {
	current, err := semver.Parse("3.12.0")
	if err != nil {
		fmt.Println("could not parse")
	}
	next := semver.BumpMinor(current)
	fmt.Println(next)
	// Output: 3.13.0
}

func ExampleBumpPatch() {
	current, err := semver.Parse("3.12.0")
	if err != nil {
		fmt.Println("could not parse")
	}
	next := semver.BumpPatch(current)
	fmt.Println(next)
	// Output: 3.12.1
}

func ExampleIsValid() {
	// Don't need the 'v' at the start
	one := semver.IsValid("1.19.0")

	// But you can have it if you want
	two := semver.IsValid("v1.19.0")

	// Can handle all the complexity of the semver 2.0.0 spec
	three := semver.IsValid("v6.35.12-rc1+build.123")

	// Obviously wrong
	four := semver.IsValid("I'm not a version")

	fmt.Println(one)
	fmt.Println(two)
	fmt.Println(three)
	fmt.Println(four)
	// Output:
	// true
	// true
	// true
	// false
}
