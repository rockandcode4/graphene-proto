package node

import (
    "fmt"
    "gfn/consensus"
)

func InitNode() error {
    fmt.Println("Initializing node...")
    consensus.InitGenesis()
    consensus.Validators = []consensus.Validator{
        {Address: "validator1", Stake: 1000, Active: true},
        {Address: "validator2", Stake: 800, Active: true},
    }
    fmt.Println("Validators initialized.")
    return nil
}

func StartNode() error {
    fmt.Println("Starting node...")
    go consensus.RunConsensus()
    select {} // keep running
}
