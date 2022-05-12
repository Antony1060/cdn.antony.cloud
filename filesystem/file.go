package filesystem

import (
	"cdn/db"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Name      string `json:"name"`
	Path      string `json:"-"`
	HashHex   string `json:"hash"`
	Password  string `json:"-"`
	DeathUnix int64  `json:"deathUnix"`
	Indexed   bool   `json:"indexed"`
	usable    bool   // tells if this structure still exists in the filesystem, Delete makes this unusable
}

func (f *File) IsUsable() bool {
	return f.usable
}

func (f *File) IsNamed() bool {
	return f.HashHex != strings.TrimSuffix(f.Name, filepath.Ext(f.Name))
}

func (f *File) Delete() error {
	if err := os.RemoveAll(f.Path); err != nil {
		return err
	}

	database := db.Get()
	delete(database.FilePasswords, f.Path)
	delete(database.FileDeathUnix, f.Path)

	f.usable = false

	return nil
}

func (f *File) HasPassword() bool {
	return f.Password != ""
}

func (f *File) Unlock(pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(f.Password), []byte(pass)) != nil
}

func FromFile(dir, name string, index bool) (*File, error) {
	database := db.Get()

	path, err := filepath.Abs(filepath.Join(dir, name))
	if err != nil {
		return nil, err
	}

	open, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer open.Close()

	h := sha256.New()
	if _, err = io.Copy(h, open); err != nil {
		return nil, err
	}
	sum := h.Sum(nil)

	deathUnix, ok := database.FileDeathUnix[path]
	if !ok {
		deathUnix = 0
	}

	pass, ok := database.FilePasswords[path]
	if !ok {
		pass = ""
	}

	f := &File{
		Name:      name,
		Path:      path,
		HashHex:   hex.EncodeToString(sum),
		Password:  pass,
		DeathUnix: deathUnix,
		Indexed:   index,
		usable:    true,
	}

	return f, nil
}
