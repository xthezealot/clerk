/*
Package clerk provides file persistence for your Go variables.

It can replace external databases for small projects that can keep data in memory and don't do a lot of writes.
*/
package clerk

import (
	"encoding/gob"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// DB provides methods to persist your data.
type DB interface {
	Save() error
	Remove() error
}

type db struct {
	filename    string
	tmpFilename string
	source      interface{}
	mu          sync.Mutex
}

// New makes a new database.
// It decodes the named file in the data source (a pointer to in-memory data).
func New(filename string, source interface{}) (DB, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	db := &db{
		filename:    filename,
		tmpFilename: tmpFilename(filename),
		source:      source,
	}
	if err = gob.NewDecoder(file).Decode(source); err != nil && err != io.EOF {
		return nil, err
	}
	return db, nil
}

// Save encodes all the source data in gob format and updates the data file if there is no error.
func (d *db) Save() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	tmpFile, err := os.OpenFile(d.tmpFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	if err = gob.NewEncoder(tmpFile).Encode(d.source); err != nil && err != io.EOF {
		tmpFile.Close()
		return err
	}
	if err = tmpFile.Close(); err != nil {
		return err
	}
	if err = os.Rename(d.tmpFilename, d.filename); err != nil {
		return err
	}
	return nil
}

// Remove deletes the database file.
func (d *db) Remove() error {
	return os.Remove(d.filename)
}

func tmpFilename(filename string) string {
	dir := filepath.Dir(filename)
	base := filepath.Base(filename)
	return filepath.Join(dir, "~"+base)
}
