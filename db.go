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

// Errors.
var (
	ErrDBInited     = errors.New("clerk: database inited multiple times")
	ErrDBNotInited  = errors.New("clerk: database not inited")
	ErrOldMigration = errors.New("clerk: database in file is newer than migration timestamp")
)

// DB provides methods to persist your data.
type DB struct {
	UpdatedAt   time.Time // UpdatedAt is the timestamp of the last data modification (by calling DB.Lock).
	filename    string
	tmpFilename string
	parent      interface{}
	mu          sync.RWMutex
}

// DBInterface represents an underlying database.
type DBInterface interface {
	updatedAt() time.Time
	setFilename(string)
	setParent(interface{})
	inited() bool
	Touch()
	Save() error
	Migrate(string, DBInterface, func() error) error
	MigrateOrRebase(string, DBInterface, func() error)
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
		panic(ErrDBInited)
	}
	db.setFilename(filename)
	db.setParent(db)
}

func (db *DB) updatedAt() time.Time {
	return db.UpdatedAt
}

func (db *DB) setFilename(name string) {
	db.filename = name
	// Set d.tmpFilename
	dir := filepath.Dir(name)
	base := filepath.Base(name)
	db.tmpFilename = filepath.Join(dir, base+"~")
}

func (db *DB) setParent(p interface{}) {
	db.parent = p
}

func (db *DB) inited() bool {
	return db.filename != "" && db.tmpFilename != "" && db.parent != nil
}

// Touch updates the UpdatedAt field.
func (db *DB) Touch() {
	db.UpdatedAt = time.Now()
}

// Save persists the database in the file set on init.
// Be sure the database is locked for writing.
func (db *DB) Save() error {
	if !db.inited() {
		return ErrDBNotInited
	}
	db.Touch()
	tmpFile, err := os.OpenFile(db.tmpFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	if err = gob.NewEncoder(tmpFile).Encode(db.parent); err != nil && err != io.EOF {
		tmpFile.Close()
		return err
	}
	if err = tmpFile.Close(); err != nil {
		return err
	}
	return os.Rename(db.tmpFilename, db.filename)
}

// Rebase replaces the data with the content of the file set on init.
// Be sure the database is locked for writing.
func (db *DB) Rebase() error {
	if !db.inited() {
		return ErrDBNotInited
	}

	file, err := os.OpenFile(db.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = gob.NewDecoder(file).Decode(db.parent); err != nil && err != io.EOF {
		return err
	}
	return nil
}

// Migrate runs a migration if provided date is after the one registered in the data file.
// It rebases the database file in oldDB, runs function do and saves database.
// timestamp format is:
//	2006-01-02 15:04:05 -07
func (db *DB) Migrate(timestamp string, oldDB DBInterface, do func() error) error {
	if !db.inited() {
		return ErrDBNotInited
	}
	ts, err := time.Parse("2006-01-02 15:04:05 -07", timestamp)
	if err != nil {
		return err
	}
	Init(db.filename, oldDB)
	if err = oldDB.Rebase(); err != nil {
		return err
	}
	if !ts.After(oldDB.updatedAt()) {
		return ErrOldMigration
	}
	if err = do(); err != nil {
		return err
	}
	return db.Save()
}

// MigrateOrRebase runs Migrate.
// If the migration has not been used, database is rebased.
// It panics on error.
func (db *DB) MigrateOrRebase(timestamp string, oldDB DBInterface, do func() error) {
	if err := db.Migrate(timestamp, oldDB, do); err != nil {
		if err != ErrOldMigration {
			panic(err)
		}
		if err = db.Rebase(); err != nil {
			panic(err)
		}
	}
}

// Remove deletes the database file.
func (db *DB) Remove() error {
	if err := os.Remove(db.filename); !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Lock locks database for reading and writing.
func (db *DB) Lock() {
	db.Touch()
	db.mu.Lock()
}

// RLock locks database for reading.
func (db *DB) RLock() {
	db.mu.RLock()
}

// Unlock unlocks database for reading and writing.
func (db *DB) Unlock() {
	db.mu.Unlock()
}

// RUnlock unlocks database for reading.
func (db *DB) RUnlock() {
	db.mu.RUnlock()
}
