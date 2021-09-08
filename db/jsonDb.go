package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"os"
	"path/filepath"
)

var data *JsonDb

type JsonDb struct {
	TTLs map[string]int `json:"ttls"`
}

func (db *JsonDb) Save() error {
	f, err := os.Create(filepath.Join(".", "data.json"))
	if err != nil {
		return err
	}

	defer f.Close()

	j, err := json.Marshal(db)
	if err != nil {
		return err
	}
	fmt.Println(string(j[:]))

	w := bufio.NewWriter(f)
	defer w.Flush()
	bytes, err := w.WriteString(string(j[:]))
	if err != nil {
		return err
	}
	log.WithField("size", bytes).Info("Saved data")
	return nil
}

func (db *JsonDb) Load() error {
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

func Get() *JsonDb {
	return data
}

func Init() error {
	db := &JsonDb{ TTLs: map[string]int{} }
	err := db.Load()
	err = db.Save()
	if err != nil {
		return err
	}
	data = db
	return nil
}