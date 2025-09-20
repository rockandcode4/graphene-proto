package state

import (
    "encoding/json"
    "fmt"
    "path/filepath"

    "github.com/syndtr/goleveldb/leveldb"
)

type Account struct {
    Address string `json:"address"`
    Balance uint64 `json:"balance"`
    Nonce   uint64 `json:"nonce"`
}

type StateDB struct {
    db *leveldb.DB
}

func NewStateDB(dir string) (*StateDB, error) {
    path := filepath.Join(dir, "db")
    d, err := leveldb.OpenFile(path, nil)
    if err != nil { return nil, err }
    return &StateDB{db: d}, nil
}

func (s *StateDB) Close() error { return s.db.Close() }

func (s *StateDB) GetAccount(addr string) (*Account, error) {
    b, err := s.db.Get([]byte("acct:"+addr), nil)
    if err != nil {
        // not found => zero account
        return &Account{Address: addr, Balance: 0, Nonce:0}, nil
    }
    var a Account
    if err := json.Unmarshal(b, &a); err != nil { return nil, err }
    return &a, nil
}

func (s *StateDB) PutAccount(a *Account) error {
    b, _ := json.Marshal(a)
    return s.db.Put([]byte("acct:"+a.Address), b, nil)
}

func (s *StateDB) Transfer(from, to string, amount uint64) error {
    fa, _ := s.GetAccount(from)
    ta, _ := s.GetAccount(to)
    if fa.Balance < amount { return fmt.Errorf("insufficient balance") }
    fa.Balance -= amount
    ta.Balance += amount
    fa.Nonce += 1
    if err := s.PutAccount(fa); err != nil { return err }
    return s.PutAccount(ta)
}
