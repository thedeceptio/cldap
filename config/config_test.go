package config

import (
	"os"
	"path/filepath"
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

	want := &Config{
		Host:           "ldap.example.com",
		Port:           389,
		BaseDN:         "dc=example,dc=com",
		UserSearchBase: "ou=users,dc=example,dc=com",
		UsernameAttr:   "uid",
		UseTLS:         false,
		UseStartTLS:    true,
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

func TestSaveCreatesDir(t *testing.T) {
	defer withTempHome(t)()

	if err := Save(&Config{Host: "h", Port: 389}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	home, _ := os.UserHomeDir()
	info, err := os.Stat(filepath.Join(home, ".cldap"))
	if err != nil {
		t.Fatalf("dir not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected directory")
	}
}

func TestLoadMissingFile(t *testing.T) {
	defer withTempHome(t)()

	_, err := Load()
	if err == nil {
		t.Error("expected error for missing config file")
	}
}

func TestSaveFilePermissions(t *testing.T) {
	defer withTempHome(t)()

	if err := Save(&Config{Host: "h", Port: 389}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	home, _ := os.UserHomeDir()
	info, err := os.Stat(filepath.Join(home, ".cldap", "config.json"))
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected 0600 permissions, got %o", perm)
	}
}
