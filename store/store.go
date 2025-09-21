package store

import (
    "encoding/json"
    "fmt"

    "github.com/syndtr/goleveldb/leveldb"
    "gfn/consensus"
)

var db *leveldb.DB

// OpenDB opens (or creates) the database
func OpenDB(path string) error {
    var err error
    db, err = leveldb.OpenFile(path, nil)
    if err != nil {
        return err
    }
    return nil
}

// CloseDB closes the database
func CloseDB() {
    if db != nil {
        db.Close()
    }
}

// SaveBlock stores a block by its hash
func SaveBlock(block *consensus.Block) error {
    data, err := json.Marshal(block)
    if err != nil {
        return err
    }
    return db.Put([]byte(block.Hash), data, nil)
}

// LoadBlock retrieves a block by its hash
func LoadBlock(hash string) (*consensus.Block, error) {
    data, err := db.Get([]byte(hash), nil)
    if err != nil {
        return nil, err
    }
    var b consensus.Block
    if err := json.Unmarshal(data, &b); err != nil {
        return nil, err
    }
    return &b, nil
}

// SaveHead stores the latest block hash
func SaveHead(hash string) error {
    return db.Put([]byte("HEAD"), []byte(hash), nil)
}

// LoadHead gets the latest block hash
func LoadHead() (string, error) {
    data, err := db.Get([]byte("HEAD"), nil)
    if err != nil {
        return "", err
    }
    return string(data), nil
}
