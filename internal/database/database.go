package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string 
	m *sync.RWMutex
}
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
} 

type DBStructureUser struct {
	Users map[string]User `json:"users"`
}


type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
}
type User struct {
	Password string `json:"password"`
	Email string `json:"email"`
	ID int `json:"id"`
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
		return db.createDBUser()
	}
	return nil
}

func (db *DB)createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
	}
	return db.writeDB(dbStructure)
} 
func (db *DB)createDBUser() error {
	dbStructure := DBStructureUser{
		Users: map[string]User{},
	}
	return db.writeDBUser(dbStructure)
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStructure.Chirps) + 1
	chirp := Chirp {
		ID : id,
		Body: body,
	}

	dbStructure.Chirps[id] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

func(db *DB)CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDBUser()
	if err != nil {
		return User{}, err
	}
	if _, ok := dbStructure.Users[email]; ok {
		return User{}, errors.New("User already exist")
	}
	id := len(dbStructure.Users) + 1
	bytePass := []byte(password)
	costFactor := 12
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePass, costFactor)
	if err != nil {
		return User{}, err
	}
	user := User {
		Password: string(hashedPassword),
		Email: email,
		ID: id,
	}

	dbStructure.Users[email] = user
	err = db.writeDBUser(dbStructure)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func(db *DB)LoginUser(email, password string) (User, error) {
	dbStructure, err := db.loadDBUser()
	if err != nil {
		return User{}, err
	}
	user, ok := dbStructure.Users[email]
	if !ok {
		return User{}, errors.New("Email didn't match")
	}
	getPassword := []byte(user.Password)
	givenPassword := []byte(password)
	err =  bcrypt.CompareHashAndPassword(getPassword, givenPassword)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func(db *DB)GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	
	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
} 

func(db *DB) GetChirpById(chirpId int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp, ok := dbStructure.Chirps[chirpId]
	if !ok {
		return Chirp{}, errors.New("Chirp doesn't exist")
	}
	return chirp, nil
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

func (db *DB)loadDBUser() (DBStructureUser, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	dbStructure := DBStructureUser{}
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
func (db *DB)writeDBUser(dbStructure DBStructureUser) error {
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