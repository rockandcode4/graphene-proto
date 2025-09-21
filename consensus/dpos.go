package consensus

import (
    "fmt"
    "time"
)

type Validator struct {
    Address string
    Stake   int64
    Active  bool
}

var Validators []Validator
var Blockchain []*Block

func InitGenesis() {
    genesis := NewBlock(0, "", "genesis", []byte("Genesis Block"))
    Blockchain = append(Blockchain, genesis)
    fmt.Println("Genesis block created:", genesis.Hash)
}

func ProduceBlock(validator Validator, data []byte) *Block {
    prev := Blockchain[len(Blockchain)-1]
    block := NewBlock(prev.Height+1, prev.Hash, validator.Address, data)
    Blockchain = append(Blockchain, block)
    return block
}

func RunConsensus() {
    for {
        for _, v := range Validators {
            if v.Active {
                b := ProduceBlock(v, []byte("tx data"))
                fmt.Println("Block produced by", v.Address, "at height", b.Height)
                time.Sleep(2 * time.Second)
            }
        }
    }
}
