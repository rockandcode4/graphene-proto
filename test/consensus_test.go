package test

import (
    "gfn/consensus"
    "testing"
)

func TestGenesisBlock(t *testing.T) {
    consensus.InitGenesis()
    if len(consensus.Blockchain) == 0 {
        t.Fatal("Genesis block not created")
    }
}
