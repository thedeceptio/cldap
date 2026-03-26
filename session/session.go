package session

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Session struct {
	BindDN   string `json:"bind_dn"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cldap", "session.json"), nil
}

func Load() (*Session, error) {
	p, err := path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func Save(s *Session) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	d := filepath.Join(home, ".cldap")
	if err := os.MkdirAll(d, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(d, "session.json"), data, 0600)
}

func Clear() error {
	p, err := path()
	if err != nil {
		return err
	}
	err = os.Remove(p)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
