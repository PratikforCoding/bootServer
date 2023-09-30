package database

import(
	"errors"
)
type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
	Author_id int `json:"author-id"`
}

func (db *DB) CreateChirp(body string, author_id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStructure.Chirps) + 1
	chirp := Chirp {
		ID : id,
		Body: body,
		Author_id: author_id,
	}

	dbStructure.Chirps[id] = chirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
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