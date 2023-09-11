package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var ErrNotExist = errors.New("resource does not exist")

type DB struct {
	path string 
	m *sync.RWMutex
}
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"user"`
} 

func NewDB(path string) (*DB,  error) {
	db := &DB {
		path: path,
		m: &sync.RWMutex{},
	}
	err := db.ensureDb()
	return db, err
}

func (db *DB)ensureDb() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return nil
}

func (db *DB)createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users: map[int]User{},
	}
	return db.writeDB(dbStructure)
} 

func (db *DB)loadDB() (DBStructure, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}
	return dbStructure, nil
}

func (db *DB) ResetDB() error {
	err := os.Remove(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return db.ensureDb()
}

func (db *DB)writeDB(dbStructure DBStructure) error {
	db.m.Lock()
	defer db.m.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}