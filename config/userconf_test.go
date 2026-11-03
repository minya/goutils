package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

type settings struct {
	Token string `json:"token"`
}

// F60: only a missing file is a NoFileError; an unreadable one is not.
func TestUnmarshalJsonDistinguishesMissingFromUnreadable(t *testing.T) {
	dir := t.TempDir()
	var s settings

	err := UnmarshalJson(&s, filepath.Join(dir, "absent.json"))
	var nf *NoFileError
	if !errors.As(err, &nf) || !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("missing file: err = %v, want a NoFileError matching fs.ErrNotExist", err)
	}

	// A path through a regular file fails with ENOTDIR, not "missing".
	file := filepath.Join(dir, "file")
	os.WriteFile(file, []byte("{}"), 0o600)
	err = UnmarshalJson(&s, filepath.Join(file, "settings.json"))
	if err == nil || errors.As(err, &nf) {
		t.Fatalf("ENOTDIR: err = %v, want a non-NoFileError error", err)
	}

	if os.Getuid() != 0 {
		locked := filepath.Join(dir, "locked.json")
		os.WriteFile(locked, []byte(`{"token":"x"}`), 0o000)
		err = UnmarshalJson(&s, locked)
		if err == nil || errors.As(err, &nf) || !errors.Is(err, fs.ErrPermission) {
			t.Fatalf("unreadable: err = %v, want a permission error, not NoFileError", err)
		}
	}
}

// F60: settings are written 0600 and atomically, over an existing file too.
func TestMarshalJsonWritesPrivateFileAtomically(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "settings.json")
	os.WriteFile(p, []byte(`{"token":"old"}`), 0o644)

	if err := MarshalJson(settings{Token: "secret"}, p); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Errorf("mode = %v, want 0600", fi.Mode().Perm())
	}
	var got settings
	if err := UnmarshalJson(&got, p); err != nil || got.Token != "secret" {
		t.Fatalf("read back %+v, %v", got, err)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 1 {
		t.Errorf("temporary files left behind: %v", entries)
	}
}
