/*
Package clerk provides file persistence for your Go structures.

It can replace external databases for small projects that can keep data in memory and don't make a lot of writes.
*/
package clerk

import (
	"encoding/gob"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	errDBInited    = errors.New("clerk: database inited multiple times")
	errDBNotInited = errors.New("clerk: database not inited")
)

// DB provides methods to persist your data.
type DB struct {
	Time        time.Time // Time is the timestamp of the last data modification (calling DB.Lock).
	filename    string
	tmpFilename string
	parent      interface{}
	mu          sync.RWMutex
}

// DBInterface represents an underlying database.
type DBInterface interface {
	setFilename(string)
	setParent(interface{})
	inited() bool
	Save() error
	Rebase() error
	Remove() error
	Lock()
	RLock()
	Unlock()
	RUnlock()
}

// Init sets the filename for source saving.
func Init(filename string, db DBInterface) {
	if db.inited() {
		panic(errDBInited)
	}
	db.setFilename(filename)
	db.setParent(db)
}

func (d *DB) setFilename(name string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.filename = name
	// Set d.tmpFilename
	dir := filepath.Dir(name)
	base := filepath.Base(name)
	d.tmpFilename = filepath.Join(dir, "~"+base)
}

func (d *DB) setParent(p interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.parent = p
}

func (d *DB) inited() bool {
	return d.filename != "" && d.tmpFilename != "" && d.parent != nil
}

// Save persists the database in the file set on init.
// Be sure the database is locked for reading.
func (d *DB) Save() error {
	if !d.inited() {
		return errDBNotInited
	}

	tmpFile, err := os.OpenFile(d.tmpFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	if err = gob.NewEncoder(tmpFile).Encode(d.parent); err != nil && err != io.EOF {
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

// Rebase replaces the data with the content of the file set on init.
// Be sure the database is locked for writing.
func (d *DB) Rebase() error {
	if !d.inited() {
		return errDBNotInited
	}

	file, err := os.OpenFile(d.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = gob.NewDecoder(file).Decode(d.parent); err != nil && err != io.EOF {
		return err
	}
	return nil
}

// Remove deletes the database file.
func (d *DB) Remove() error {
	return os.Remove(d.filename)
}

// Lock locks database for reading and writing.
func (d *DB) Lock() {
	d.Time = time.Now()
	d.mu.Lock()
}

// RLock locks database for reading.
func (d *DB) RLock() {
	d.mu.RLock()
}

// Unlock unlocks database for reading and writing.
func (d *DB) Unlock() {
	d.mu.Unlock()
}

// RUnlock unlocks database for reading.
func (d *DB) RUnlock() {
	d.mu.RUnlock()
}
