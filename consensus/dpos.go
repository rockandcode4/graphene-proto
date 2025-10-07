package consensus

import (
    "fmt"
    "time"

    "gfn/store"
    "gfn/state"
)

// ----------------- Core Types -----------------

type Validator struct {
    Address string
    Stake   int64
    Active  bool
}

type Block struct {
    Height    int
    PrevHash  string
    Producer  string
    Data      []byte
    Hash      string
    Timestamp int64
}

type Transaction struct {
    From   string
    To     string
    Amount float64
}

// ----------------- Globals -----------------

var Validators []Validator
var Blockchain []*Block

// ----------------- Block Functions -----------------

// NewBlock creates a new block
func NewBlock(height int, prevHash string, producer string, data []byte) *Block {
    block := &Block{
        Height:    height,
        PrevHash:  prevHash,
        Producer:  producer,
        Data:      data,
        Timestamp: time.Now().Unix(),
    }
    block.Hash = fmt.Sprintf("%x", time.Now().UnixNano()) // simple unique hash
    return block
}

// ----------------- Genesis + Load -----------------

// InitGenesis creates the genesis block
func InitGenesis() {
    genesis := NewBlock(0, "", "genesis", []byte("Genesis Block"))
    Blockchain = append(Blockchain, genesis)
    store.SaveBlock(genesis)
    store.SaveHead(genesis.Hash)
    fmt.Println("Genesis block created:", genesis.Hash)
}

// LoadBlockchain reloads the chain from LevelDB
func LoadBlockchain() error {
    headHash, err := store.LoadHead()
    if err != nil {
        return fmt.Errorf("no chain found, need to create genesis: %v", err)
    }

    // Walk backwards from head block until genesis
    hash := headHash
    var chain []*Block
    for {
        block, err := store.LoadBlock(hash)
        if err != nil {
            return fmt.Errorf("failed to load block: %v", err)
        }
        chain = append([]*Block{block}, chain...) // prepend

        if block.Height == 0 {
            break
        }
        hash = block.PrevHash
    }

    Blockchain = chain
    fmt.Println("Blockchain loaded, height:", len(Blockchain)-1)
    return nil
}

// ----------------- Transactions -----------------

func ApplyTransaction(tx Transaction) {
    fmt.Println("Applying transaction:", tx.From, "→", tx.To, tx.Amount)
    if err := state.Transfer(tx.From, tx.To, tx.Amount); err != nil {
        fmt.Println("Transaction failed:", err)
    }
}

// ----------------- Block Production -----------------

func ProduceBlock(validator Validator, data []byte) *Block {
    prev := Blockchain[len(Blockchain)-1]

    // Example: validator pays user1 in each block
    tx := Transaction{From: validator.Address, To: "user1", Amount: 1.5}
    ApplyTransaction(tx)

    block := NewBlock(prev.Height+1, prev.Hash, validator.Address, data)
    Blockchain = append(Blockchain, block)

    store.SaveBlock(block)
    store.SaveHead(block.Hash)

    return block
}

// ----------------- Consensus Loop -----------------

func RunConsensus() {
    for {
        for _, v := range Validators {
            if v.Active {
                b := ProduceBlock(v, []byte("tx data"))
                fmt.Println("Block produced by", v.Address, "at height", b.Height, "hash:", b.Hash)
                time.Sleep(2 * time.Second)
            }
        }
    }
}
