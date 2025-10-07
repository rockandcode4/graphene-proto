package state

import (
    "encoding/json"
    "fmt"

    "github.com/syndtr/goleveldb/leveldb"
    "gfn/store"
)

type Account struct {
    Address string  `json:"address"`
    Balance float64 `json:"balance"`
}

var accounts = make(map[string]*Account)
var db *leveldb.DB

// Initialize the state (open DB reference)
func InitState() {
    db = store.GetDB()
    if db == nil {
        fmt.Println("state: no database found")
        return
    }
    loadAccounts()
}

func loadAccounts() {
    iter := db.NewIterator(nil, nil)
    for iter.Next() {
        key := string(iter.Key())
        if len(key) > 4 && key[:4] == "acc:" {
            var acc Account
            json.Unmarshal(iter.Value(), &acc)
            accounts[acc.Address] = &acc
        }
    }
    iter.Release()
    fmt.Printf("Loaded %d accounts\n", len(accounts))
}

func SaveAccount(acc *Account) error {
    data, _ := json.Marshal(acc)
    return db.Put([]byte("acc:"+acc.Address), data, nil)
}

func GetBalance(addr string) float64 {
    if acc, ok := accounts[addr]; ok {
        return acc.Balance
    }
    return 0
}

func Transfer(from, to string, amount float64) error {
    if GetBalance(from) < amount {
        return fmt.Errorf("insufficient balance")
    }

    // Deduct
    accounts[from].Balance -= amount

    // Add to receiver
    if _, ok := accounts[to]; !ok {
        accounts[to] = &Account{Address: to, Balance: 0}
    }
    accounts[to].Balance += amount

    // Save both
    SaveAccount(accounts[from])
    SaveAccount(accounts[to])

    fmt.Printf("Transferred %.2f GFN from %s to %s\n", amount, from, to)
    return nil
}

func Credit(addr string, amount float64) {
    if _, ok := accounts[addr]; !ok {
        accounts[addr] = &Account{Address: addr, Balance: 0}
    }
    accounts[addr].Balance += amount
    SaveAccount(accounts[addr])
}

func Accounts() map[string]*Account {
    return accounts
}
