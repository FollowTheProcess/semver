package semver_test

import (
	"reflect"
	"testing"

	"github.com/FollowTheProcess/semver"
)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.version.Tag(); got != tt.want {
				t.Errorf("got %q, wanted %q", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    semver.Version
		wantErr bool
	}{
		{
			name:    "simple",
			text:    "1.2.4",
			want:    semver.Version{Major: 1, Minor: 2, Patch: 4},
			wantErr: false,
		},
		{
			name:    "simple with v",
			text:    "v1.2.4",
			want:    semver.Version{Major: 1, Minor: 2, Patch: 4},
			wantErr: false,
		},
		{
			name:    "prerelease",
			text:    "v2.3.7-rc.1",
			want:    semver.Version{Major: 2, Minor: 3, Patch: 7, Prerelease: "rc.1"},
			wantErr: false,
		},
		{
			name:    "prerelease and build",
			text:    "v8.1.0-rc.1+build.123",
			want:    semver.Version{Major: 8, Minor: 1, Patch: 0, Prerelease: "rc.1", Build: "build.123"},
			wantErr: false,
		},
		{
			name:    "beta",
			text:    "1.2.3-beta",
			want:    semver.Version{Major: 1, Minor: 2, Patch: 3, Prerelease: "beta"},
			wantErr: false,
		},
		{
			name:    "obviously wrong",
			text:    "moby dick",
			want:    semver.Version{},
			wantErr: true,
		},
		{
			name:    "invalid",
			text:    "1",
			want:    semver.Version{},
			wantErr: true,
		},
		{
			name:    "prerelease digits",
			text:    "1.2.3-0123",
			want:    semver.Version{},
			wantErr: true,
		},
		{
			name:    "extra parts",
			text:    "1.2.3.4",
			want:    semver.Version{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := semver.Parse(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() returned %v, wanted %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, wanted %#v", got, tt.want)
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
