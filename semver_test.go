package semver

import "testing"

func TestHello(t *testing.T) {
	got := Hello()
	want := "Hello semver"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}
