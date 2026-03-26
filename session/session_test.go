package session

import (
	"os"
	"testing"
)

func withTempHome(t *testing.T) func() {
	t.Helper()
	tmp := t.TempDir()
	orig := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	return func() { os.Setenv("HOME", orig) }
}

func TestSaveAndLoad(t *testing.T) {
	defer withTempHome(t)()

	want := &Session{
		BindDN:   "uid=jsmith,ou=users,dc=example,dc=com",
		Password: "secret",
		Username: "jsmith",
	}

	if err := Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if *got != *want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestClear(t *testing.T) {
	defer withTempHome(t)()

	if err := Save(&Session{Username: "u"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	_, err := Load()
	if err == nil {
		t.Error("expected error after clearing session")
	}
}

func TestClearWhenNoSession(t *testing.T) {
	defer withTempHome(t)()

	// Clear on a non-existent session should not error
	if err := Clear(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLoadMissingSession(t *testing.T) {
	defer withTempHome(t)()

	_, err := Load()
	if err == nil {
		t.Error("expected error for missing session file")
	}
}

func TestSessionFilePermissions(t *testing.T) {
	defer withTempHome(t)()

	if err := Save(&Session{Username: "u", Password: "p"}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	p, _ := path()
	info, err := os.Stat(p)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected 0600 permissions, got %o", perm)
	}
}
