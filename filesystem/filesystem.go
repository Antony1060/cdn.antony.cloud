package filesystem

import (
	"cdn/db"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/karrick/godirwalk"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileSystem struct {
	IndexDir  string
	SecretDir string
}

type FileLookup struct {
	Index  []File `json:"index"`
	Secret []File `json:"secret"`
}

func (fs *FileSystem) GetAll() (*FileLookup, error) {
	index := make([]File, 0)
	secret := make([]File, 0)

	wg := new(errgroup.Group)
	c := make(chan File)

	if err := filesFromDir(fs.IndexDir, true, wg, c); err != nil {
		return nil, err
	}

	if err := filesFromDir(fs.SecretDir, false, wg, c); err != nil {
		return nil, err
	}

	go func() {
		_ = wg.Wait()
		close(c)
	}()

	for f := range c {
		if f.Indexed {
			index = append(index, f)
		} else {
			secret = append(secret, f)
		}
	}

	return &FileLookup{
		Index:  index,
		Secret: secret,
	}, nil
}

func (fs *FileSystem) Exists(name string, index bool) bool {
	dir := fs.SecretDir
	if index {
		dir = fs.IndexDir
	}

	path, err := filepath.Abs(filepath.Join(dir, name))
	if err != nil {
		return false
	}

	_, err = os.Stat(path)
	return err == nil
}

func (fs *FileSystem) CreateFile(header *multipart.FileHeader, name, password string, timeTillDeath int, indexed, override bool) error {
	database := db.Get()

	dir := fs.SecretDir
	if indexed {
		dir = fs.IndexDir
	}

	f, err := header.Open()
	if err != nil {
		return err
	}

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return err
	}

	sum := h.Sum(nil)

	if name == "" {
		name = hex.EncodeToString(sum) + filepath.Ext(header.Filename)
	}

	if fs.Exists(name, indexed) && !override {
		return fmt.Errorf("file already exists")
	}

	path, err := filepath.Abs(filepath.Join(dir, name))
	if err != nil {
		return err
	}

	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if _, err = io.Copy(w, f); err != nil {
		return err
	}

	if timeTillDeath > 0 {
		database.FileDeathUnix[path] = timeTillDeath
	} else {
		delete(database.FileDeathUnix, path)
	}

	if password != "" {
		pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		database.FilePasswords[path] = string(pass)
	} else {
		delete(database.FilePasswords, path)
	}

	if err = database.Save(); err != nil {
		return err
	}

	return nil
}

func (fs *FileSystem) Get(name string, index bool) (*File, error) {
	dir := fs.SecretDir
	if index {
		dir = fs.IndexDir
	}

	f, err := FromFile(dir, name, index)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func filesFromDir(dir string, indexed bool, wg *errgroup.Group, files chan File) error {
	err := godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(osPathname string, directoryEntry *godirwalk.Dirent) error {
			wg.Go(func() error {
				if filepath.Clean(osPathname) == filepath.Clean(dir) {
					return nil
				}

				f, err := FromFile(dir, directoryEntry.Name(), indexed)
				if err != nil {
					return err
				}

				files <- *f

				return nil
			})
			return nil
		},
	})

	if err != nil {
		return err
	}

	return nil
}
