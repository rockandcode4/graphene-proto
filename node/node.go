package node

import (
    "fmt"
    "gfn/consensus"
    "gfn/store"
)

func InitNode() error {
    fmt.Println("Initializing node...")
    if err := store.OpenDB("data"); err != nil {
        return err
    }
    defer store.CloseDB()

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
    if err := store.OpenDB("data"); err != nil {
        return err
    }
    defer store.CloseDB()

    go consensus.RunConsensus()
    select {} // block forever
}
