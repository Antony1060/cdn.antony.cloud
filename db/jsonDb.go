package db

import (
	"bufio"
	"encoding/json"
	"github.com/apex/log"
	"os"
	"path/filepath"
)

var data *Database

type Database struct {
	FileDeathUnix map[string]int `json:"file_death_unix"`
	FilePasswords map[string]int `json:"file_passwords"`
}

func (db *Database) Save() error {
	f, err := os.Create(filepath.Join(".", "data.json"))
	if err != nil {
		return err
	}

	defer f.Close()

	j, err := json.Marshal(db)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	defer w.Flush()
	bytes, err := w.WriteString(string(j[:]))
	if err != nil {
		return err
	}
	log.WithField("size", bytes).Info("Saved data")
	return nil
}

func (db *Database) Load() error {
	bytes, err := os.ReadFile(filepath.Join(".", "data.json"))
	if err != nil {
		return nil
	}
	err = json.Unmarshal(bytes, db)
	if err != nil {
		return err
	}
	return nil
}

func Get() *Database {
	return data
}

func Init() error {
	db := &Database{
		FileDeathUnix: map[string]int{},
		FilePasswords: map[string]int{},
	}
	err := db.Load()
	err = db.Save()
	if err != nil {
		return err
	}
	data = db
	return nil
}
