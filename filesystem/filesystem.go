package filesystem

import (
	"cdn/db"
	"crypto/sha256"
	"encoding/hex"
	"github.com/karrick/godirwalk"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileSystem struct {
	IndexDir  string
	SecretDir string
}

type File struct {
	Name      string
	HashHex   string
	Password  string
	DeathUnix int
	Indexed   bool
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
	return false // TODO
}

func (fs *FileSystem) CreateFile(f *os.File, password string, timeTillDeath int, saveNamed, indexed bool) error {
	return nil // TODO
}

func (fs *FileSystem) Get(name string, index bool) (*File, error) {
	return nil, nil // TODO
}

func (f *File) IsNamed() bool {
	return f.HashHex != strings.TrimSuffix(f.Name, filepath.Ext(f.Name))
}

func (f *File) Delete() error {
	return nil // TODO
}

func (f *File) HasPassword() bool {
	return f.Password != ""
}

func (f *File) Unlock(pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(f.Password), []byte(pass)) != nil
}

func filesFromDir(dir string, indexed bool, wg *errgroup.Group, files chan File) error {
	database := db.Get()

	err := godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(osPathname string, directoryEntry *godirwalk.Dirent) error {
			wg.Go(func() error {
				if filepath.Clean(osPathname) == filepath.Clean(dir) {
					return nil
				}

				data, err := ioutil.ReadFile(osPathname)
				if err != nil {
					return err
				}

				sum := sha256.Sum256(data)

				path := filepath.Join(dir, directoryEntry.Name())

				deathUnix, ok := database.FileDeathUnix[path]
				if !ok {
					deathUnix = 0
				}

				pass, ok := database.FilePasswords[path]
				if !ok {
					pass = ""
				}

				f := File{
					Name:      directoryEntry.Name(),
					HashHex:   hex.EncodeToString(sum[:]),
					Password:  pass,
					DeathUnix: deathUnix,
					Indexed:   indexed,
				}

				files <- f

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
