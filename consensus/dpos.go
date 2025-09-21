package consensus

import (
    "fmt"
    "time"
    "gfn/store"
)

type Validator struct {
    Address string
    Stake   int64
    Active  bool
}

var Validators []Validator
var Blockchain []*Block

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

// ProduceBlock makes a new block from a validator
func ProduceBlock(validator Validator, data []byte) *Block {
    prev := Blockchain[len(Blockchain)-1]
    block := NewBlock(prev.Height+1, prev.Hash, validator.Address, data)
    Blockchain = append(Blockchain, block)

    // persist
    store.SaveBlock(block)
    store.SaveHead(block.Hash)

    return block
}

// RunConsensus loops through validators to produce blocks
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
